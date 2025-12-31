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

func (f *KubeConfigFactory) GetKubeConfig(ctx context.Context, id int) (*KubeConfig, error) {
	println(id)
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

func (f *KubeConfigFactory) GetRestConfig(ctx context.Context, id int) (*rest.Config, error) {
	clusterWithToken, err := f.cluster.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("kubernetes cluster get failed, err: %v", err)
	}
	decryptToken, err := f.encryptor.Decrypt(clusterWithToken.Token)
	if err != nil {
		return nil, fmt.Errorf("kubernetes decrypt token failed, err: %v", err)
	}
	clusterWithToken.Token = decryptToken
	restConfig := &rest.Config{
		Host:        clusterWithToken.Host,
		BearerToken: clusterWithToken.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // 设置为 true 时，不需要 CA
		},
	}
	return restConfig, nil
}

func (f *KubeConfigFactory) GetClientSet(ctx context.Context, id int) (*kubernetes.Clientset, error) {
	kubeConfig, err := f.GetKubeConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	return kubeConfig.ClientSet()
}

func (f *KubeConfigFactory) GetTektonClientSet(ctx context.Context, id int) (*clientsetversioned.Clientset, error) {
	kubeConfig, err := f.GetKubeConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	return kubeConfig.TektonClientSet()
}
