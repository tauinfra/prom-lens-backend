package repository

import (
	"context"
	"fmt"

	clientsetversioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"valyria-backend/internal/pkg/encryption"
)

type KubeConfigFactory struct {
	cluster   ClusterRepository
	encryptor encryption.Encryptor
}

func NewKubeConfigFactory(cluster ClusterRepository, encryptor encryption.Encryptor) *KubeConfigFactory {
	return &KubeConfigFactory{
		cluster:   cluster,
		encryptor: encryptor,
	}
}

func (f *KubeConfigFactory) GetKubeConfig(ctx context.Context, id uint) (*KubeConfig, error) {
	clusterWithToken, err := f.cluster.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("kubernetes cluster get failed, err: %v", err)
	}
	decryptToken, err := f.encryptor.Decrypt(clusterWithToken.Token)
	if err != nil {
		return nil, fmt.Errorf("kubernetes decrypt token failed, err: %v", err)
	}
	clusterWithToken.Token = decryptToken
	kubeConfig := NewKubeConfig(clusterWithToken.Host, clusterWithToken.Token)
	return kubeConfig, nil
}

func (f *KubeConfigFactory) GetRestConfig(ctx context.Context, id uint) (*rest.Config, error) {
	clusterWithToken, err := f.cluster.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("kubernetes cluster get failed, err: %v", err)
	}
	decryptToken, err := f.encryptor.Decrypt(clusterWithToken.Token)
	if err != nil {
		return nil, fmt.Errorf("kubernetes decrypt token failed, err: %v", err)
	}
	clusterWithToken.Token = decryptToken
	kubeQPS, kubeBurst := kubeClientLimits()
	restConfig := &rest.Config{
		Host:        clusterWithToken.Host,
		BearerToken: clusterWithToken.Token,
		QPS:         kubeQPS,
		Burst:       kubeBurst,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // 设置为 true 时，不需要 CA
		},
	}
	return restConfig, nil
}

func (f *KubeConfigFactory) GetClientSet(ctx context.Context, id uint) (*kubernetes.Clientset, error) {
	kubeConfig, err := f.GetKubeConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	return kubeConfig.ClientSet()
}

// GetRestConfigAsUser 返回带 Impersonate 的 rest.Config，username 为空时等同于 GetRestConfig
func (f *KubeConfigFactory) GetRestConfigAsUser(ctx context.Context, id uint, username string) (*rest.Config, error) {
	restConfig, err := f.GetRestConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	if username == "" {
		return restConfig, nil
	}
	cfg := *restConfig
	cfg.Impersonate = rest.ImpersonationConfig{UserName: username}
	return &cfg, nil
}

// GetClientSetAsUser 返回以指定用户身份访问集群的 ClientSet（Impersonate），用于用户维度的 RBAC
// username 为空时行为等同于 GetClientSet（不 Impersonate）
func (f *KubeConfigFactory) GetClientSetAsUser(ctx context.Context, id uint, username string) (*kubernetes.Clientset, error) {
	restConfig, err := f.GetRestConfigAsUser(ctx, id, username)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

func (f *KubeConfigFactory) GetTektonClientSet(ctx context.Context, id uint) (*clientsetversioned.Clientset, error) {
	kubeConfig, err := f.GetKubeConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	return kubeConfig.TektonClientSet()
}
