package pkg

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/simulate"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	"strconv"
)

var (
	TopicNode = "node"
	TopicPod  = "pod"
)

// For Kubernetes scheduler use only. ONLY implements those used by scheduler.
// All options are ignored, for simplicity.
type simClient struct {
	sim *SchedSim
}

func (client *simClient) RESTClient() rest.Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) ComponentStatuses() deprecatedv1.ComponentStatusInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) ConfigMaps(_ string) deprecatedv1.ConfigMapInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Endpoints(_ string) deprecatedv1.EndpointsInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Events(_ string) deprecatedv1.EventInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) LimitRanges(_ string) deprecatedv1.LimitRangeInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Namespaces() deprecatedv1.NamespaceInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Nodes() deprecatedv1.NodeInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) PersistentVolumes() deprecatedv1.PersistentVolumeInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) PersistentVolumeClaims(_ string) deprecatedv1.PersistentVolumeClaimInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Pods(_ string) deprecatedv1.PodInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) PodTemplates(_ string) deprecatedv1.PodTemplateInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) ReplicationControllers(_ string) deprecatedv1.ReplicationControllerInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) ResourceQuotas(_ string) deprecatedv1.ResourceQuotaInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Secrets(_ string) deprecatedv1.SecretInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Services(_ string) deprecatedv1.ServiceInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) ServiceAccounts(_ string) deprecatedv1.ServiceAccountInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) DiscoveryV1alpha1() discoveryv1alpha1.DiscoveryV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) Discovery() discovery.DiscoveryInterface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AppsV1() appsv1.AppsV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AuditregistrationV1alpha1() auditregistrationv1alpha1.AuditregistrationV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) BatchV1() batchv1.BatchV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) BatchV2alpha1() batchv2alpha1.BatchV2alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) CoordinationV1() coordinationv1.CoordinationV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) CoreV1() corev1.CoreV1Interface {
	return &coreV1Client{sim: client.sim}
}

func (client *simClient) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) NetworkingV1() networkingv1.NetworkingV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) RbacV1() rbacv1.RbacV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) SchedulingV1() schedulingv1.SchedulingV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) SettingsV1alpha1() settingsv1alpha1.SettingsV1alpha1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) StorageV1() storagev1.StorageV1Interface {
	panic("Using this interface is not allowed.")
}

func (client *simClient) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
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
	panic("Using this interface is not allowed.")
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

func (client *coreV1NodeClient) Create(_ context.Context, node *apicorev1.Node, _ apimachineryv1.CreateOptions) (*apicorev1.Node, error) {
	// 创建CoreScheduler
	schedulerName := node.Annotations[simulate.NodeAnnotationCoreScheduler]
	scheduler, exist := simulate.GetCoreScheduler(schedulerName)
	if !exist {
		return nil, fmt.Errorf("No CoreScheduler %s", schedulerName)
	}

	numCpu, ok := node.Status.Capacity.Cpu().AsInt64()
	if !ok || numCpu == 0 {
		return nil, fmt.Errorf("cpu num must larger than 0")
	}

	clone := node.DeepCopy()
	simNode := &simulate.Node{
		Node:         *clone,
		Scheduler:    scheduler,
		Pods:         map[string]*simulate.Pod{},
		CpuState:     make([][]*simulate.RunEntity, numCpu),
		LastCpuUsage: 0,
	}

	err := client.sim.Nodes.Add(simNode)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error adding node %s to store", node.Name))
	}

	ev := &watch.Event{
		Type:   watch.Added,
		Object: node,
	}
	err = util.GetMessageQueue().Publish(TopicNode, ev)
	if err != nil {
		logrus.Errorf("Error publishing add event: %v", err)
	}

	return node, nil
}

func (client *coreV1NodeClient) Update(_ context.Context, node *apicorev1.Node, _ apimachineryv1.UpdateOptions) (*apicorev1.Node, error) {
	item, exists, err := client.sim.Nodes.Get(node)
	if !exists {
		return nil, fmt.Errorf("no node name %s", node.Name)
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error getting oldNode %s", node.Name))
	}

	storeNode := item.(*simulate.Node)
	storeNode.Node = *(node.DeepCopy())

	// 这里暂时没有更改simulate.Node的额外属性

	err = client.sim.Nodes.Update(storeNode)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error updating node %s", node.Name))
	}

	ev := &watch.Event{
		Type:   watch.Modified,
		Object: node,
	}
	err = util.GetMessageQueue().Publish(TopicNode, ev)
	if err != nil {
		logrus.Errorf("error publishing update event %v", err)
	}

	return &storeNode.Node, nil
}

