package watcher

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type Client struct {
	discoveryClient *discovery.DiscoveryClient
	discoveryMapper *restmapper.DeferredDiscoveryRESTMapper
	dynamicClient   dynamic.Interface
}

func newClient(config *rest.Config) (*Client, error) {

	disC, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	dynC, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cacheC := memory.NewMemCacheClient(disC)
	cacheC.Invalidate()

	dm := restmapper.NewDeferredDiscoveryRESTMapper(cacheC)

	return &Client{disC, dm, dynC}, nil
}
