package executor

import (
	"bufio"
	"fmt"
	"time"
	"valyria-backend/internal/core/logger"

	"context"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

func (p *PipelineExecutor) getSucceededCondition(pr *tektonv1.PipelineRun) *apis.Condition {
	for i := range pr.Status.Conditions {
		if pr.Status.Conditions[i].Type == apis.ConditionSucceeded {
			return &pr.Status.Conditions[i]
		}
	}
	return nil
}

func (p *PipelineExecutor) waitPipelineInitialized(name string) error {
	ctx, cancel := context.WithTimeout(p.ctx, 5*time.Minute)
	defer cancel()

	watcher, err := p.tektonClient.TektonV1().
		PipelineRuns(p.cfg.K8s.Tekton.Namespace).
		Watch(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%s", name),
		})
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("tekton pipelineRun %s start timeout", name)
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return fmt.Errorf("pipelineRun watch channel closed")
			}
			pr, ok := event.Object.(*tektonv1.PipelineRun)
			if !ok {
				continue
			}
			cond := p.getSucceededCondition(pr)
			if cond == nil {
				// controller 还没写状态
				continue
			}
			switch cond.Status {
		case corev1.ConditionUnknown:
			// Pipeline 已开始运行
			if pr.Status.PipelineSpec != nil &&
				len(pr.Status.PipelineSpec.Tasks) > 0 {
				p.taskOrderMutex.Lock()
				p.taskOrder = make([]string, 0, len(pr.Status.PipelineSpec.Tasks))
				for _, task := range pr.Status.PipelineSpec.Tasks {
					p.taskOrder = append(p.taskOrder, task.Name)
				}
				p.taskOrderMutex.Unlock()
				return nil
			}
			case corev1.ConditionFalse:
				return fmt.Errorf("%s", cond.Message)
			case corev1.ConditionTrue:
				// 极少数情况下（task 非常快）会直接成功
				return nil
			}
		}
	}
}

func (p *PipelineExecutor) waitTaskRuns(pipelineRun string) error {
	p.taskOrderMutex.RLock()
	taskOrder := make([]string, len(p.taskOrder))
	copy(taskOrder, p.taskOrder)
	p.taskOrderMutex.RUnlock()

	for i, task := range taskOrder {
		taskName := pipelineRun + "-" + task
		logger.Infof("tekton taskRun '%s' (step %d/%d) is waiting to complete.", taskName, i+1, len(taskOrder))
		if err := p.watchTaskRun(taskName, task); err != nil {
			return err
		}
		logger.Infof("tekton taskRun '%s' completed successfully", taskName)
		// 如果有下一个任务，等待前一个任务的输出成为下一个任务的输入
		if i < len(taskOrder)-1 {
			logger.Infof("tekton taskRun '%s' completed, preparing for next taskRun '%s'",
				taskName, taskOrder[i+1])
		}
	}
	return nil
}

func (p *PipelineExecutor) watchTaskRun(taskRun, name string) error {
	watcher, err := p.client.CoreV1().
		Pods(p.cfg.K8s.Tekton.Namespace).
		Watch(p.ctx, metav1.ListOptions{
			LabelSelector: "tekton.dev/taskRun=" + taskRun,
		})
	if err != nil {
		logger.Errorf("error watching Pods for taskRun %s: %v", name, err)
		return err
	}
	defer watcher.Stop() // 保证退出

	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			logger.Errorf("failed to cast event object to Pod %s.", name)
			continue
		}

		switch pod.Status.Phase {
		case corev1.PodRunning:
			for _, c := range pod.Spec.Containers {
				if c.Name == "step-"+name {
					key := pod.Name + "/" + c.Name
					if _, ok = p.processedPods.LoadOrStore(key, true); !ok {
						go p.streamLogs(pod.Namespace, pod.Name, c.Name)
					}
				}
			}
		case corev1.PodSucceeded:
			return nil // 结束函数，触发 defer
		case corev1.PodFailed:
			return fmt.Errorf("pod '%v' completed with failure", pod.Name)
		}
	}
	return nil
}

func (p *PipelineExecutor) streamLogs(ns, pod, container string) {
	req := p.client.CoreV1().
		Pods(ns).
		GetLogs(pod, &corev1.PodLogOptions{
			Container: container,
			Follow:    true,
		})

	stream, err := req.Stream(p.ctx)
	if err != nil {
		return
	}
	defer stream.Close() // 确保关闭

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		logMsg := fmt.Sprintf("[%s] [%s]  %s\n", time.Now().Format("2006-01-02 15:04:05"), container, scanner.Text())
		if err = p.logWriter.AppendLog(logMsg); err != nil {
			logger.Errorf("Failed to write log: %v", err)
			return
		}
	}
}
