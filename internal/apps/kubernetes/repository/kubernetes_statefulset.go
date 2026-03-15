package repository

import (
	"context"
	"fmt"
	"time"
	"valyria-backend/internal/pkg/k8s/factory"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type StatefulSetRepository interface {
	List(ctx context.Context, id uint, ns string) ([]StatefulSet, error)
	Get(ctx context.Context, id uint, ns, name string) (*appsv1.StatefulSet, error)
	GetDetail(ctx context.Context, id uint, ns, name string) (StatefulSet, error)
	Create(ctx context.Context, id uint, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Update(ctx context.Context, id uint, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Delete(ctx context.Context, id uint, ns, name string) error
	Scale(ctx context.Context, id uint, ns string, name string, replicas int32) (*autoscalingv1.Scale, error)
	Restart(ctx context.Context, id uint, ns string, name string) (*appsv1.StatefulSet, error)
}

type StatefulSet struct {
	Namespace         string                           `json:"namespace"`
	Name              string                           `json:"name"`
	Replicas          int32                            `json:"replicas"`
	ReadyReplicas     int32                            `json:"readyReplicas"`
	CurrentReplicas   int32                            `json:"currentReplicas"`
	AvailableReplicas int32                            `json:"availableReplicas"`
	UpdatedReplicas   int32                            `json:"updatedReplicas"`
	MatchLabels       map[string]string                `json:"matchLabels"`
	Strategy          appsv1.StatefulSetUpdateStrategy `json:"strategy"`
	CreatedAt         string                           `json:"createdAt"`
}

type statefulSetRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewStatefulSetRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) StatefulSetRepository {
	return &statefulSetRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *statefulSetRepository) List(ctx context.Context, id uint, ns string) (statefulSets []StatefulSet, err error) {
	var (
		statefulSet StatefulSet
		matchLabels string
	)
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().StatefulSets(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 statefulSets list failed. err: %v1", err)
	}
	for _, item := range response.Items {
		for k, v := range item.Spec.Selector.MatchLabels {
			matchLabels += fmt.Sprintf("%v1=%v1,", k, v) // 标签
		}
		statefulSet.Namespace = item.Namespace
		statefulSet.Name = item.Name
		statefulSet.Replicas = item.Status.Replicas
		statefulSet.ReadyReplicas = item.Status.ReadyReplicas
		statefulSet.CurrentReplicas = item.Status.CurrentReplicas
		statefulSet.UpdatedReplicas = item.Status.UpdatedReplicas
		statefulSet.AvailableReplicas = item.Status.AvailableReplicas
		statefulSet.MatchLabels = item.Spec.Selector.MatchLabels
		statefulSet.Strategy = item.Spec.UpdateStrategy
		statefulSet.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		statefulSets = append(statefulSets, statefulSet)
	}
	return statefulSets, nil
}

func (r *statefulSetRepository) GetDetail(ctx context.Context, id uint, ns, name string) (statefulSet StatefulSet, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return statefulSet, err
	}
	response, err := client.AppsV1().StatefulSets(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return statefulSet, fmt.Errorf("kubernetes AppsV1 statefulSet get failed. err: %v", err)
	}

	statefulSet.Namespace = response.Namespace
	statefulSet.Name = response.Name
	statefulSet.Replicas = response.Status.Replicas
	statefulSet.ReadyReplicas = response.Status.ReadyReplicas
	statefulSet.CurrentReplicas = response.Status.CurrentReplicas
	statefulSet.UpdatedReplicas = response.Status.UpdatedReplicas
	statefulSet.AvailableReplicas = response.Status.AvailableReplicas
	statefulSet.MatchLabels = response.Spec.Selector.MatchLabels
	statefulSet.Strategy = response.Spec.UpdateStrategy
	statefulSet.CreatedAt = response.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
	return statefulSet, nil
}

func (r *statefulSetRepository) Get(ctx context.Context, id uint, ns, name string) (*appsv1.StatefulSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().StatefulSets(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 statefulSet get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *statefulSetRepository) Create(ctx context.Context, id uint, ns string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().StatefulSets(ns).Create(ctx, statefulSet, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 statefulSet get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *statefulSetRepository) Update(ctx context.Context, id uint, ns string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().StatefulSets(ns).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 statefulSet update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *statefulSetRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.AppsV1().StatefulSets(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes AppsV1 statefulSet delete failed. err: %v", err)
	}
	return nil
}

// Scale 更新副本数量
func (r *statefulSetRepository) Scale(ctx context.Context, id uint, ns, name string, replicas int32) (*autoscalingv1.Scale, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	// 获取副本数量
	scale, err := client.AppsV1().StatefulSets(ns).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 statefulSet get scale failed. err: %v", err)
	}
	scale.Spec.Replicas = replicas // 更新副本数量
	response, err := client.AppsV1().StatefulSets(ns).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 statefulSet update scale failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

// Restart 重启
func (r *statefulSetRepository) Restart(ctx context.Context, id uint, ns, name string) (*appsv1.StatefulSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	data := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
	response, err := client.AppsV1().StatefulSets(ns).Patch(ctx, name, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 deployment restart failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, err
}
