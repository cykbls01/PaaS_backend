package client

import (
	"Cloud/pkg/util"
)
import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"sync"
)

// use sync.Once, so clientset can only init once
func InitClient() (*kubernetes.Clientset, error) {
	var once sync.Once
	var clientset *kubernetes.Clientset
	once.Do(func() {
		log.Print("start doInit()")
		clientset, _ = doInit()
	})
	return clientset, nil
}

func CoreV1Client() corev1.CoreV1Interface {
	clientset, e := InitClient()
	if e != nil {
		log.Fatalf("init client error %s\n", e)
	}
	return clientset.CoreV1()
}

func doInit() (*kubernetes.Clientset, error) {

	// create the clientset
	clientset := util.GetClientset()
	return clientset, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if h := os.Getenv("USERPROFILE"); h != "" { // windows
		return h
	}

	return "/root"
}

func InitRestClient() (*rest.Config, error, *corev1client.CoreV1Client) {
	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig := util.GetRestConfig()
	// Create a Kubernetes core/v1 client.
	coreclient := util.GetCoreClient()
	return restconfig, nil, coreclient
}
