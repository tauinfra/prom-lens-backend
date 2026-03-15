package repository

import (
	"fmt"

	clientsetversioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeConfig struct {
	Host  string  `gorm:"host" json:"host"`
	Token string  `gorm:"token" json:"token"`
	QPS   float32 `gorm:"qps" json:"qps"`
	Burst int     `gorm:"burst" json:"burst"`
}

func NewKubeConfig(host, token string) *KubeConfig {
	qps, burst := kubeClientLimits()
	return &KubeConfig{
		Host:  host,
		Token: token,
		QPS:   qps,
		Burst: burst,
	}
}

// Validate 验证配置
func (c *KubeConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("kubernetes host is required")
	}
	if c.Token == "" {
		return fmt.Errorf("kubernetes token is required")
	}
	return nil
}

func (c *KubeConfig) RestConfig() (*rest.Config, error) {
	// 验证配置
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid kubeconfig: %w", err)
	}
	return &rest.Config{
		Host:        c.Host,
		BearerToken: c.Token,
		QPS:         c.QPS,
		Burst:       c.Burst,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}

func (c *KubeConfig) ClientSet() (*kubernetes.Clientset, error) {
	config, err := c.RestConfig()
	if err != nil {
		return nil, fmt.Errorf("create rest config failed: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes clientset failed: %w", err)
	}

	return client, nil
}

func (c *KubeConfig) TektonClientSet() (*clientsetversioned.Clientset, error) {
	config, err := c.RestConfig()
	if err != nil {
		return nil, fmt.Errorf("create rest config failed: %w", err)
	}

	client, err := clientsetversioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create tekton clientset failed: %w", err)
	}

	return client, nil
}
