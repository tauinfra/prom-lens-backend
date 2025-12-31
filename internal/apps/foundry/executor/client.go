package executor

import (
	"fmt"
	"os"

	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newClients() (*tektonclient.Clientset, *kubernetes.Clientset, error) {
	var (
		config *rest.Config
		err    error
	)

	// 优先使用集群内配置
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		// 如果不在集群内，则使用本地 kubeconfig
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.ExpandEnv("$HOME/.kube/config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build config: %v", err)
	}

	// 创建 Tekton 客户端
	tektonClient, err := tektonclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Tekton client: %v", err)
	}

	// 创建 Kubernetes 客户端（用于获取 Pod 和日志）
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}
	return tektonClient, client, nil
}
