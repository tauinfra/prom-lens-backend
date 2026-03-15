package repository

import (
	"context"
	"fmt"
	"io"
	"strings"
	"valyria-backend/internal/pkg/k8s/factory"
	"valyria-backend/internal/pkg/k8s/helper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

type PodRepository interface {
	List(ctx context.Context, id uint, ns, labelSelector string) ([]Pod, error)
	// ListAll 集群级别 Pod 列表，支持 labelSelector / fieldSelector
	ListAll(ctx context.Context, id uint, labelSelector, fieldSelector string) ([]Pod, error)
	Get(ctx context.Context, id uint, ns, name string) (pod *corev1.Pod, err error)
	GetDetail(ctx context.Context, id uint, ns, name string) (Pod, error)
	GetLogs(ctx context.Context, id uint, ns, name, container string, follow bool, tailLines, sinceSeconds *int64) (io.ReadCloser, context.CancelFunc, error)
	Delete(ctx context.Context, id uint, ns, name string) error
	Executor(ctx context.Context, id uint, ns, name, container string) (remotecommand.Executor, error)
	DebugExecutor(ctx context.Context, id uint, ns, name, container string) (remotecommand.Executor, error)
}

type Pod struct {
	Namespace    string                `json:"namespace"`
	Name         string                `json:"name"`
	PodIP        string                `json:"podIP"`
	HostIP       string                `json:"hostIP"`
	NodeName     string                `json:"nodeName"`
	Status       string                `json:"status"`
	RestartCount int                   `json:"restartCount"`
	Ready        string                `json:"ready"`
	Containers   []helper.Container    `json:"containers"`
	Conditions   []corev1.PodCondition `json:"conditions"`
	Labels       map[string]string     `json:"labels"`
	CreatedAt    string                `json:"createdAt"`
}

type Container struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	Ready        bool   `json:"ready"`
	State        string `json:"state"`
	RestartCount int32  `json:"restartCount"`
}

type podRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewPodRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) PodRepository {
	return &podRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *podRepository) List(ctx context.Context, id uint, ns, labelSelector string) (pods []Pod, err error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	podHelper := helper.NewPodHelper()

	response, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 pods list failed. err: %v", err)
	}
	for _, item := range response.Items {
		var pod Pod
		pod.Namespace = item.Namespace
		pod.Name = item.Name
		pod.PodIP = item.Status.PodIP
		pod.HostIP = item.Status.HostIP
		pod.NodeName = item.Spec.NodeName
		pod.Status = podHelper.GetPodStatus(item)
		pod.RestartCount = podHelper.GetMaxRestartCount(item)
		pod.Labels = item.Labels
		pod.Containers = podHelper.GetContainerStatuses(item)
		pod.Ready = fmt.Sprintf("%v/%v", podHelper.GetMaxReadyCount(item), len(item.Status.ContainerStatuses))
		pod.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		pods = append(pods, pod)
	}
	return pods, nil
}

func (r *podRepository) ListAll(ctx context.Context, id uint, labelSelector, fieldSelector string) (pods []Pod, err error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	podHelper := helper.NewPodHelper()
	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
	}
	response, err := client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 pods list failed. err: %v", err)
	}
	for _, item := range response.Items {
		var pod Pod
		pod.Namespace = item.Namespace
		pod.Name = item.Name
		pod.PodIP = item.Status.PodIP
		pod.HostIP = item.Status.HostIP
		pod.NodeName = item.Spec.NodeName
		pod.Status = podHelper.GetPodStatus(item)
		pod.RestartCount = podHelper.GetMaxRestartCount(item)
		pod.Labels = item.Labels
		pod.Containers = podHelper.GetContainerStatuses(item)
		pod.Ready = fmt.Sprintf("%v/%v", podHelper.GetMaxReadyCount(item), len(item.Status.ContainerStatuses))
		pod.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		pods = append(pods, pod)
	}
	return pods, nil
}

