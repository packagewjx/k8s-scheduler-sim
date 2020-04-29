package mock

import (
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
)

// For Kubernetes scheduler use only. ONLY implements those used by scheduler.
type Client struct {
}

func NewClient() kubernetes.Interface {
	return &Client{}
}

func (client *Client) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	panic("implement me")
}

func (client *Client) DiscoveryV1alpha1() discoveryv1alpha1.DiscoveryV1alpha1Interface {
	panic("implement me")
}

func (client *Client) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	panic("implement me")
}

func (client *Client) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	panic("implement me")
}

func (client *Client) Discovery() discovery.DiscoveryInterface {
	panic("implement me")
}

func (client *Client) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	panic("implement me")
}

func (client *Client) AppsV1() appsv1.AppsV1Interface {
	panic("implement me")
}

func (client *Client) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	panic("implement me")
}

func (client *Client) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	panic("implement me")
}

func (client *Client) AuditregistrationV1alpha1() auditregistrationv1alpha1.AuditregistrationV1alpha1Interface {
	panic("implement me")
}

func (client *Client) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	panic("implement me")
}

func (client *Client) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	panic("implement me")
}

func (client *Client) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	panic("implement me")
}

func (client *Client) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	panic("implement me")
}

func (client *Client) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	panic("implement me")
}

func (client *Client) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	panic("implement me")
}

func (client *Client) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	panic("implement me")
}

func (client *Client) BatchV1() batchv1.BatchV1Interface {
	panic("implement me")
}

func (client *Client) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	panic("implement me")
}

func (client *Client) BatchV2alpha1() batchv2alpha1.BatchV2alpha1Interface {
	panic("implement me")
}

func (client *Client) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	panic("implement me")
}

func (client *Client) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	panic("implement me")
}

func (client *Client) CoordinationV1() coordinationv1.CoordinationV1Interface {
	panic("implement me")
}

func (client *Client) CoreV1() corev1.CoreV1Interface {
	panic("implement me")
}

func (client *Client) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	panic("implement me")
}

func (client *Client) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	panic("implement me")
}

func (client *Client) NetworkingV1() networkingv1.NetworkingV1Interface {
	panic("implement me")
}

func (client *Client) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	panic("implement me")
}

func (client *Client) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	panic("implement me")
}

func (client *Client) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	panic("implement me")
}

func (client *Client) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	panic("implement me")
}

func (client *Client) RbacV1() rbacv1.RbacV1Interface {
	panic("implement me")
}

func (client *Client) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	panic("implement me")
}

func (client *Client) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	panic("implement me")
}

func (client *Client) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	panic("implement me")
}

func (client *Client) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	panic("implement me")
}

func (client *Client) SchedulingV1() schedulingv1.SchedulingV1Interface {
	panic("implement me")
}

func (client *Client) SettingsV1alpha1() settingsv1alpha1.SettingsV1alpha1Interface {
	panic("implement me")
}

func (client *Client) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	panic("implement me")
}

func (client *Client) StorageV1() storagev1.StorageV1Interface {
	panic("implement me")
}

func (client *Client) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
	panic("implement me")
}
