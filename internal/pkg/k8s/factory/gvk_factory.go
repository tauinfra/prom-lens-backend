package factory

import (
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	eventsv1 "k8s.io/api/events/v1"
	netv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// =======================================================
// GVKFactory —— 工厂模式
// 提供 GVK 补全、Scheme 管理、CRD 注册
// =======================================================

// GVKFactory 负责管理 scheme，并提供统一的 GVK 完成工具
type GVKFactory struct {
	scheme *runtime.Scheme
	once   sync.Once
}

// NewGVKFactory 工厂构造方法
func NewGVKFactory() *GVKFactory {
	return &GVKFactory{}
}

// initScheme 初始化 Kubernetes Scheme（只初始化一次）
func (f *GVKFactory) initScheme() {
	f.once.Do(func() {
		f.scheme = runtime.NewScheme()
		// Kubernetes API group 注册
		_ = corev1.AddToScheme(f.scheme)
		_ = appsv1.AddToScheme(f.scheme)
		_ = rbacv1.AddToScheme(f.scheme)
		_ = storagev1.AddToScheme(f.scheme)
		_ = batchv1.AddToScheme(f.scheme)
		_ = netv1.AddToScheme(f.scheme)
		_ = policyv1.AddToScheme(f.scheme)
		_ = eventsv1.AddToScheme(f.scheme)
		_ = coordinationv1.AddToScheme(f.scheme)
		// Meta（List / Status 等）
		_ = metav1.AddMetaToScheme(f.scheme)
	})
}

// RegisterCRD 支持向工厂注入任何 CRD（扩展性强）
func (f *GVKFactory) RegisterCRD(addToScheme func(*runtime.Scheme) error) error {
	f.initScheme()
	return addToScheme(f.scheme)
}

// Complete 自动补齐 Kind、APIVersion
func (f *GVKFactory) Complete(obj runtime.Object) {
	f.initScheme()
	gvks, _, err := f.scheme.ObjectKinds(obj)
	if err != nil || len(gvks) == 0 {
		return
	}

	gvk := gvks[0]

	obj.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	})
}
