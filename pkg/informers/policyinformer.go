package informers

import (
	"k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers/policy/v1beta1"
	beta1 "k8s.io/client-go/listers/policy/v1beta1"
	"k8s.io/client-go/tools/cache"
)

type policyInformer struct {
}

func (p *policyInformer) PodDisruptionBudgets() v1beta1.PodDisruptionBudgetInformer {
	return &podDisruptionBudgetInformer{}
}

func (p *policyInformer) PodSecurityPolicies() v1beta1.PodSecurityPolicyInformer {
	panic("implement me")
}

func (p *policyInformer) V1beta1() v1beta1.Interface {
	return p
}

type podDisruptionBudgetInformer struct {
}

func (p *podDisruptionBudgetInformer) List(selector labels.Selector) (ret []*v1beta12.PodDisruptionBudget, err error) {
	panic("implement me")
}

func (p *podDisruptionBudgetInformer) PodDisruptionBudgets(namespace string) beta1.PodDisruptionBudgetNamespaceLister {
	panic("implement me")
}

func (p *podDisruptionBudgetInformer) GetPodPodDisruptionBudgets(pod *v1.Pod) ([]*v1beta12.PodDisruptionBudget, error) {
	panic("implement me")
}

func (p *podDisruptionBudgetInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (p *podDisruptionBudgetInformer) Lister() beta1.PodDisruptionBudgetLister {
	return p
}
