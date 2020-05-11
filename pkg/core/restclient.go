package core

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
)

type restClient struct {
}

func (rc *restClient) GetRateLimiter() flowcontrol.RateLimiter {
	return nil
}

func (rc *restClient) Verb(verb string) *rest.Request {
	panic("implement me")
}

func (rc *restClient) Post() *rest.Request {
	return rc.Verb("Post")
}

func (rc *restClient) Put() *rest.Request {
	return rc.Verb("Put")
}

func (rc *restClient) Patch(pt types.PatchType) *rest.Request {
	return rc.Verb("Patch")
}

func (rc *restClient) Get() *rest.Request {
	return rc.Verb("Get")
}

func (rc *restClient) Delete() *rest.Request {
	return rc.Verb("Delete")
}

func (rc *restClient) APIVersion() schema.GroupVersion {
	return schema.GroupVersion{
		Group:   "core",
		Version: "v1",
	}
}
