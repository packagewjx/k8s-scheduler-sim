module github.com/packagewjx/k8s-scheduler-sim

go 1.14

require (
    k8s.io/kubernetes v1.18.0
    k8s.io/api v0.0.0
    k8s.io/apimachinery v0.0.0
    k8s.io/client-go v0.0.0
)

replace (
	k8s.io/api => ../../../k8s.io/kubernetes/staging/src/k8s.io/api
	k8s.io/apiextensions-apiserver => ../../../k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver
	k8s.io/apimachinery => ../../../k8s.io/kubernetes/staging/src/k8s.io/apimachinery
	k8s.io/apiserver => ../../../k8s.io/kubernetes/staging/src/k8s.io/apiserver
	k8s.io/cli-runtime => ../../../k8s.io/kubernetes/staging/src/k8s.io/cli-runtime
	k8s.io/client-go => ../../../k8s.io/kubernetes/staging/src/k8s.io/client-go
	k8s.io/cloud-provider => ../../../k8s.io/kubernetes/staging/src/k8s.io/cloud-provider
	k8s.io/cluster-bootstrap => ../../../k8s.io/kubernetes/staging/src/k8s.io/cluster-bootstrap
	k8s.io/code-generator => ../../../k8s.io/kubernetes/staging/src/k8s.io/code-generator
	k8s.io/component-base => ../../../k8s.io/kubernetes/staging/src/k8s.io/component-base
	k8s.io/cri-api => ../../../k8s.io/kubernetes/staging/src/k8s.io/cri-api
	k8s.io/csi-translation-lib => ../../../k8s.io/kubernetes/staging/src/k8s.io/csi-translation-lib
	k8s.io/gengo => k8s.io/gengo v0.0.0-20200114144118-36b2048a9120
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => ../../../k8s.io/kubernetes/staging/src/k8s.io/kube-aggregator
	k8s.io/kube-controller-manager => ../../../k8s.io/kubernetes/staging/src/k8s.io/kube-controller-manager
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	k8s.io/kube-proxy => ../../../k8s.io/kubernetes/staging/src/k8s.io/kube-proxy
	k8s.io/kube-scheduler => ../../../k8s.io/kubernetes/staging/src/k8s.io/kube-scheduler
	k8s.io/kubectl => ../../../k8s.io/kubernetes/staging/src/k8s.io/kubectl
	k8s.io/kubelet => ../../../k8s.io/kubernetes/staging/src/k8s.io/kubelet
	k8s.io/legacy-cloud-providers => ../../../k8s.io/kubernetes/staging/src/k8s.io/legacy-cloud-providers
	k8s.io/metrics => ../../../k8s.io/kubernetes/staging/src/k8s.io/metrics
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.1-alpha.1
	k8s.io/sample-apiserver => ../../../k8s.io/kubernetes/staging/src/k8s.io/sample-apiserver
	k8s.io/system-validators => k8s.io/system-validators v1.0.4
	k8s.io/utils => k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
)