func (r *podRepository) Get(ctx context.Context, id uint, ns, name string) (pod *corev1.Pod, err error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 pods get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *podRepository) GetDetail(ctx context.Context, id uint, ns, name string) (pod Pod, err error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return pod, err
	}
	response, err := client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return pod, fmt.Errorf("kubernetes CoreV1 pods get failed. err: %v", err)
	}

	podHelper := helper.NewPodHelper()
	pod.Namespace = response.Namespace
	pod.Name = response.Name
	pod.PodIP = response.Status.PodIP
	pod.HostIP = response.Status.HostIP
	pod.NodeName = response.Spec.NodeName
	pod.Status = podHelper.GetPodStatus(*response)
	pod.RestartCount = podHelper.GetMaxRestartCount(*response)
	pod.Labels = response.Labels
	pod.Containers = podHelper.GetContainerStatuses(*response)
	pod.Conditions = response.Status.Conditions
	pod.Ready = fmt.Sprintf("%v/%v", podHelper.GetMaxReadyCount(*response), len(response.Status.ContainerStatuses))
	pod.CreatedAt = response.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
	return pod, nil
}

func (r *podRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return err
	}
	err = client.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 pods delete failed. err: %v", err)
	}
	return nil
}

func (r *podRepository) GetLogs(ctx context.Context, id uint, ns, name, container string, follow bool, tailLines, sinceSeconds *int64) (io.ReadCloser, context.CancelFunc, error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, nil, err
	}
	response := client.CoreV1().Pods(ns).GetLogs(name, &corev1.PodLogOptions{
		Container:    container,
		Follow:       follow,       // 如果为 true 默认直接从第一行读取，日志量大的时候会有加载问题
		SinceSeconds: sinceSeconds, // 日志开始时间
		TailLines:    tailLines,    // 从倒数第几行开始读取日志
	})
	// 独立 context，不受 Gin 限制，但支持 cancel
	streamCtx, cancel := context.WithCancel(context.Background())
	stream, err := response.Stream(streamCtx)
	if err != nil {
		cancel() // 避免泄漏
		return nil, nil, fmt.Errorf("kubernetes CoreV1 pods get logs stream failed. err: %v", err)
	}
	return stream, cancel, err
}

func (r *podRepository) Executor(ctx context.Context, id uint, ns, name, container string) (remotecommand.Executor, error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}

	// 1. 首先检查 Pod 状态
	pod, err := client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get pod failed: %w", err)
	}

	if pod.Status.Phase != corev1.PodRunning {
		return nil, fmt.Errorf("pod is not running, current phase: %s", pod.Status.Phase)
	}

	// 2. 检查容器是否存在
	containerExists := false
	for _, c := range pod.Spec.Containers {
		if c.Name == container {
			containerExists = true
			break
		}
	}
	if !containerExists {
		return nil, fmt.Errorf("container %s not found in pod", container)
	}

	// 注意: Pod 和 container 不存在也会创建连接，而且不报错
	response := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(name).
		Namespace(ns).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command: []string{
				"/bin/sh",
				"-c",
				"TERM=xterm-256color; export TERM; [ -x /bin/bash ] " +
					"&& ([ -x /usr/bin/script ] " +
					"&& /usr/bin/script -q -c '/bin/bash' /dev/null " +
					"|| exec /bin/bash) || exec /bin/sh",
			},
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			TTY:    true, // 开启终端, 默认为 false
		}, scheme.ParameterCodec) // scheme.ParameterCodec 应该是 pod 的 GVK(GroupVersion & Kind) 之类的

	cfg, err := r.cfgFactory.GetRestConfigAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", response.URL())
	if err != nil {
		return nil, fmt.Errorf("create spdy executor failed: %w", err)
	}

	return executor, nil
}

func (r *podRepository) DebugExecutor(ctx context.Context, id uint, ns, name, container string) (remotecommand.Executor, error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}

	req := client.CoreV1().RESTClient().Get().
		AbsPath(fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec", ns, name)).
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   strings.Fields("ls -la"), // 测试命令
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	cfg, err := r.cfgFactory.GetRestConfigAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("create spdy executor failed: %w", err)
	}

	return executor, nil
}
