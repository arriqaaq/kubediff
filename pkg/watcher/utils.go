package watcher

import (
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func newKubeConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func getGVRFromResource(disco *restmapper.DeferredDiscoveryRESTMapper, resource string) (schema.GroupVersionResource, error) {
	var gvr schema.GroupVersionResource
	if strings.Count(resource, "/") >= 2 {
		s := strings.SplitN(resource, "/", 3)
		gvr = schema.GroupVersionResource{Group: s[0], Version: s[1], Resource: s[2]}
	} else if strings.Count(resource, "/") == 1 {
		s := strings.SplitN(resource, "/", 2)
		gvr = schema.GroupVersionResource{Group: "", Version: s[0], Resource: s[1]}
	}

	if _, err := disco.ResourcesFor(gvr); err != nil {
		return schema.GroupVersionResource{}, err
	}
	return gvr, nil
}

func getNamespace(ns string) string {
	switch ns {
	case "all":
		return metav1.NamespaceAll
	default:
		return ns
	}
}