func (client *coreV1NodeClient) UpdateStatus(_ context.Context, node *apicorev1.Node, _ apimachineryv1.UpdateOptions) (*apicorev1.Node, error) {
	item, exists, err := client.sim.Nodes.Get(node)
	if !exists {
		return nil, fmt.Errorf("no node name %s", node.Name)
	}
	if err != nil {
		return nil, err
	}

	storeNode := item.(*simulate.Node)
	storeNode.Status = *(node.Status.DeepCopy())

	// 这里暂时没有更改simulate.Node的额外属性

	err = client.sim.Nodes.Update(storeNode)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error updating node %s", node.Name))
	}

	ev := &watch.Event{
		Type:   watch.Modified,
		Object: node,
	}
	err = util.GetMessageQueue().Publish(TopicNode, ev)
	if err != nil {
		logrus.Errorf("error publishing update event: %v", err)
	}

	return &storeNode.Node, nil
}

func (client *coreV1NodeClient) Delete(_ context.Context, name string, _ apimachineryv1.DeleteOptions) error {
	item, exists, err := client.sim.Nodes.GetByKey(name)
	if !exists {
		return fmt.Errorf("no node %s", name)
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error getting node %s", name))
	}

	err = client.sim.Nodes.Delete(item)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error deleting node %s", name))
	}

	// 发送删除通知
	ev := &watch.Event{
		Type:   watch.Deleted,
		Object: &item.(*simulate.Node).Node,
	}
	err = util.GetMessageQueue().Publish(TopicNode, ev)
	if err != nil {
		logrus.Errorf("Error publishing delete event: %v", err)
	}

	return nil
}

func (client *coreV1NodeClient) DeleteCollection(_ context.Context, _ apimachineryv1.DeleteOptions, _ apimachineryv1.ListOptions) error {
	panic("Using this interface is not allowed.")
}

func (client *coreV1NodeClient) Get(_ context.Context, name string, _ apimachineryv1.GetOptions) (*apicorev1.Node, error) {
	item, exists, err := client.sim.Nodes.GetByKey(name)
	if !exists {
		return nil, fmt.Errorf("no node name %s", name)
	}
	if err != nil {
		return nil, err
	}
	node := item.(*simulate.Node)
	return &node.Node, nil
}

func (client *coreV1NodeClient) List(_ context.Context, _ apimachineryv1.ListOptions) (*apicorev1.NodeList, error) {
	zero := int64(0)

	nodes := make([]apicorev1.Node, 0, 10)
	list := client.sim.Nodes.List()
	for _, node := range list {
		nodes = append(nodes, node.(*simulate.Node).Node)
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

func (c *coreV1PodClient) Create(_ context.Context, pod *apicorev1.Pod, _ apimachineryv1.CreateOptions) (*apicorev1.Pod, error) {
	cpuLimit, err := strconv.ParseFloat(pod.Annotations[simulate.PodAnnotationCpuLimit], 64)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing cpulimit")
	}

	memLimit, err := strconv.ParseInt(pod.Annotations[simulate.PodAnnotationMemLimit], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing memlimit")
	}

	algName, ok := pod.Annotations[simulate.PodAnnotationAlgorithm]
	if !ok {
		return nil, fmt.Errorf("pod must have algorithm to run")
	}

	factory, exist := simulate.GetPodAlgorithmFactory(algName)
	if !exist {
		return nil, fmt.Errorf("no pod algorithm %s", algName)
	}

	stateString, _ := pod.Annotations[simulate.PodAnnotationInitialState]

	_, ok = pod.Annotations[simulate.PodAnnotationDeploymentController]
	if !ok {
		return nil, fmt.Errorf("pod must have deployment controller name")
	}

	clone := pod.DeepCopy()
	simPod := &simulate.Pod{
		Pod:       *clone,
		CpuLimit:  cpuLimit,
		MemLimit:  memLimit,
		Algorithm: nil,
	}
	algorithm, err := factory(stateString, simPod)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating pod algorithm")
	}
	simPod.Algorithm = algorithm

	err = c.sim.Pods.Add(simPod)
	if err != nil {
		return nil, errors.Wrap(err, "Error adding to store")
	}

	// 发送添加通知
	ev := &watch.Event{
		Type:   watch.Added,
		Object: pod,
	}
	err = util.GetMessageQueue().Publish(TopicPod, ev)
	if err != nil {
		logrus.Errorf("Error publishing add event: %v", err)
	}

	return pod, nil
}

