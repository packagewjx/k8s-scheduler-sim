// fake 本包提供测试使用的一些简单的无其他依赖的接口实现，不会产生和修改实际数据
package fake

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
	"k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	apimachineryv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	deprecatedv1 "k8s.io/client-go/deprecated/typed/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	appsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	auditregistrationv1alpha1 "k8s.io/client-go/kubernetes/typed/auditregistration/v1alpha1"
	authenticationv1 "k8s.io/client-go/kubernetes/typed/authentication/v1"
	authenticationv1beta1 "k8s.io/client-go/kubernetes/typed/authentication/v1beta1"
	authorizationv1 "k8s.io/client-go/kubernetes/typed/authorization/v1"
	authorizationv1beta1 "k8s.io/client-go/kubernetes/typed/authorization/v1beta1"
	autoscalingv1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta2"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	batchv1beta1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
	batchv2alpha1 "k8s.io/client-go/kubernetes/typed/batch/v2alpha1"
	certificatesv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	coordinationv1beta1 "k8s.io/client-go/kubernetes/typed/coordination/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	discoveryv1alpha1 "k8s.io/client-go/kubernetes/typed/discovery/v1alpha1"
	discoveryv1beta1 "k8s.io/client-go/kubernetes/typed/discovery/v1beta1"
	eventsv1beta1 "k8s.io/client-go/kubernetes/typed/events/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	flowcontrolv1alpha1 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1alpha1"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	networkingv1beta1 "k8s.io/client-go/kubernetes/typed/networking/v1beta1"
	nodev1alpha1 "k8s.io/client-go/kubernetes/typed/node/v1alpha1"
	nodev1beta1 "k8s.io/client-go/kubernetes/typed/node/v1beta1"
	policyv1beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	rbacv1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	schedulingv1 "k8s.io/client-go/kubernetes/typed/scheduling/v1"
	schedulingv1alpha1 "k8s.io/client-go/kubernetes/typed/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/client-go/kubernetes/typed/scheduling/v1beta1"
	settingsv1alpha1 "k8s.io/client-go/kubernetes/typed/settings/v1alpha1"
	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	storagev1alpha1 "k8s.io/client-go/kubernetes/typed/storage/v1alpha1"
	storagev1beta1 "k8s.io/client-go/kubernetes/typed/storage/v1beta1"
	"k8s.io/client-go/rest"
)

func NewFakeKubernetesInterface() kubernetes.Interface {
	_ = util.GetMessageQueue().NewTopic(util.TopicPod)
	_ = util.GetMessageQueue().NewTopic(util.TopicNode)
	return &fakeKubernetesInterface{}
}

// 无需任何依赖的kubernetes.Interface，主要测试各个Informer的功能
type fakeKubernetesInterface struct {
}

func (f *fakeKubernetesInterface) RESTClient() rest.Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) ComponentStatuses() deprecatedv1.ComponentStatusInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) ConfigMaps(_ string) deprecatedv1.ConfigMapInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Endpoints(_ string) deprecatedv1.EndpointsInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Events(_ string) deprecatedv1.EventInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) LimitRanges(_ string) deprecatedv1.LimitRangeInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Namespaces() deprecatedv1.NamespaceInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Nodes() deprecatedv1.NodeInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) PersistentVolumes() deprecatedv1.PersistentVolumeInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) PersistentVolumeClaims(_ string) deprecatedv1.PersistentVolumeClaimInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Pods(_ string) deprecatedv1.PodInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) PodTemplates(_ string) deprecatedv1.PodTemplateInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) ReplicationControllers(_ string) deprecatedv1.ReplicationControllerInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) ResourceQuotas(_ string) deprecatedv1.ResourceQuotaInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Secrets(_ string) deprecatedv1.SecretInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Services(_ string) deprecatedv1.ServiceInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) ServiceAccounts(_ string) deprecatedv1.ServiceAccountInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) Discovery() discovery.DiscoveryInterface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AppsV1() appsv1.AppsV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AuditregistrationV1alpha1() auditregistrationv1alpha1.AuditregistrationV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) BatchV1() batchv1.BatchV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) BatchV2alpha1() batchv2alpha1.BatchV2alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) CoordinationV1() coordinationv1.CoordinationV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) CoreV1() corev1.CoreV1Interface {
	return &fakeCoreV1Interface{}
}

