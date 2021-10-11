package k8s

import (
	"context"
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type NodeInfo struct{
	Node		string		`json:"name"`
	Usedmem		int64		`json:"usedMem"`
	Totmem		int64		`json:"totalMem"`

	UsedCPU		int64		`json:"usedCPU"`
	TotCPU		int64		`json:"totalCPU"`

	UsedSto		int64		`json:"usedStorage"`
	TotSto		int64		`json:"totalStorage"`
}

type PodInfo struct{
	Pode		string		`json:"name"`
	Usedmem		int64		`json:"usedMem"`
	Totmem		int64		`json:"totalMem"`

	UsedCPU		int64		`json:"usedCPU"`
	TotCPU		int64		`json:"totalCPU"`
}

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

func GetMetrics(pods *PodMetricsList,username string) error {
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/"+username+"/pods").DoRaw(context.TODO())
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	return err
}

func GetPodMetric(pod *PodMetrics,username string,podname string) error{
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/"+username+"/pods/" + podname).DoRaw(context.TODO())
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pod)
	return err
}

func AnalyzePorMetric(podName,username string)PodInfo{
	var result PodInfo
	var pod PodMetrics
	GetPodMetric(&pod,username,podName)
	if(len(pod.Containers)>0) {
		mem := resource.MustParse(pod.Containers[0].Usage.Memory)
		cpu := resource.MustParse(pod.Containers[0].Usage.CPU)
		result.Usedmem = mem.ScaledValue(resource.Mega)
		result.UsedCPU = cpu.MilliValue()
	}


	podItf := GetPodItf(username,true)
	info,_ := podItf.Get(context.TODO(),podName,meta_v1.GetOptions{})
	result.TotCPU = info.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
	result.Totmem = info.Spec.Containers[0].Resources.Requests.Memory().ScaledValue(resource.Mega)
	return result
}

func GetNodeMetric(nodes *NodeMetricsList) error {
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").DoRaw(context.TODO())
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &nodes)
	return err
}
