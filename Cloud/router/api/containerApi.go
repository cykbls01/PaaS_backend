package api

import (
	"Cloud/k8s"
	"Cloud/model"
	"Cloud/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "k8s.io/api/apps/v1beta1"
	v12 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"log"
	"strconv"
	"strings"
)
var(
	deployment = appsv1.Deployment{}

	service = v12.Service{}

	YAML []byte
	JSON []byte

	err error
	code int
)
type Children struct {	//label-value结构
	Value 	string `json:"value"`
	Label 	string `json:"label"`
	Comment string `json:"comment,omitempty"`
}
type form struct {					//创建容器的表单
	Name 		string 				`json:"name"`
	Podname		string				`json:"podname"`
	Image		string			`json:"image"`
	Env			[]Children			`json:"env"`
	Port 		int				`json:"port"`
	Url			string				`json:"url"`
	Msg			string				`json:"msg"`
	EndTime     string              `json:"end_time"`
}

type podInfo struct{
	Name        string		`json:"name"`
	Podname		string		`json:"podname"`
	Image		string	    `json:"image"`
	ApplyTime	string		`json:"apply_time"`
	DueTime     string      `json:"due_time"`
	Url         string      `json:"url"`
	Port        string      `json:"port"`
	NameSpace   string      `json:"namespace"`
	Status      string      `json:"status"`
	Info        k8s.PodInfo `json:"info"`
}

func analyzeDeployment(forminfo form){
	var replicas int32
	var container v12.Container
	var envs []v12.EnvVar
	var env []Children
	deployment.APIVersion = "apps/v1"
	deployment.Kind = "Deployment"
	deployment.SetName(forminfo.Name)
	env = forminfo.Env
	replicas = 1
	tmp := make(map[string]string)
	tmp["app"] = forminfo.Name
	container.Image = forminfo.Image
	container.Name = forminfo.Name

	for i := 0;i < len(forminfo.Env);i++{
		var envvar v12.EnvVar
		envvar.Name = env[i].Label
		envvar.Value = env[i].Value
		envs = append(envs,envvar)
	}
	container.Env = envs

	deploySpec := appsv1.DeploymentSpec{
		Replicas: &replicas,
		Template: v12.PodTemplateSpec{
			ObjectMeta: meta_v1.ObjectMeta{
				Labels: tmp,
			},
			Spec: v12.PodSpec{
				Containers: []v12.Container{container},
			},
		},
	}
	deployment.Spec = deploySpec
	deployment.Spec.Selector = &meta_v1.LabelSelector{
		MatchLabels: tmp,
	}
}

func analyzeService(forminfo form){
	tmp := make(map[string]string)
	tmp["app"] = forminfo.Name
	port := forminfo.Port
	svcport := []v12.ServicePort{}
	//for i:=0 ; i < len(port) ; i++{
		//porti := port[i]
	porti := port
		svcport=append(svcport,v12.ServicePort{Port: int32(porti),TargetPort: intstr.FromInt(porti),Name: fmt.Sprint(porti)})
	//}
	service.APIVersion = "v1"
	service.Kind = "Service"
	service.SetName(forminfo.Name)
	service.Spec = v12.ServiceSpec{
		Type: v12.ServiceTypeNodePort,
		Ports: svcport,
		Selector: tmp,
	}
}


func ApplyContainer(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	var forminfo form
	err := c.BindJSON(&forminfo)
    log.Println(err)

	// 解析deployment
    analyzeDeployment(forminfo)
	b,_ :=json.Marshal(deployment)
	util.WriteFile(id+"/"+forminfo.Name+"_deployment.json",b)

	// 解析service
    analyzeService(forminfo)
	b,_ =json.Marshal(service)
	util.WriteFile(id+"/"+forminfo.Name+"_service.json",b)

	// 生成申请
	type pod struct{
		Podname		string				`json:"podname"`
		Image		string			`json:"image"`
		//Port 		[]int				`json:"port"`
		Port 		int				`json:"port"`
	}
	var p pod
	p.Image = forminfo.Image
	p.Podname = forminfo.Podname
	p.Port = forminfo.Port
	d,_ := json.Marshal(&p)
	model.AddApply(id,forminfo.Name,string(d),forminfo.Msg,forminfo.EndTime)

	util.ReturnSuccess(c,"成功申请",nil)
	return
}