func (f *fakeKubernetesInterface) DiscoveryV1alpha1() discoveryv1alpha1.DiscoveryV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) NetworkingV1() networkingv1.NetworkingV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) RbacV1() rbacv1.RbacV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) SchedulingV1() schedulingv1.SchedulingV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) SettingsV1alpha1() settingsv1alpha1.SettingsV1alpha1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) StorageV1() storagev1.StorageV1Interface {
	panic("implement me")
}

func (f *fakeKubernetesInterface) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
	panic("implement me")
}

type fakeCoreV1Interface struct {
}

func (f *fakeCoreV1Interface) RESTClient() rest.Interface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) ComponentStatuses() corev1.ComponentStatusInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) ConfigMaps(_ string) corev1.ConfigMapInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Endpoints(_ string) corev1.EndpointsInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Events(_ string) corev1.EventInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) LimitRanges(_ string) corev1.LimitRangeInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Namespaces() corev1.NamespaceInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Nodes() corev1.NodeInterface {
	return &fakeNodeInterface{}
}

func (f *fakeCoreV1Interface) PersistentVolumes() corev1.PersistentVolumeInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) PersistentVolumeClaims(_ string) corev1.PersistentVolumeClaimInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Pods(_ string) corev1.PodInterface {
	return &fakePodInterface{}
}

func (f *fakeCoreV1Interface) PodTemplates(_ string) corev1.PodTemplateInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) ReplicationControllers(_ string) corev1.ReplicationControllerInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) ResourceQuotas(_ string) corev1.ResourceQuotaInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Secrets(_ string) corev1.SecretInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) Services(_ string) corev1.ServiceInterface {
	panic("implement me")
}

func (f *fakeCoreV1Interface) ServiceAccounts(_ string) corev1.ServiceAccountInterface {
	panic("implement me")
}

type fakeNodeInterface struct {
}

func (f *fakeNodeInterface) Create(_ context.Context, node *v1.Node, _ apimachineryv1.CreateOptions) (*v1.Node, error) {
	addEvent := &watch.Event{
		Type:   watch.Added,
		Object: node,
	}
	_ = util.GetMessageQueue().Publish(util.TopicNode, addEvent)
	return node, nil
}

func (f *fakeNodeInterface) Update(_ context.Context, node *v1.Node, _ apimachineryv1.UpdateOptions) (*v1.Node, error) {
	updateEvent := &watch.Event{
		Type:   watch.Modified,
		Object: node,
	}
	_ = util.GetMessageQueue().Publish(util.TopicNode, updateEvent)
	return node, nil
}

func (f *fakeNodeInterface) UpdateStatus(_ context.Context, node *v1.Node, _ apimachineryv1.UpdateOptions) (*v1.Node, error) {
	updateEvent := &watch.Event{
		Type:   watch.Modified,
		Object: node,
	}
	_ = util.GetMessageQueue().Publish(util.TopicNode, updateEvent)
	return node, nil
}

func (f *fakeNodeInterface) Delete(_ context.Context, _ string, _ apimachineryv1.DeleteOptions) error {
	deleteEvent := &watch.Event{
		Type:   watch.Deleted,
		Object: &v1.Node{},
	}
	_ = util.GetMessageQueue().Publish(util.TopicNode, deleteEvent)
	return nil
}

func (f *fakeNodeInterface) DeleteCollection(_ context.Context, _ apimachineryv1.DeleteOptions, _ apimachineryv1.ListOptions) error {
	panic("implement me")
}

func (f *fakeNodeInterface) Get(_ context.Context, _ string, _ apimachineryv1.GetOptions) (*v1.Node, error) {
	return &v1.Node{}, nil
}

