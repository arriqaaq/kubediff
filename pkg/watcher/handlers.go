package watcher

import (
	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/event"
	"github.com/arriqaaq/kubediff/pkg/log"
	"github.com/arriqaaq/kubediff/pkg/notify"
	"github.com/go-test/deep"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

const (
	EventAdd    string = "EventAdd"
	EventUpdate string = "EventUpdate"
	EventDelete string = "EventDelete"
)

type eventHandler func(resourceType string, notifier notify.Notifier) cache.ResourceEventHandlerFuncs

func watchHandler(resourceType string, notifier notify.Notifier) cache.ResourceEventHandlerFuncs {

	var handler cache.ResourceEventHandlerFuncs
	handler.AddFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).WithField("obj", obj).Info("add event")
		notifier.Handle(event.NewEvent(EventAdd, resourceType, obj, nil))
	}
	handler.UpdateFunc = func(old, new interface{}) {
		log.WithField("resourceType", resourceType).WithField("old", old).WithField("new", new).Info("update event")
		notifier.Handle(event.NewEvent(EventUpdate, resourceType, new, nil))
	}
	handler.DeleteFunc = func(obj interface{}) {
		log.WithField("resourceType", resourceType).WithField("obj", obj).Info("delete event")
		notifier.Handle(event.NewEvent(EventDelete, resourceType, obj, nil))
	}
	return handler
}

func diffHandler(resourceType string, notifier notify.Notifier) cache.ResourceEventHandlerFuncs {

	var handler cache.ResourceEventHandlerFuncs
	handler.UpdateFunc = func(old, new interface{}) {
		oldObj := old.(*unstructured.Unstructured)
		newObj := new.(*unstructured.Unstructured)

		if !equality.Semantic.DeepEqual(old, new) {
			diff := deep.Equal(oldObj, newObj)
			log.WithField("resourceType", resourceType).WithField("diff", diff).Info("update event")
			notifier.Handle(event.NewEvent(EventUpdate, resourceType, old, diff))
		}
	}
	return handler
}

func noOpHandler(resourceType string, notifier notify.Notifier) cache.ResourceEventHandlerFuncs {

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
