package util

import (
	"Cloud/pkg/consts"
	"bytes"
	"flag"
	"github.com/wzyonggege/logger"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
)

var(
	clientset *kubernetes.Clientset
	config *rest.Config
	err error
)

func init(){
	//集群外配置client-go
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig",consts.PATH + "config1","absolute path to the kubeconfig file")

	flag.Parse()


	config, err = clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil{
		panic(err.Error())
	}

	clientset,err = kubernetes.NewForConfig(config)
	if err != nil{
		panic(err.Error())
	}
}


func CreateNamespaceIfNotExist(username string,Type bool) bool{
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: username,
		},
		Status: apiv1.NamespaceStatus{
			Phase: apiv1.NamespaceActive,
		},
	}
	namespacesClient := clientset.CoreV1().Namespaces()
	_,err := namespacesClient.Create(namespace)
	if err != nil{
		return false
	}

	if Type{
		return true
	}

	err = createResourceQuota(username)
	if err != nil{
		logger.Error("创建命名空间资源限制失败")
		return false
	}

	err = createLimitRange(username)
	if err != nil{
		logger.Error("创建命名空间容器默认配额失败")
		return false
	}

	return true
}

func GetPodItf(username string,Type bool) v1.PodInterface{
	CreateNamespaceIfNotExist(username,Type)
	return clientset.CoreV1().Pods(username)
}

func GetSvcItf(username string,Type bool) v1.ServiceInterface{
	CreateNamespaceIfNotExist(username,Type)
	return clientset.CoreV1().Services(username)
}

func GetDeployItf(username string,Type bool) v1beta1.DeploymentInterface {
	CreateNamespaceIfNotExist(username,Type)
	return clientset.AppsV1beta1().Deployments(username)
}

func GetPVItf(username string,Type bool) v1.PersistentVolumeInterface{
	CreateNamespaceIfNotExist(username,Type)
	return clientset.CoreV1().PersistentVolumes()
}

func GetPVCItf(username string,Type bool) v1.PersistentVolumeClaimInterface{
	CreateNamespaceIfNotExist(username,Type)
	return clientset.CoreV1().PersistentVolumeClaims(username)
}

func GetClientset()*kubernetes.Clientset{
	return clientset
}

func GetRestConfig()*rest.Config{
	return config
}

func CreateQuota(username string){
	CreateNamespaceIfNotExist(username,consts.NORMAL)
}

func createResourceQuota(username string)error{
	resourcequota := apiv1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			Kind: "ResourceQuota",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "resource-" + username,
			Namespace: username,
		},
		Spec: apiv1.ResourceQuotaSpec{
			Hard: apiv1.ResourceList{
				apiv1.ResourceLimitsCPU : resource.MustParse("2.4"),
				apiv1.ResourceRequestsCPU: resource.MustParse("2"),
				apiv1.ResourceLimitsMemory : resource.MustParse("4.8G"),
				apiv1.ResourceRequestsMemory : resource.MustParse("4G"),
				apiv1.ResourceLimitsEphemeralStorage: resource.MustParse("120G"),
				apiv1.ResourceRequestsEphemeralStorage: resource.MustParse("100G"),
			},
		},
	}
	rqItf := clientset.CoreV1().ResourceQuotas(username)
	_ , err = rqItf.Create(&resourcequota)
	return err
}

func createLimitRange(username string)error{
	limitrange := apiv1.LimitRange{
		TypeMeta: metav1.TypeMeta{
			Kind: "LimitRange",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "limit-" + username,
			Namespace: username,
		},
		Spec: apiv1.LimitRangeSpec{
			Limits: []apiv1.LimitRangeItem{
				{
					Type: apiv1.LimitTypeContainer,
					Default: apiv1.ResourceList{
						apiv1.ResourceCPU : resource.MustParse("0.6"),
						apiv1.ResourceMemory : resource.MustParse("0.6G"),
						apiv1.ResourceEphemeralStorage : resource.MustParse("0.6G"),
					},
					DefaultRequest: apiv1.ResourceList{
						apiv1.ResourceCPU: resource.MustParse("0.5"),
						apiv1.ResourceMemory : resource.MustParse("0.5G"),
						apiv1.ResourceEphemeralStorage: resource.MustParse("0.5G"),
					},
				},
			},
		},
	}
	lrItf := clientset.CoreV1().LimitRanges(username)
	_ , err = lrItf.Create(&limitrange)
	return err
}

func GetRQItf(username string,Type bool) v1.ResourceQuotaInterface {
	CreateNamespaceIfNotExist(username,Type)
	return  clientset.CoreV1().ResourceQuotas(username)
}

func GetNodeItf() v1.NodeInterface {
	return  clientset.CoreV1().Nodes()
}

func GetNSItf() v1.NamespaceInterface {
	return clientset.CoreV1().Namespaces()
}

func GetCoreClient()*v1.CoreV1Client{
	coreclient,_ := v1.NewForConfig(GetRestConfig())
	return coreclient
}

func ExecuteRemoteCommand(pod *apiv1.Pod,command []string)(string,error){
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	coreclient := GetCoreClient()

	buf := &bytes.Buffer{}

	req := coreclient.RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&apiv1.PodExecOptions{
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: os.Stderr,
	})
	if err != nil {
		return "",err
	}
	return buf.String(),nil
}