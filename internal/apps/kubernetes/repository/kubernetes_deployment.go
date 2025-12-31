package repository

import (
	"context"
	"fmt"
	"time"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type DeploymentRepository interface {
	List(ctx context.Context, id int, ns string) ([]Deployment, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.Deployment, error)
	GetDetail(ctx context.Context, id int, ns, name string) (Deployment, error)
	Create(ctx context.Context, id int, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(ctx context.Context, id int, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(ctx context.Context, id int, ns string, name string) error
	Scale(ctx context.Context, id int, ns string, name string, replicas int32) (*autoscalingv1.Scale, error)
	Restart(ctx context.Context, id int, ns string, name string) (*appsv1.Deployment, error)
	Rollout(ctx context.Context, id int, ns string, name, rsName string) (*appsv1.Deployment, error)
}

type Deployment struct {
	Namespace         string                       `json:"namespace"`
	Name              string                       `json:"name"`
	Replicas          int32                        `json:"replicas"`          // 部署期望副本数
	ReadyReplicas     int32                        `json:"readyReplicas"`     // 部署期望副本数
	AvailableReplicas int32                        `json:"availableReplicas"` // 可用副本
	UpdatedReplicas   int32                        `json:"updatedReplicas"`   // 已更新副本
	Paused            bool                         `json:"paused"`
	MatchLabels       map[string]string            `json:"matchLabels"` // Selector 调度标签
	Strategy          appsv1.DeploymentStrategy    `json:"strategy"`
	Conditions        []appsv1.DeploymentCondition `json:"conditions"`
	CreatedAt         string                       `json:"createdAt"`
	UpdatedAt         string                       `json:"updatedAt,omitempty"`
}

type deploymentRepository struct {
	kubeFactory *KubeConfigFactory
	gvkFactory  *factory.GVKFactory
}

func NewDeploymentRepository(kubeFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) DeploymentRepository {
	return &deploymentRepository{
		kubeFactory: kubeFactory,
		gvkFactory:  gvkFactory,
	}
}

func (r *deploymentRepository) List(ctx context.Context, id int, ns string) (deployments []Deployment, err error) {
	var (
		deployment Deployment
		updatedAt  string
	)
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployments list failed. err: %v", err)
		return nil, err
	}
	for _, item := range response.Items {
		var matchLabels string
		// Selector 调度标签
		for k, v := range item.Spec.Selector.MatchLabels {
			matchLabels += fmt.Sprintf("%v=%v,", k, v) // 标签
		}
		// 更新时间
		for _, condition := range item.Status.Conditions {
			if condition.Type == appsv1.DeploymentProgressing && condition.Reason == "NewReplicaSetAvailable" {
				updatedAt = condition.LastUpdateTime.Time.Format("2006-01-02 15:04:05")
			}
		}

		deployment.Namespace = item.Namespace
		deployment.Name = item.Name
		deployment.Paused = item.Spec.Paused
		deployment.Strategy = item.Spec.Strategy
		deployment.MatchLabels = item.Spec.Selector.MatchLabels
		deployment.Replicas = item.Status.Replicas
		deployment.ReadyReplicas = item.Status.ReadyReplicas
		deployment.UpdatedReplicas = item.Status.UpdatedReplicas
		deployment.AvailableReplicas = item.Status.AvailableReplicas
		deployment.Conditions = item.Status.Conditions
		deployment.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		deployment.UpdatedAt = updatedAt
		deployments = append(deployments, deployment)
	}
	return deployments, nil
}

func (r *deploymentRepository) GetDetail(ctx context.Context, id int, ns, name string) (deployment Deployment, err error) {
	var updatedAt string
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return deployment, err
	}
	response, err := client.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment get failed. err: %v", err)
	}

	// 更新时间
	for _, condition := range response.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing && condition.Reason == "NewReplicaSetAvailable" {
			updatedAt = condition.LastUpdateTime.Time.Format("2006-01-02 15:04:05")
		}
	}

	deployment.Namespace = response.Namespace
	deployment.Name = response.Name
	deployment.Paused = response.Spec.Paused
	deployment.Strategy = response.Spec.Strategy
	deployment.MatchLabels = response.Spec.Selector.MatchLabels
	deployment.Replicas = response.Status.Replicas
	deployment.ReadyReplicas = response.Status.ReadyReplicas
	deployment.UpdatedReplicas = response.Status.UpdatedReplicas
	deployment.AvailableReplicas = response.Status.AvailableReplicas
	deployment.Conditions = response.Status.Conditions
	deployment.CreatedAt = response.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
	deployment.UpdatedAt = updatedAt
	return deployment, nil
}

func (r *deploymentRepository) Get(ctx context.Context, id int, ns, name string) (*appsv1.Deployment, error) {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *deploymentRepository) Create(ctx context.Context, id int, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().Deployments(ns).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *deploymentRepository) Update(ctx context.Context, id int, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().Deployments(ns).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

// Scale 更新副本数量
func (r *deploymentRepository) Scale(ctx context.Context, id int, ns, name string, replicas int32) (*autoscalingv1.Scale, error) {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	// 获取副本数量
	scale, err := client.AppsV1().Deployments(ns).GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment get scale failed. err: %v", err)
		return nil, err
	}
	scale.Spec.Replicas = replicas // 更新副本数量
	response, err := client.AppsV1().Deployments(ns).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment update scale failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *deploymentRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.AppsV1().Deployments(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment delete failed. err: %v", err)
	}
	return err
}

func (r *deploymentRepository) Restart(ctx context.Context, id int, ns, name string) (*appsv1.Deployment, error) {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	data := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
	response, err := client.AppsV1().Deployments(ns).Patch(ctx, name, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment restart failed. err: %v", err)

	}
	r.gvkFactory.Complete(response)
	return response, err
}

func (r *deploymentRepository) Rollout(ctx context.Context, id int, ns, name, rsName string) (*appsv1.Deployment, error) {
	client, err := r.kubeFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	// 1. 获取当前 deployment 版本
	deployment, err := client.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment get failed. err: %v", err)
		return nil, err
	}
	// 2. 获取当前 deployment 的 rs 版本(在 rs 方法中获取)，并通过 rs 名称找到指定版本
	replica, err := client.AppsV1().ReplicaSets(ns).Get(context.TODO(), rsName, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 replicaSet get failed. err: %v", err)
		return nil, err
	}
	// 3. 替换指定版本 template
	deployment.Spec.Template = replica.Spec.Template
	// 4. 执行 deployment 的 update 方法实现回滚
	response, err := client.AppsV1().Deployments(ns).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 deployment rollback update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}
