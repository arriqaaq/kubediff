package watcher

import (
	"time"

	"github.com/arriqaaq/kdiff/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Informer interface {
	AddEventHandler(handler eventHandler)
	HasSynced() bool
	Start(ch <-chan struct{})
}

type NewInformerFunc func(config *rest.Config) (*multiResourceInformer, error)

func NewMultiResourceInformer(cfg *config.Config, client *Client, resyncPeriod time.Duration) NewInformerFunc {
	return func(kubeConfig *rest.Config) (*multiResourceInformer, error) {
		informers := make(map[string]map[string]cache.SharedIndexInformer)
		dynamicClient, err := dynamic.NewForConfig(kubeConfig)
		if err != nil {
			return nil, err
		}

		resources := make(map[string]schema.GroupVersionResource)
		for _, r := range cfg.Resources {
			gvr, err := getGVRFromResource(client.DiscoveryMapper, r.Kind)
			if err != nil {
				return nil, err
			}
			resources[r.Kind] = gvr
		}

		dynamicInformers := make([]dynamicinformer.DynamicSharedInformerFactory, 0, len(cfg.Namespaces))

		for _, ns := range cfg.Namespaces {

			namespace := getNamespace(ns)
			di := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
				dynamicClient,
				resyncPeriod,
				namespace,
				nil,
			)

			for r, gvr := range resources {
				if _, ok := informers[namespace]; !ok {
					informers[namespace] = make(map[string]cache.SharedIndexInformer)
				}
				informers[ns][r] = di.ForResource(gvr).Informer()
			}

			dynamicInformers = append(dynamicInformers, di)
		}

		return &multiResourceInformer{
			resourceToGVR:      resources,
			resourceToInformer: informers,
			informerFactory:    dynamicInformers,
		}, nil
	}
}

type multiResourceInformer struct {
	resourceToGVR      map[string]schema.GroupVersionResource
	resourceToInformer map[string]map[string]cache.SharedIndexInformer
	informerFactory    []dynamicinformer.DynamicSharedInformerFactory
}

var _ Informer = &multiResourceInformer{}

// AddEventHandler adds the handler to each namespaced informer
func (i *multiResourceInformer) AddEventHandler(handler eventHandler) {
	for _, ki := range i.resourceToInformer {
		for kind, informer := range ki {
			informer.AddEventHandler(handler(kind))
		}
	}
}

// HasSynced checks if each namespaced informer has synced
func (i *multiResourceInformer) HasSynced() bool {
	for _, ki := range i.resourceToInformer {
		for _, informer := range ki {
			if ok := informer.HasSynced(); !ok {
				return ok
			}
		}
	}

	return true
}

func (i *multiResourceInformer) Start(stopCh <-chan struct{}) {
	for _, informer := range i.informerFactory {
		informer.Start(stopCh)
	}
}
