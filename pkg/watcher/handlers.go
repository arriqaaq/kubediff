package watcher

import (
	"github.com/arriqaaq/kdiff/config"
	"github.com/arriqaaq/kdiff/pkg/log"
	"github.com/go-test/deep"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

type eventHandler func(resourceType string) cache.ResourceEventHandlerFuncs

func watchHandler(resourceType string) cache.ResourceEventHandlerFuncs {

	var handler cache.ResourceEventHandlerFuncs
	handler.AddFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).WithField("obj", obj).Info("add event")
	}
	handler.UpdateFunc = func(old, new interface{}) {
		log.WithField("resourceType", resourceType).WithField("old", old).WithField("new", new).Info("update event")
	}
	handler.DeleteFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).WithField("obj", obj).Info("delete event")
	}
	return handler
}

func diffHandler(resourceType string) cache.ResourceEventHandlerFuncs {

	var handler cache.ResourceEventHandlerFuncs
	handler.AddFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).WithField("obj", obj).Info("add event")
	}
	handler.UpdateFunc = func(old, new interface{}) {
		oldObj := old.(*unstructured.Unstructured)
		newObj := new.(*unstructured.Unstructured)

		if !equality.Semantic.DeepEqual(old, new) {
			diff := deep.Equal(oldObj, newObj)
			log.WithField("resourceType", resourceType).WithField("diff", diff).Info("update event")
		}
	}
	handler.DeleteFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).WithField("obj", obj).Info("delete event")
	}
	return handler
}

func noOpHandler(resourceType string) cache.ResourceEventHandlerFuncs {

	var handler cache.ResourceEventHandlerFuncs
	handler.AddFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).Info("delete event")
	}
	handler.UpdateFunc = func(old, new interface{}) {
		log.WithField("resourceType", resourceType).Info("delete event")
	}
	handler.DeleteFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).Info("delete event")
	}
	return handler
}

func getEventHandler(mode config.RunMode) eventHandler {
	switch mode {
	case config.DiffMode:
		return diffHandler
	default:
		return watchHandler
	}
}