func (f *fakeNodeInterface) List(_ context.Context, _ apimachineryv1.ListOptions) (*v1.NodeList, error) {
	nodes := []v1.Node{{}}
	list := &v1.NodeList{
		TypeMeta: apimachineryv1.TypeMeta{},
		ListMeta: apimachineryv1.ListMeta{},
		Items:    nodes,
	}
	return list, nil
}

func (f *fakeNodeInterface) Watch(_ context.Context, _ apimachineryv1.ListOptions) (watch.Interface, error) {
	return util.GetMessageQueue().Subscribe(util.TopicNode)
}

func (f *fakeNodeInterface) Patch(_ context.Context, _ string, _ types.PatchType, _ []byte, _ apimachineryv1.PatchOptions, _ ...string) (result *v1.Node, err error) {
	panic("implement me")
}

func (f *fakeNodeInterface) PatchStatus(_ context.Context, _ string, _ []byte) (*v1.Node, error) {
	panic("implement me")
}

type fakePodInterface struct {
}

func (f *fakePodInterface) Create(_ context.Context, pod *v1.Pod, _ apimachineryv1.CreateOptions) (*v1.Pod, error) {
	ev := &watch.Event{
		Type:   watch.Added,
		Object: pod,
	}
	_ = util.GetMessageQueue().Publish(util.TopicPod, ev)
	return pod, nil
}

func (f *fakePodInterface) Update(_ context.Context, pod *v1.Pod, _ apimachineryv1.UpdateOptions) (*v1.Pod, error) {
	ev := &watch.Event{
		Type:   watch.Modified,
		Object: pod,
	}
	_ = util.GetMessageQueue().Publish(util.TopicPod, ev)
	return pod, nil
}

func (f *fakePodInterface) UpdateStatus(ctx context.Context, pod *v1.Pod, opts apimachineryv1.UpdateOptions) (*v1.Pod, error) {
	return f.Update(ctx, pod, opts)
}

func (f *fakePodInterface) Delete(_ context.Context, name string, _ apimachineryv1.DeleteOptions) error {
	ev := &watch.Event{
		Type: watch.Deleted,
		Object: &v1.Pod{
			ObjectMeta: apimachineryv1.ObjectMeta{
				Name: name,
			},
		},
	}
	_ = util.GetMessageQueue().Publish(util.TopicPod, ev)
	return nil
}

func (f *fakePodInterface) DeleteCollection(_ context.Context, _ apimachineryv1.DeleteOptions, _ apimachineryv1.ListOptions) error {
	panic("implement me")
}

func (f *fakePodInterface) Get(_ context.Context, _ string, _ apimachineryv1.GetOptions) (*v1.Pod, error) {
	panic("implement me")
}

func (f *fakePodInterface) List(_ context.Context, _ apimachineryv1.ListOptions) (*v1.PodList, error) {
	panic("implement me")
}

func (f *fakePodInterface) Watch(_ context.Context, _ apimachineryv1.ListOptions) (watch.Interface, error) {
	return util.GetMessageQueue().Subscribe(util.TopicPod)
}

func (f *fakePodInterface) Patch(_ context.Context, _ string, _ types.PatchType, _ []byte, _ apimachineryv1.PatchOptions, _ ...string) (result *v1.Pod, err error) {
	panic("implement me")
}

func (f *fakePodInterface) GetEphemeralContainers(_ context.Context, _ string, _ apimachineryv1.GetOptions) (*v1.EphemeralContainers, error) {
	panic("implement me")
}

func (f *fakePodInterface) UpdateEphemeralContainers(_ context.Context, _ string, _ *v1.EphemeralContainers, _ apimachineryv1.UpdateOptions) (*v1.EphemeralContainers, error) {
	panic("implement me")
}

func (f *fakePodInterface) Bind(_ context.Context, _ *v1.Binding, _ apimachineryv1.CreateOptions) error {
	panic("implement me")
}

func (f *fakePodInterface) Evict(_ context.Context, _ *v1beta1.Eviction) error {
	panic("implement me")
}

func (f *fakePodInterface) GetLogs(_ string, _ *v1.PodLogOptions) *rest.Request {
	panic("implement me")
}