func CreateContainer(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		ApplyId int `json:"applyId"`
	}
	var JsonData Json
	err := c.BindJSON(&JsonData)
	log.Println(err)
	apply := model.GetApply(JsonData.ApplyId)

	// 部署容器
	data := util.ReadFile("./"+id+"/"+apply.VmName+"_deployment.json")
	json.Unmarshal(data,&deployment)
	DeployItf := k8s.GetDeployItf(id,true)
    _,err = DeployItf.Create(context.TODO(),&deployment,meta_v1.CreateOptions{})

    if err != nil{
    	log.Println(err)
    	util.ReturnError(c,"部署容器失败")
    	return
	}

	// 部署端口
	data = util.ReadFile("./"+id+"/"+apply.VmName+"_service.json")
	json.Unmarshal(data,&service)
	SvcItf := k8s.GetSvcItf(id,true)
	_,err = SvcItf.Create(context.TODO(),&service,meta_v1.CreateOptions{})

	if err != nil{
		fmt.Println(err)
		util.ReturnError(c,"部署端口失败")
		return
	}

	// 更新申请
	apply.FinishTime = util.GetNowTime()
	apply.ReplyMsg = "创建成功"
	apply.Status = 2
	model.UpdateApply(apply)

	util.ReturnSuccess(c,"成功创建容器",nil)
	return
}

func DeleteContainer(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		Name string `json:"name"`
	}
	var JsonData Json
	err := c.BindJSON(&JsonData)
	log.Println(err)

	// 删除容器和端口
	SvcItf := k8s.GetSvcItf(id,true)
	err = SvcItf.Delete(context.TODO(),JsonData.Name,meta_v1.DeleteOptions{})
	log.Println(err)
	k8s.DeleteDeployment(JsonData.Name,id)
	log.Println(err)

	util.ReturnSuccess(c,"成功删除容器",nil)
	return
}

func UpdateContainer(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		PodName string `json:"podName"`
		Port    int    `json:"port"`
	}
	var JsonData Json
	err := c.BindJSON(&JsonData)
	log.Println(err)


	// 部署端口
	data := util.ReadFile("./"+id+"/"+JsonData.PodName+"_service.json")
	err = json.Unmarshal(data,&service)

	porti := JsonData.Port
	service.Spec.Ports[0].Port = int32(porti)
	service.Spec.Ports[0].TargetPort = intstr.FromInt(porti)
	SvcItf := k8s.GetSvcItf(id,true)
	SvcItf.Delete(context.TODO(),JsonData.PodName,meta_v1.DeleteOptions{})
	SvcItf.Create(context.TODO(),&service,meta_v1.CreateOptions{})

	if err != nil{
		util.ReturnError(c,"部署端口失败")
		return
	}

	util.ReturnSuccess(c,"成功修改容器",nil)
	return
}

func analyzePod(podList *v12.PodList)[]*podInfo{
	var pl []*podInfo
	for _,v := range podList.Items{
		var p podInfo
		SvcItf := k8s.GetSvcItf(v.Namespace,true)
		p.Name = v.Name
		p.Podname = strings.Split(v.Name,"-")[0];
		p.NameSpace = v.Namespace
		p.Status = string(v.Status.Phase)
		apply := model.GetApplyByName(v.Namespace,p.Podname)
		service,_ := SvcItf.Get(context.TODO(),p.Podname,meta_v1.GetOptions{})
		if(len(service.Spec.Ports)>0) {
			p.Url = "http://10.251.253.51:"+strconv.Itoa(int(service.Spec.Ports[0].NodePort))
			p.Port = strconv.Itoa(int(service.Spec.Ports[0].Port))
		}
		p.ApplyTime = apply.ApplyTime
		p.DueTime = apply.DueTime
		if(len(v.Spec.Containers)>0) {
			p.Image = v.Spec.Containers[0].Image
		}
		p.Info = k8s.AnalyzePorMetric(v.Name,v.Namespace)
		pl = append(pl,&p)
	}
	return pl
}

func ListContainerByUser(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 封装参数
	PodItf := k8s.GetPodItf(id,true)
	pl,_ := PodItf.List(context.TODO(),meta_v1.ListOptions{})
    podList := analyzePod(pl)

	util.ReturnSuccess(c,"成功获取容器列表",podList)
	return
}

func ListContainerByAdmin(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 封装参数
	user := model.GetUser(id)
	PodItf := k8s.GetPodItf("",true)
	pl,_ := PodItf.List(context.TODO(),meta_v1.ListOptions{})

	if user.Role == 3{
		for _,v := range pl.Items{
			use := model.GetUser(v.Namespace)
			if use == nil || use.Department != user.Department{
			}
		}
	}else if user.Role == 4{
	}

	podList := analyzePod(pl)


	util.ReturnSuccess(c,"成功获取容器列表",podList)
	return
}
