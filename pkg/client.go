package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	apimachineryv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	deprecatedv1 "k8s.io/client-go/deprecated/typed/core/v1"
	"k8s.io/client-go/discovery"
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

var (
	TopicNode = "node"
	TopicPod  = "pod"
)

// For Kubernetes scheduler use only. ONLY implements those used by scheduler.
type SimClient struct {
	Sim *SchedSim
}

func (client *SimClient) RESTClient() rest.Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) ComponentStatuses() deprecatedv1.ComponentStatusInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) ConfigMaps(_ string) deprecatedv1.ConfigMapInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Endpoints(_ string) deprecatedv1.EndpointsInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Events(_ string) deprecatedv1.EventInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) LimitRanges(_ string) deprecatedv1.LimitRangeInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Namespaces() deprecatedv1.NamespaceInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Nodes() deprecatedv1.NodeInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) PersistentVolumes() deprecatedv1.PersistentVolumeInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) PersistentVolumeClaims(_ string) deprecatedv1.PersistentVolumeClaimInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Pods(_ string) deprecatedv1.PodInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) PodTemplates(_ string) deprecatedv1.PodTemplateInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) ReplicationControllers(_ string) deprecatedv1.ReplicationControllerInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) ResourceQuotas(_ string) deprecatedv1.ResourceQuotaInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Secrets(_ string) deprecatedv1.SecretInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Services(_ string) deprecatedv1.ServiceInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) ServiceAccounts(_ string) deprecatedv1.ServiceAccountInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) DiscoveryV1alpha1() discoveryv1alpha1.DiscoveryV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) Discovery() discovery.DiscoveryInterface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AppsV1() appsv1.AppsV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AuditregistrationV1alpha1() auditregistrationv1alpha1.AuditregistrationV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) BatchV1() batchv1.BatchV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) BatchV2alpha1() batchv2alpha1.BatchV2alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) CoordinationV1() coordinationv1.CoordinationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) CoreV1() corev1.CoreV1Interface {
	return &coreV1Client{sim: client.Sim}
}

func (client *SimClient) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) NetworkingV1() networkingv1.NetworkingV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) RbacV1() rbacv1.RbacV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) SchedulingV1() schedulingv1.SchedulingV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) SettingsV1alpha1() settingsv1alpha1.SettingsV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) StorageV1() storagev1.StorageV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *SimClient) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

// coreV1Client 实现k8s.io/client-go/kubernetes/typed/core/v1.Interface
type coreV1Client struct {
	sim *SchedSim
}

func (client *coreV1Client) RESTClient() rest.Interface {
	return &restClient{}
}

func (client *coreV1Client) ComponentStatuses() corev1.ComponentStatusInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) ConfigMaps(_ string) corev1.ConfigMapInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) Endpoints(_ string) corev1.EndpointsInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) Events(_ string) corev1.EventInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) LimitRanges(_ string) corev1.LimitRangeInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) Namespaces() corev1.NamespaceInterface {
	panic("implement me")
}

func (client *coreV1Client) Nodes() corev1.NodeInterface {
	return &coreV1NodeClient{sim: client.sim}
}

func (client *coreV1Client) PersistentVolumes() corev1.PersistentVolumeInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) PersistentVolumeClaims(_ string) corev1.PersistentVolumeClaimInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) Pods(_ string) corev1.PodInterface {
	return &coreV1PodClient{sim: client.sim}
}

func (client *coreV1Client) PodTemplates(_ string) corev1.PodTemplateInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) ReplicationControllers(_ string) corev1.ReplicationControllerInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) ResourceQuotas(_ string) corev1.ResourceQuotaInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) Secrets(_ string) corev1.SecretInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) Services(_ string) corev1.ServiceInterface {
	panic("Using this interface is not allowed.")
}

func (client *coreV1Client) ServiceAccounts(_ string) corev1.ServiceAccountInterface {
	panic("Using this interface is not allowed.")
}

// coreV1NodeClient 实现corev1.NodeInterface
type coreV1NodeClient struct {
	sim *SchedSim
}

func (client *coreV1NodeClient) Create(_ context.Context, _ *apicorev1.Node, _ apimachineryv1.CreateOptions) (*apicorev1.Node, error) {
	panic("Using this interface is not allowed.")
}

func (client *coreV1NodeClient) Update(_ context.Context, node *apicorev1.Node, _ apimachineryv1.UpdateOptions) (*apicorev1.Node, error) {
	for _, simNode := range client.sim.Nodes {
		if simNode.Name == node.Name {
			simNode.Node = *node

			// 事件通知
			event := &watch.Event{
				Type:   watch.Modified,
				Object: node,
			}
			err := util.GetMessageQueue().Publish(TopicNode, event)
			if err != nil {
				fmt.Println(err)
			}

			return &simNode.Node, nil
		}
	}
	return nil, errors.New("No node named " + node.Name)
}

