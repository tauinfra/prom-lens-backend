package executor

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func clientSet() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.ExpandEnv("$HOME/.kube/config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, fmt.Errorf("build k8s config: %w", err)
	}
	return kubernetes.NewForConfig(config)
}
