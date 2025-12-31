package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DaemonSetRepository interface {
	List(ctx context.Context, id int, ns string) ([]DaemonSet, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.DaemonSet, error)
	GetDetail(ctx context.Context, id int, ns, name string) (DaemonSet, error)
	Create(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Update(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Delete(ctx context.Context, id int, ns, name string) error
}

type DaemonSet struct {
	Namespace              string                         `json:"namespace"`
	Name                   string                         `json:"name"`
	CurrentNumberScheduled int32                          `json:"currentNumberScheduled"`
	DesiredNumberScheduled int32                          `json:"desiredNumberScheduled"`
	NumberReady            int32                          `json:"numberReady"`
	NumberAvailable        int32                          `json:"numberAvailable"`
	UpdatedNumberScheduled int32                          `json:"updatedNumberScheduled"`
	MatchLabels            map[string]string              `json:"matchLabels"`
	Strategy               appsv1.DaemonSetUpdateStrategy `json:"strategy"`
	CreatedAt              string                         `json:"createdAt"`
}

type daemonSetRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewDaemonSetRepository(kubeFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) DaemonSetRepository {
	return &daemonSetRepository{
		cfgFactory: kubeFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *daemonSetRepository) List(ctx context.Context, id int, ns string) (daemonSets []DaemonSet, err error) {
	var (
		daemonSet   DaemonSet
		matchLabels string
	)
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().DaemonSets(ns).List(context.TODO(), metav1.ListOptions{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		logger.ZapLogger.Error(fmt.Sprintf("kubernetes AppsV1 daemonSet list failed. err: %v1", err))
		return nil, err
	}
	for _, item := range response.Items {
		for k, v := range item.Spec.Selector.MatchLabels {
			matchLabels += fmt.Sprintf("%v=%v,", k, v) // 标签
		}
		daemonSet.Namespace = item.Namespace
		daemonSet.Name = item.Name
		daemonSet.MatchLabels = item.Spec.Selector.MatchLabels
		daemonSet.CurrentNumberScheduled = item.Status.CurrentNumberScheduled
		daemonSet.DesiredNumberScheduled = item.Status.DesiredNumberScheduled // Desired 期望调度的 Pod 数量
		daemonSet.NumberReady = item.Status.NumberReady                       // Ready 状态的 Pod 数量。 所有容器都处于运行状态，并且通过了就绪探针
		daemonSet.NumberAvailable = item.Status.NumberAvailable               // Available 状态的 Pod 数量。不仅要通过就绪探针，还需要运行一段时间
		daemonSet.UpdatedNumberScheduled = item.Status.UpdatedNumberScheduled
		daemonSet.Strategy = item.Spec.UpdateStrategy
		daemonSet.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		daemonSets = append(daemonSets, daemonSet)
	}
	return
}

func (r *daemonSetRepository) Get(ctx context.Context, id int, ns, name string) (*appsv1.DaemonSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().DaemonSets(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.ZapLogger.Error(fmt.Sprintf("kubernetes AppsV1 daemonSet get failed. err: %v1", err))
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *daemonSetRepository) GetDetail(ctx context.Context, id int, ns, name string) (daemonSet DaemonSet, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return daemonSet, err
	}
	response, err := client.AppsV1().DaemonSets(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.ZapLogger.Error(fmt.Sprintf("kubernetes AppsV1 daemonSet list failed. err: %v1", err))
		return daemonSet, err
	}
	daemonSet.Namespace = response.Namespace
	daemonSet.Name = response.Name
	daemonSet.MatchLabels = response.Spec.Selector.MatchLabels
	daemonSet.CurrentNumberScheduled = response.Status.CurrentNumberScheduled
	daemonSet.DesiredNumberScheduled = response.Status.DesiredNumberScheduled // Desired 期望调度的 Pod 数量
	daemonSet.NumberReady = response.Status.NumberReady                       // Ready 状态的 Pod 数量。 所有容器都处于运行状态，并且通过了就绪探针
	daemonSet.NumberAvailable = response.Status.NumberAvailable               // Available 状态的 Pod 数量。不仅要通过就绪探针，还需要运行一段时间
	daemonSet.UpdatedNumberScheduled = response.Status.UpdatedNumberScheduled
	daemonSet.Strategy = response.Spec.UpdateStrategy
	daemonSet.CreatedAt = response.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
	return daemonSet, nil
}

func (r *daemonSetRepository) Create(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().DaemonSets(ns).Create(context.TODO(), body, metav1.CreateOptions{})
	if err != nil {
		logger.ZapLogger.Error(fmt.Sprintf("kubernetes AppsV1 daemonSet get failed. err: %v1", err))
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *daemonSetRepository) Update(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().DaemonSets(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.ZapLogger.Error(fmt.Sprintf("kubernetes AppsV1 daemonSet update failed. err: %v1", err))
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

// Delete 删除
func (r *daemonSetRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.AppsV1().DaemonSets(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.ZapLogger.Error(fmt.Sprintf("kubernetes AppsV1 daemonSet delete failed. err: %v1", err))
	}
	return err
}
