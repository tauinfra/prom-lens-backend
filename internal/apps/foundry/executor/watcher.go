package executor

import (
	"bufio"
	"fmt"
	"time"
	"valyria-backend/internal/core/logger"

	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (t *PipelineExecutor) waitPipelineInitialized(name string) error {
	timeout := time.After(1 * time.Minute)    // 持续 1 分钟，表示超时
	ticker := time.NewTicker(2 * time.Second) // 4 秒检测一次
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("tekton pipelineRun %s start timeout", name)
		case <-ticker.C:
			response, err := t.tektonClient.TektonV1().PipelineRuns("default").Get(context.Background(), name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					continue // PipelineRun 还没创建，继续等待
				}
				return err
			}
			if response.Status.PipelineSpec != nil && len(response.Status.PipelineSpec.Tasks) > 0 {
				// 拿到 task 顺序
				t.taskOrder = make([]string, 0, len(response.Status.PipelineSpec.Tasks))
				for _, task := range response.Status.PipelineSpec.Tasks {
					t.taskOrder = append(t.taskOrder, task.Name)
				}
				return nil // 成功获取 taskOrder
			}
		}
	}
}
func (t *PipelineExecutor) waitTaskRuns(pipelineRun string) error {
	for i, task := range t.taskOrder {
		taskName := pipelineRun + "-" + task
		logger.Infof("tekton taskRun '%s' (step %d/%d) is waiting to complete.", taskName, i+1, len(t.taskOrder))
		if err := t.watchTaskRun(taskName, task); err != nil {
			return err
		}
		logger.Infof("tekton taskRun '%s' completed successfully", taskName)
		// 如果有下一个任务，等待前一个任务的输出成为下一个任务的输入
		if i < len(t.taskOrder)-1 {
			logger.Infof("tekton taskRun '%s' completed, preparing for next taskRun '%s'",
				taskName, t.taskOrder[i+1])
		}
	}
	return nil
}

func (t *PipelineExecutor) watchTaskRun(taskRun, name string) error {
	watcher, err := t.kubeClient.CoreV1().
		Pods("default").
		Watch(context.Background(), metav1.ListOptions{
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
					if _, ok := t.processedPods.LoadOrStore(key, true); !ok {
						go t.streamLogs(pod.Namespace, pod.Name, c.Name)
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
	req := p.kubeClient.CoreV1().
		Pods(ns).
		GetLogs(pod, &corev1.PodLogOptions{
			Container: container,
			Follow:    true,
		})

	stream, err := req.Stream(context.Background())
	if err != nil {
		return
	}
	defer stream.Close() // 确保关闭

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		logMsg := fmt.Sprintf("[%s] [%s]  %s\n", time.Now().Format("2006-01-02 15:04:05"), container, scanner.Text())
		fmt.Println(logMsg)
		if err = p.log.AppendLog(logMsg); err != nil {
			logger.Errorf("Failed to write log: %v", err)
			return
		}
	}
}