func (client *coreV1NodeClient) UpdateStatus(_ context.Context, node *apicorev1.Node, _ apimachineryv1.UpdateOptions) (*apicorev1.Node, error) {
	for _, simNode := range client.sim.Nodes {
		if simNode.Name == node.Name {
			simNode.Node.Status = node.Status

			// 事件通知
			event := &watch.Event{
				Type:   watch.Modified,
				Object: node,
			}
			err := util.GetMessageQueue().Publish(TopicNode, event)
			if err != nil {
				fmt.Println(err)
			}

			return &simNode.Node, nil
		}
	}
	return nil, errors.New("No node named " + node.Name)
}

func (client *coreV1NodeClient) Delete(_ context.Context, _ string, _ apimachineryv1.DeleteOptions) error {
	panic("Using this interface is not allowed.")
}

func (client *coreV1NodeClient) DeleteCollection(_ context.Context, _ apimachineryv1.DeleteOptions, _ apimachineryv1.ListOptions) error {
	panic("Using this interface is not allowed.")
}

func (client *coreV1NodeClient) Get(_ context.Context, name string, _ apimachineryv1.GetOptions) (*apicorev1.Node, error) {
	for _, simNode := range client.sim.Nodes {
		if simNode.Name == name {
			return &simNode.Node, nil
		}
	}
	return nil, errors.New("No node named " + name)
}

func (client *coreV1NodeClient) List(_ context.Context, _ apimachineryv1.ListOptions) (*apicorev1.NodeList, error) {
	zero := int64(0)

	nodes := make([]apicorev1.Node, len(client.sim.Nodes))
	for i := 0; i < len(client.sim.Nodes); i++ {
		nodes[i] = client.sim.Nodes[i].Node
	}

	return &apicorev1.NodeList{
		TypeMeta: apimachineryv1.TypeMeta{},
		ListMeta: apimachineryv1.ListMeta{
			ResourceVersion:    "v1",
			RemainingItemCount: &zero,
		},
		Items: nodes,
	}, nil
}

func (client *coreV1NodeClient) Watch(_ context.Context, _ apimachineryv1.ListOptions) (watch.Interface, error) {
	return util.GetMessageQueue().Subscribe(TopicNode)
}

func (client *coreV1NodeClient) Patch(_ context.Context, _ string, _ types.PatchType, _ []byte, _ apimachineryv1.PatchOptions, _ ...string) (result *apicorev1.Node, err error) {
	panic("Using this interface is not allowed.")
}

func (client *coreV1NodeClient) PatchStatus(_ context.Context, _ string, _ []byte) (*apicorev1.Node, error) {
	panic("Using this interface is not allowed.")
}

type coreV1PodClient struct {
	sim *SchedSim
}

func (c *coreV1PodClient) Create(_ context.Context, _ *apicorev1.Pod, _ apimachineryv1.CreateOptions) (*apicorev1.Pod, error) {
	panic("implement me")
}

func (c *coreV1PodClient) Update(_ context.Context, _ *apicorev1.Pod, _ apimachineryv1.UpdateOptions) (*apicorev1.Pod, error) {
	panic("implement me")
}

func (c *coreV1PodClient) UpdateStatus(_ context.Context, _ *apicorev1.Pod, _ apimachineryv1.UpdateOptions) (*apicorev1.Pod, error) {
	panic("implement me")
}

func (c *coreV1PodClient) Delete(_ context.Context, _ string, _ apimachineryv1.DeleteOptions) error {
	panic("implement me")
}

func (c *coreV1PodClient) DeleteCollection(_ context.Context, _ apimachineryv1.DeleteOptions, _ apimachineryv1.ListOptions) error {
	panic("implement me")
}

func (c *coreV1PodClient) Get(_ context.Context, _ string, _ apimachineryv1.GetOptions) (*apicorev1.Pod, error) {
	panic("implement me")
}

func (c *coreV1PodClient) List(_ context.Context, _ apimachineryv1.ListOptions) (*apicorev1.PodList, error) {
	panic("implement me")
}

func (c *coreV1PodClient) Watch(_ context.Context, _ apimachineryv1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (c *coreV1PodClient) Patch(_ context.Context, _ string, _ types.PatchType, _ []byte, _ apimachineryv1.PatchOptions, _ ...string) (result *apicorev1.Pod, err error) {
	panic("implement me")
}

func (c *coreV1PodClient) GetEphemeralContainers(_ context.Context, _ string, _ apimachineryv1.GetOptions) (*apicorev1.EphemeralContainers, error) {
	panic("implement me")
}

func (c *coreV1PodClient) UpdateEphemeralContainers(_ context.Context, _ string, _ *apicorev1.EphemeralContainers, _ apimachineryv1.UpdateOptions) (*apicorev1.EphemeralContainers, error) {
	panic("implement me")
}

func (c *coreV1PodClient) Bind(_ context.Context, _ *apicorev1.Binding, _ apimachineryv1.CreateOptions) error {
	panic("implement me")
}

func (c *coreV1PodClient) Evict(_ context.Context, _ *v1beta1.Eviction) error {
	panic("implement me")
}

func (c *coreV1PodClient) GetLogs(_ string, _ *apicorev1.PodLogOptions) *rest.Request {
	panic("implement me")
}
