package watcher

import (
	"time"

	"github.com/arriqaaq/kdiff/config"
)

const (
	resyncPeriod = time.Duration(1) * time.Minute
)

func NewWatcher(cfg *config.Config) (*Watcher, error) {
	kubeconfig, err := newKubeConfig()
	if err != nil {
		return nil, err
	}

	client, err := newClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	informer, err := NewMultiResourceInformer(cfg, resyncPeriod)(client)
	if err != nil {
		return nil, err
	}
	informer.AddEventHandler(getEventHandler(cfg.Mode))

	return &Watcher{client, informer}, nil
}

type Watcher struct {
	client   *Client
	informer Informer
}

func (w *Watcher) Run(stopCh chan struct{}) {
	w.informer.Start(stopCh)
}