func (c *coreV1PodClient) Update(_ context.Context, pod *apicorev1.Pod, _ apimachineryv1.UpdateOptions) (*apicorev1.Pod, error) {
	item, exists, err := c.sim.Pods.Get(pod)
	if !exists {
		return nil, fmt.Errorf("no pod %s", pod.Name)
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error getting pod %s", pod.Name))
	}

	simPod := item.(*simulate.Pod)
	simPod.Pod = *(pod.DeepCopy())

	err = c.sim.Pods.Update(simPod)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error updating pod %s", pod.Name))
	}

	ev := &watch.Event{
		Type:   watch.Modified,
		Object: pod,
	}
	err = util.GetMessageQueue().Publish(TopicPod, ev)
	if err != nil {
		logrus.Errorf("Error publishing update event: %v", err)
	}

	return pod, nil
}

func (c *coreV1PodClient) UpdateStatus(_ context.Context, pod *apicorev1.Pod, _ apimachineryv1.UpdateOptions) (*apicorev1.Pod, error) {
	item, exists, err := c.sim.Pods.Get(pod)
	if !exists {
		return nil, fmt.Errorf("no pod %s", pod.Name)
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error getting pod %s", pod.Name))
	}

	simPod := item.(*simulate.Pod)
	simPod.Pod.Status = *(pod.Status.DeepCopy())

	err = c.sim.Pods.Update(simPod)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error updating status of pod %s", pod.Name))
	}

	ev := &watch.Event{
		Type:   watch.Modified,
		Object: pod,
	}
	err = util.GetMessageQueue().Publish(TopicPod, ev)
	if err != nil {
		logrus.Errorf("Error publishing update event: %v", err)
	}

	return pod, nil
}

func (c *coreV1PodClient) Delete(_ context.Context, name string, _ apimachineryv1.DeleteOptions) error {
	item, exists, err := c.sim.Pods.GetByKey(name)
	if !exists {
		return fmt.Errorf("no pod %s", name)
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error getting pod %s", name))
	}

	err = c.sim.Pods.Delete(item)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error deleting pod %s", name))
	}

	ev := &watch.Event{
		Type:   watch.Deleted,
		Object: &item.(*simulate.Pod).Pod,
	}
	err = util.GetMessageQueue().Publish(TopicPod, ev)
	if err != nil {
		logrus.Errorf("Error publishing delete event: %v", err)
	}

	return nil
}

func (c *coreV1PodClient) DeleteCollection(_ context.Context, _ apimachineryv1.DeleteOptions, _ apimachineryv1.ListOptions) error {
	panic("Using this interface is not allowed.")
}

func (c *coreV1PodClient) Get(_ context.Context, name string, _ apimachineryv1.GetOptions) (*apicorev1.Pod, error) {
	item, exists, err := c.sim.Pods.GetByKey(name)
	if !exists {
		return nil, fmt.Errorf("no pod %s", name)
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error getting pod %s", name))
	}

	return &item.(*simulate.Pod).Pod, nil
}

func (c *coreV1PodClient) List(_ context.Context, _ apimachineryv1.ListOptions) (*apicorev1.PodList, error) {
	zero := int64(0)
	list := c.sim.Pods.List()

	arr := make([]apicorev1.Pod, 0, 10)
	for _, pod := range list {
		arr = append(arr, pod.(*simulate.Pod).Pod)
	}

	podList := &apicorev1.PodList{
		TypeMeta: apimachineryv1.TypeMeta{},
		ListMeta: apimachineryv1.ListMeta{
			ResourceVersion:    "",
			RemainingItemCount: &zero,
		},
		Items: arr,
	}
	return podList, nil
}

func (c *coreV1PodClient) Watch(_ context.Context, _ apimachineryv1.ListOptions) (watch.Interface, error) {
	return util.GetMessageQueue().Subscribe(TopicPod)
}

func (c *coreV1PodClient) Patch(_ context.Context, _ string, _ types.PatchType, _ []byte, _ apimachineryv1.PatchOptions, _ ...string) (result *apicorev1.Pod, err error) {
	panic("Using this interface is not allowed.")
}

func (c *coreV1PodClient) GetEphemeralContainers(_ context.Context, _ string, _ apimachineryv1.GetOptions) (*apicorev1.EphemeralContainers, error) {
	panic("Using this interface is not allowed.")
}

func (c *coreV1PodClient) UpdateEphemeralContainers(_ context.Context, _ string, _ *apicorev1.EphemeralContainers, _ apimachineryv1.UpdateOptions) (*apicorev1.EphemeralContainers, error) {
	panic("Using this interface is not allowed.")
}

func (c *coreV1PodClient) Bind(_ context.Context, binding *apicorev1.Binding, _ apimachineryv1.CreateOptions) error {
	fmt.Println(binding)
	return nil
}

func (c *coreV1PodClient) Evict(_ context.Context, _ *v1beta1.Eviction) error {
	panic("implement me")
}

func (c *coreV1PodClient) GetLogs(_ string, _ *apicorev1.PodLogOptions) *rest.Request {
	panic("Using this interface is not allowed.")
}
