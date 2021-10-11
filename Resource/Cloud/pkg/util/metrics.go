package util

import (
	"encoding/json"
	"k8s.io/client-go/kubernetes"
	"time"
)

type PodMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Timestamp  time.Time `json:"timestamp"`
		Window     string    `json:"window"`
		Containers []struct {
			Name  string `json:"name"`
			Usage struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"usage"`
		} `json:"containers"`
	} `json:"items"`
}

type PodMetrics struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`

	// The following fields define time interval from which metrics were
	// collected from the interval [Timestamp-Window, Timestamp].
	Timestamp  time.Time `json:"timestamp"`
	Window     string    `json:"window"`

	// Metrics for all containers are collected within the same time window.
	Containers []struct {
		Name  string `json:"name"`
		Usage struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"containers"`
}

type NodeMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			SelfLink          string    `json:"selfLink"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Timestamp  time.Time `json:"timestamp"`
		Window     string    `json:"window"`
		Usage struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"items"`
}

func GetMetrics(clientset *kubernetes.Clientset, pods *PodMetricsList,username string) error {
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/"+username+"/pods").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	return err
}

func GetPodMetric(clientset *kubernetes.Clientset, pod *PodMetrics,username string,podname string) error{
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/"+username+"/pods/" + podname).DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pod)
	return err
}

func GetNodeMetric(clientset *kubernetes.Clientset, nodes *NodeMetricsList) error {
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &nodes)
	return err
}
