package api

import (
	"Cloud/kubectl"
	"Cloud/models"
	"Cloud/pkg/consts"
	"Cloud/pkg/util"
	"Cloud/pkg/util/ws"
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wzyonggege/logger"
	"go.uber.org/zap"
	"io"
	"k8s.io/api/apps/v1beta1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	v1beta12 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"strings"
)

func createDeploy(username string,deployItf v1beta12.DeploymentInterface, name string,repo string,env []Children,limit map[string]int64,Type bool){
	code = consts.SUCCESS
	var replicas int32
	var container v12.Container
	var envs []v12.EnvVar
	deployment.APIVersion = "extensions/v1beta1"
	deployment.Kind = "Deployment"
	image := repo[strings.LastIndex(repo,"/")+1:]
	image = image[0:strings.Index(image,":")]
	path := models.GetImage(image).Mntpath
	deployment.SetName(name)
	replicas = 1
	tmp := make(map[string]string)
	tmp["app"] = name
	container.Image = repo
	container.Name = name
	container.VolumeMounts = []v12.VolumeMount{
		{
			MountPath: path,
			Name: name + "-pvc",
		},
	}
	var requestcpu,requestmemory,requeststorage float64
	flag := 0
	if limit != nil && len(limit) != 0{
		for k,v:=range(limit){
			switch k {
			case "requestCPU" :
				flag = flag | 1
				requestcpu = float64(v)
			case "requestMem":
				flag = flag | 2
				requestmemory = float64(v)
			case "requestStorage":
				flag = flag | 4
				requeststorage = float64(v)
			}
		}
		rlist := make(map[v12.ResourceName]resource.Quantity)
		llist := make(map[v12.ResourceName]resource.Quantity)
		if flag & 1 >0{
			rlist[v12.ResourceCPU] = resource.MustParse(fmt.Sprint(requestcpu) + "m")
			llist[v12.ResourceCPU ] = resource.MustParse(fmt.Sprint(requestcpu*1.2) + "m")
		}
		if flag & 2 > 0{
			rlist[v12.ResourceMemory] = resource.MustParse(fmt.Sprint(requestmemory)+"M")
			llist[v12.ResourceMemory] = resource.MustParse(fmt.Sprint(requestmemory*1.2)+"M")
		}
		if flag & 4 > 0{
			createPVC(util.GetPVCItf(username,consts.CLASS),name,requeststorage)
			if code != consts.SUCCESS{
				return
			}
		}
		container.Resources = v12.ResourceRequirements{
			Requests: rlist,
			Limits: llist,
		}
	}
	for i := 0;i < len(env);i++{
		var envvar v12.EnvVar
		envvar.Name = env[i].Label
		envvar.Value = env[i].Value
		envs = append(envs,envvar)
	}
	container.Env = envs
	deploySpec := v1beta1.DeploymentSpec{
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
	deploySpec.Template.Spec.Volumes = []v12.Volume{
		{
			Name:         name + "-pvc",
			VolumeSource: v12.VolumeSource{
				PersistentVolumeClaim: &v12.PersistentVolumeClaimVolumeSource{
					ClaimName: name,
				},
			},
		},
	}
	deployment.Spec = deploySpec
	if _,err = deployItf.Get(deployment.Name,meta_v1.GetOptions{});err != nil{

		if _,err = deployItf.Create(&deployment); err != nil {
			code = consts.ERROR_DEPLOY_CREATE
			return
		}
	}else {
		code = consts.ERROR_DEPLOY_ALREADY_EXIST
		return
	}
	return
}

func createDeploy1(username string,deployItf v1beta12.DeploymentInterface, name string,repo string,env []Children,limit map[string]int64,Type bool){
	code = consts.SUCCESS
	var replicas int32
	var container v12.Container
	var envs []v12.EnvVar
	deployment.APIVersion = "extensions/v1beta1"
	deployment.Kind = "Deployment"
	deployment.SetName(name)
	replicas = 1
	tmp := make(map[string]string)
	tmp["app"] = name
	container.Image = repo
	container.Name = name
	var requestcpu,requestmemory,requeststorage float64
	flag := 0
	if Type{
		goto KO
	}
	if limit != nil && len(limit) != 0{
		for k,v:=range(limit){
			switch k {
			case "requestCPU" :
				flag = flag | 1
				requestcpu = float64(v)
			case "requestMem":
				flag = flag | 2
				requestmemory = float64(v)
			case "requestStorage":
				flag = flag | 4
				requeststorage = float64(v)
			}
		}
		usage := GetResourceUsage(username)
		used := usage["used"]
		total := usage["total"]
		rlist := make(map[v12.ResourceName]resource.Quantity)
		llist := make(map[v12.ResourceName]resource.Quantity)
		if flag & 1 >0{
			usedCPU := used.Requestcpu
			totalCPU := total.Requestcpu
			if usedCPU + int64(requestcpu) > totalCPU{
				code = consts.ERROR_RESOURCE_OUT
				return
			}
			rlist[v12.ResourceCPU] = resource.MustParse(fmt.Sprint(requestcpu) + "m")
			llist[v12.ResourceCPU ] = resource.MustParse(fmt.Sprint(requestcpu*1.2) + "m")
		}
		if flag & 2 > 0{
			usedMem := used.Requestmem
			totalMem := total.Requestmem
			if usedMem + int64(requestmemory) > totalMem{
				code = consts.ERROR_RESOURCE_OUT
				return
			}
			rlist[v12.ResourceMemory] = resource.MustParse(fmt.Sprint(requestmemory)+"M")
			llist[v12.ResourceMemory] = resource.MustParse(fmt.Sprint(requestmemory*1.2)+"M")
		}
		if flag & 4 > 0{
			usedSto := used.Requeststo
			totalSto := total.Requeststo
			if usedSto + int64(requeststorage) > totalSto{
				code = consts.ERROR_RESOURCE_OUT
				return
			}
			rlist[v12.ResourceEphemeralStorage] = resource.MustParse(fmt.Sprint(requeststorage)+"M")
			llist[v12.ResourceEphemeralStorage] = resource.MustParse(fmt.Sprint(requeststorage*1.2)+"M")
		}
		container.Resources = v12.ResourceRequirements{
			Requests: rlist,
			Limits: llist,
		}
	}
KO:
	for i := 0;i < len(env);i++{
		var envvar v12.EnvVar
		envvar.Name = env[i].Label
		envvar.Value = env[i].Value
		envs = append(envs,envvar)
	}
	container.Env = envs
	deploySpec := v1beta1.DeploymentSpec{
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
	if _,err = deployItf.Get(deployment.Name,meta_v1.GetOptions{});err != nil{

		if _,err = deployItf.Create(&deployment); err != nil {
			code = consts.ERROR_DEPLOY_CREATE
			return
		}
	}else {
		code = consts.ERROR_DEPLOY_ALREADY_EXIST
		return
	}
	return
}

func deleteDeploy(deployItf v1beta12.DeploymentInterface,name string){

	code = consts.SUCCESS
	deletePlicy := meta_v1.DeletePropagationForeground
	err = deployItf.Delete(name,&meta_v1.DeleteOptions{
		PropagationPolicy: &deletePlicy,
	})
	if err != nil{
		code = consts.ERROR_DEPLOY_NOT_EXIST
		return
	}
	return
}

func createSvc(svcItf v1.ServiceInterface,name string,repo string,port []int){
	code = consts.SUCCESS
	tmp := make(map[string]string)
	tmp["app"] = name
	//svc := repo[:strings.Index(repo,":")]
	//port := consts.Getport(svc)
	svcport := []v12.ServicePort{}
	for i:=0 ; i < len(port) ; i++{
		porti := port[i]
		svcport=append(svcport,v12.ServicePort{Port: int32(porti),TargetPort: intstr.FromInt(porti),Name: fmt.Sprint(porti)})
	}
	service.APIVersion = "v1"
	service.Kind = "Service"
	service.SetName(name)
	service.Spec = v12.ServiceSpec{
		Type: v12.ServiceTypeNodePort,
		Ports: svcport,
		Selector: tmp,
	}
	if _,err = svcItf.Get(service.Name,meta_v1.GetOptions{});err != nil{
		if _,err = svcItf.Create(&service); err != nil {
			code = consts.ERROR_SERVICE_CREATE
			return
		}
	}else {
		code = consts.ERROR_SERVICE_ALREADY_EXIST
		return
	}
	return
}

func deleteSvc(svcItf v1.ServiceInterface,name string){
	code = consts.SUCCESS
	err = svcItf.Delete(name,&meta_v1.DeleteOptions{})
	if err != nil{
		code = consts.ERROR_SERVICE_NOT_EXIST
		return
	}
	return
}

func createPVC(pvcItf v1.PersistentVolumeClaimInterface,name string,size float64){
	code = consts.SUCCESS
	rlist := make(map[v12.ResourceName]resource.Quantity)
	llist := make(map[v12.ResourceName]resource.Quantity)
	rlist[v12.ResourceStorage] = resource.MustParse(fmt.Sprint(int64(size))+"M")
	llist[v12.ResourceStorage] = resource.MustParse(fmt.Sprint(int64(size*1.2))+"M")
	pvc := v12.PersistentVolumeClaim{
		TypeMeta:   meta_v1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
		Spec:       v12.PersistentVolumeClaimSpec{
			AccessModes: []v12.PersistentVolumeAccessMode{v12.ReadWriteMany},
			Resources:   v12.ResourceRequirements{
				Requests: rlist,
				Limits: llist,
			},
		},
	}
	if _,err = pvcItf.Get(name,meta_v1.GetOptions{});err != nil{

		if _,err = pvcItf.Create(&pvc); err != nil {
			code = consts.ERROR_PVC_CREATE
			return
		}
	}else {
		code = consts.ERROR_PVC_EXIST
		return
	}
	return
}

func deletePVC(pvcItf v1.PersistentVolumeClaimInterface,name string){
	code = consts.SUCCESS
	err = pvcItf.Delete(name,&meta_v1.DeleteOptions{})
	if err != nil{
		code = consts.ERROR_PVC_NOT_EXIST
		return
	}
	return
}

func Create(c *gin.Context){
	var data bool
	claims := checkToken(c)
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		var forminfo form
		_ = c.BindJSON(&forminfo)
		name := forminfo.Name
		imagepiece := forminfo.Image
		limit := forminfo.Limit
		env := forminfo.Env
		port := forminfo.Port
		var image string

		image2 := forminfo.Url
		if image2 != ""{
			image = image2
		}else{
			image = consts.HOST + ":5000/" + imagepiece[0] + ":" + imagepiece[1]
		}

		var deployItf v1beta12.DeploymentInterface
		var svcItf v1.ServiceInterface
		var podItf v1.PodInterface

		if username == "admin"{
			name = name + "-share"

			deployItf = util.GetDeployItf("share",consts.CLASS)
			svcItf = util.GetSvcItf("share",consts.CLASS)
			podItf = util.GetPodItf("share",consts.CLASS)
			createDeploy("share",deployItf,name,image,env,limit,consts.CLASS)
			if(strings.Contains(image,"mysql")){
				var password string
				for i := range env{
					if(env[i].Label=="MYSQL_ROOT_PASSWORD"){
						password = env[i].Value
						break
					}
				}
				models.AddShareContainer(name,password)
			}
		}else{
			name = name + "-" + username
			deployItf = util.GetDeployItf(username,consts.NORMAL)
			svcItf = util.GetSvcItf(username,consts.NORMAL)
			podItf = util.GetPodItf(username,consts.NORMAL)
			createDeploy1(username,deployItf,name,image,env,limit,consts.NORMAL)
		}
		if code != consts.SUCCESS{
			logger.Error("创建deployment "+name+" 失败 ",zap.String("reason",consts.GetMsg(code)),zap.String("user",username),zap.String("type","创建"))
			models.DeleteShareContainer(name)
			data = false
			goto LAST
		}
		createSvc(svcItf,name,image,port)
		if code != consts.SUCCESS{
			logger.Error("创建service "+name+" 失败 ",zap.String("reason",consts.GetMsg(code)),zap.String("user",username),zap.String("type","创建"))
			data = false
			deleteDeploy(deployItf,name)
			models.DeleteShareContainer(name)
			goto LAST
		}
		logger.Info("创建容器 "+name+" 成功",zap.String("user",username),zap.String("type","创建"))

		if(username != "admin"){
			if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
				for _, pod := range podList.Items {
					if(strings.Contains(pod.Name,name)){
						svc,_ := svcItf.Get(name,meta_v1.GetOptions{})
						size := len(svc.Spec.Ports)
						ports := []string{}
						svcports := svc.Spec.Ports
						for i:=0 ; i < size ; i++{
							ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
						}
						models.AddPodInfo(pod.Name,name,image,consts.HOST,ports,pod.GetCreationTimestamp().Unix(),username)
					}
				}
			}
		}
		UpdateUserQuota(models.GetUserInfo(username).ID)
		data = true
	}
LAST:
	//c.JSON(http.StatusOK, gin.H{
	//	"code" : code,
	//	"msg" :  consts.GetMsg(code),
	//	"data" : data,
	//})
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func Delete(c *gin.Context){
	claims := checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		username := claims.Username
		username2 := c.Param("username")
		if username2 != ""{
			username = username2
		}
		deployItf := util.GetDeployItf(username,consts.NORMAL)
		var delForm deleteForm
		_ = c.BindJSON(&delForm)
		names := delForm.Name

		svcItf := util.GetSvcItf(username,consts.NORMAL)
		podItf :=util.GetPodItf(username,consts.NORMAL)


		for i := 0 ; i < len(names);i++{
			podname := names[i]
			models.PodTerminating(podname)
			pod,err := podItf.Get(podname,meta_v1.GetOptions{})
			if err != nil{
				logger.Error("删除容器 "+podname+" 失败",zap.String("reason","容器不存在"),zap.String("user",username),zap.String("type","删除"))
				continue
			}
			if pod.Status.Phase == v12.PodFailed{
				if pod.Status.Reason == "Evicted"{
					podItf.Delete(podname,&meta_v1.DeleteOptions{})
					continue
				}
			}
			name := podname[:strings.LastIndex(podname,username)] + username

			deleteDeploy(deployItf,name)
			if code != consts.SUCCESS{
				logger.Error("删除deployment "+name+" 失败 ",zap.String("reason",consts.GetMsg(code)),zap.String("user",username),zap.String("type","删除"))
				goto LAST
			}
			deleteSvc(svcItf,name)
			if code != consts.SUCCESS{
				logger.Error("删除service "+name+" 失败",zap.String("reason",consts.GetMsg(code)),zap.String("user",username),zap.String("type","删除"))
				goto LAST
			}
			logger.Info("删除容器 "+name+" 成功",zap.String("user",username),zap.String("type","删除"))
			models.DeletePod(podname)
			go UpdateUserQuota(models.GetUserInfo(username).ID)
		}
	}
LAST:
	ListPods(c)
	return
}

//func CreateSharePod(c *gin.Context){
//	var data bool
//	_ = checkToken(c)
//	if code != consts.SUCCESS{
//		goto LAST
//	}else{
//		var forminfo form
//		_ = c.BindJSON(&forminfo)
//		name := forminfo.Name
//		imagepiece := forminfo.Image
//		limit := forminfo.Limit
//		env := forminfo.Env
//		port := forminfo.Port
//		image := imagepiece[0] + ":" + imagepiece[1]
//
//		name = name + "-" + "share"
//
//		deployItf := util.GetDeployItf("share",consts.CLASS)
//		svcItf := util.GetSvcItf("share",consts.CLASS)
//		createDeploy("share",deployItf,name,image,env,limit,consts.CLASS)
//		if code != consts.SUCCESS{
//			logger.Error("创建deployment "+name+" 失败 ",zap.String("reason",consts.GetMsg(code)),zap.String("user","admin"),zap.String("type","创建"))
//			data = false
//			goto LAST
//		}
//		createSvc(svcItf,name,image,port)
//		if code != consts.SUCCESS{
//			logger.Error("创建service "+name+" 失败 ",zap.String("reason",consts.GetMsg(code)),zap.String("user","admin"),zap.String("type","创建"))
//			data = false
//			deleteDeploy(deployItf,name)
//			goto LAST
//		}
//		logger.Info("创建共享容器 "+name+" 成功",zap.String("user","admin"),zap.String("type","创建"))
//		data = true
//	}
//LAST:
//	c.Set("code",code)
//	c.Set("msg",consts.GetMsg(code))
//	c.Set("data",data)
//	return
//}

func ListSharePod(c *gin.Context){
	var data []podinfo
	_ = checkToken(c)
	if code != consts.SUCCESS{
		goto LAST
	}else{
		data = getNormalPods("share",consts.CLASS)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

//管理员删除某共享容器
func DeleteSharePod(c *gin.Context){
	_ = checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		deployItf := util.GetDeployItf("share",consts.CLASS)

		var info struct{
			Podname 	string		`json:"podname"`
		}

		c.BindJSON(&info)

		svcItf := util.GetSvcItf("share",consts.CLASS)

		podname := info.Podname
		name := GetContainerName(podname)

		deleteDeploy(deployItf,name)
		if code != consts.SUCCESS{
			logger.Error("删除deployment "+name+" 失败 ",zap.String("reason",consts.GetMsg(code)),zap.String("user","admin"),zap.String("type","删除"))
			goto LAST
		}
		deleteSvc(svcItf,name)
		if code != consts.SUCCESS{
			logger.Error("删除service "+name+" 失败",zap.String("reason",consts.GetMsg(code)),zap.String("user","admin"),zap.String("type","删除"))
			goto LAST
		}
		deletePVC(util.GetPVCItf("share",consts.CLASS),name)
		models.DeleteShareContainer(name)
		models.DeleteSharepod(podname)
		logger.Info("删除容器 "+name+" 成功",zap.String("user","admin"),zap.String("type","删除"))
	}
LAST:
	ListSharePod(c)
	return
}

func getSharePods(username string)[]podinfo{
	var pods []podinfo
	sharedpods := models.GetPodIfExist(username)
	poditf := util.GetPodItf("share",consts.CLASS)//不用给default添加limitrange和resourcequota，选择CLASS模式
	svcitf := util.GetSvcItf("share",consts.CLASS)
	for i:=range sharedpods{
		var tmpinfo podinfo
		podname := sharedpods[i].Podname
		pod,_ := poditf.Get(podname,meta_v1.GetOptions{})
		svcname := GetContainerName(podname)
		svc,_ := svcitf.Get(svcname,meta_v1.GetOptions{})
		tmpinfo.Name = svcname + "-" + sharedpods[i].Dbname
		tmpinfo.Type = "共享"
		tmpinfo.Podname = podname
		tmpinfo.Image = pod.Spec.Containers[0].Image
		tmpinfo.Image = tmpinfo.Image[strings.Index(tmpinfo.Image,"/")+1:]
		tmpinfo.Createtime = sharedpods[i].Createtime
		imagename := tmpinfo.Image[0:strings.Index(tmpinfo.Image,":")]
		image := models.GetImage(imagename)
		if(len(image.Acceptfile)!=0){
			tmpinfo.Upload = uploadinfo{true,image.Acceptfile}
		}else{
			tmpinfo.Upload = uploadinfo{false,""}
		}
		tmpinfo.Addr = consts.HOST
		tmpinfo.Extrainfo = extrainfo{sharedpods[i].Dbname,sharedpods[i].Password}
		podStatus := string(pod.Status.Phase)
		if svc != nil && err == nil{
			size := len(svc.Spec.Ports)
			ports := []string{}
			svcports := svc.Spec.Ports
			for i:=0 ; i < size ; i++{
				ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
			}
			tmpinfo.Port = ports
		}else {
			tmpinfo.Port = nil
			podStatus = "Terminating"
			goto KO
		}
	KO:
		tmpinfo.Status = podStatus
		switch tmpinfo.Status {
		case "Running": tmpinfo.Status = "运行中"
		case "Pending": tmpinfo.Status = "创建中"
		case "Terminating" : tmpinfo.Status = "关闭中"
		case "Unschedulable" : tmpinfo.Status = "无法分配"
		case "Unkown" : tmpinfo.Status = "未知错误"
		case "Failed" : tmpinfo.Status = "容器出错"
		case "ContainersNotReady" : tmpinfo.Status = "容器未准备好"
		case "Succeeded" : tmpinfo.Status = "容器成功终止"
		case "Evicted" : tmpinfo.Status = "容器被驱逐"
		case "Error" : tmpinfo.Status = "错误"
		}
		pods = append(pods, tmpinfo)
	}
	return pods
}

func getClassPods(username string)[]podinfo{
	var pods []podinfo
	classpods := models.GetClassPods(username)
	for i:=range classpods{
		podname := classpods[i].Podname
		pod,_ := util.GetPodItf(classpods[i].Class,consts.CLASS).Get(podname,meta_v1.GetOptions{})
		svcname := podname[0:strings.LastIndex(podname[0:strings.LastIndex(podname,"-")],"-")]
		svc,_ := util.GetSvcItf(classpods[i].Class,consts.CLASS).Get(svcname,meta_v1.GetOptions{})
		var tmpinfo podinfo
		tmpinfo.Name = svcname
		tmpinfo.Type = "课程"
		tmpinfo.Podname = podname
		tmpinfo.Image = pod.Spec.Containers[0].Image
		tmpinfo.Image = tmpinfo.Image[strings.Index(tmpinfo.Image,"/")+1:]
		tmpinfo.Createtime = pod.CreationTimestamp.Unix()
		imagename := tmpinfo.Image[0:strings.Index(tmpinfo.Image,":")]
		image := models.GetImage(imagename)
		if(len(image.Acceptfile)!=0){
			tmpinfo.Upload = uploadinfo{true,image.Acceptfile}
		}else{
			tmpinfo.Upload = uploadinfo{false,""}
		}
		tmpinfo.Addr = consts.HOST
		podStatus := string(pod.Status.Phase)
		if svc != nil && err == nil{
			size := len(svc.Spec.Ports)
			ports := []string{}
			svcports := svc.Spec.Ports
			for i:=0 ; i < size ; i++{
				ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
			}
			tmpinfo.Port = ports
		}else {
			tmpinfo.Port = nil
			podStatus = "Terminating"
			goto KO
		}
	KO:
		tmpinfo.Status = podStatus
		switch tmpinfo.Status {
		case "Running": tmpinfo.Status = "运行中"
		case "Pending": tmpinfo.Status = "创建中"
		case "Terminating" : tmpinfo.Status = "关闭中"
		case "Unschedulable" : tmpinfo.Status = "无法分配"
		case "Unkown" : tmpinfo.Status = "未知错误"
		case "Failed" : tmpinfo.Status = "容器出错"
		case "ContainersNotReady" : tmpinfo.Status = "容器未准备好"
		case "Succeeded" : tmpinfo.Status = "容器成功终止"
		case "Evicted" : tmpinfo.Status = "容器被驱逐"
		case "Error" : tmpinfo.Status = "错误"
		}
		pods = append(pods, tmpinfo)
	}
	return pods
}

func getNormalPodsFromDB(username string,Type bool)[]podinfo{
	var pods []podinfo
	podItf := util.GetPodItf(username,Type)

	var podinfo1 podinfo

	podlist := models.GetUserPods(models.GetUserInfo(username).ID)
	for i := range podlist{
		pod := podlist[i]
		realpod,_ := podItf.Get(pod.Podname,meta_v1.GetOptions{})
		podStatus := string(realpod.Status.Phase)
		podinfo1.Name = pod.Name
		podinfo1.Podname = pod.Podname
		podinfo1.Type = "普通"
		podinfo1.Image = pod.Image
		podinfo1.Image = podinfo1.Image[strings.Index(podinfo1.Image,"/")+1:]
		imagename := podinfo1.Image[0:strings.Index(podinfo1.Image,":")]
		image := models.GetImage(imagename)
		if(len(image.Acceptfile)!=0){
			podinfo1.Upload = uploadinfo{true,image.Acceptfile}
		}else{
			podinfo1.Upload = uploadinfo{false,""}
		}
		podinfo1.Createtime = pod.Createtime
		podinfo1.Addr = pod.Addr
		ports := strings.Split(pod.Port+",",",")
		podinfo1.Port = ports[0:len(ports)-1]
		if(len(pod.Port)==0){
			podinfo1.Port = nil
			podStatus = "Terminating"
			goto KO
		}
		if realpod.Status.Phase == v12.PodFailed{
			podStatus = realpod.Status.Reason
			goto KO
		}

		//内存超标
		if len(realpod.Status.ContainerStatuses) != 0{
			for i:=0; i<len(realpod.Status.ContainerStatuses);i++{
				cs := realpod.Status.ContainerStatuses[i]
				if &cs.LastTerminationState != nil{
					if cs.LastTerminationState.Terminated != nil{
						podStatus = cs.LastTerminationState.Terminated.Reason
						goto KO
					}
				}
			}
		}
	KO:
		podinfo1.Status = podStatus
		switch podinfo1.Status {
		case "Running": podinfo1.Status = "运行中"
		case "Pending": podinfo1.Status = "创建中"
		case "Terminating" : podinfo1.Status = "关闭中"
		case "Unschedulable" : podinfo1.Status = "无法分配"
		case "Unkown" : podinfo1.Status = "未知错误"
		case "Failed" : podinfo1.Status = "容器出错"
		case "ContainersNotReady" : podinfo1.Status = "容器未准备好"
		case "Succeeded" : podinfo1.Status = "容器成功终止"
		case "Evicted" : podinfo1.Status = "容器被驱逐"
		case "Error" : podinfo1.Status = "错误"
		}
		pods = append(pods, podinfo1)
	}
	//if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
	//	for _, pod := range podList.Items {
	//		podName := pod.Name
	//		podStatus := string(pod.Status.Phase)
	//		idx := strings.LastIndex(podName,username)
	//		svcname := podName[:idx] + username
	//		svc,err := svcItf.Get(svcname,meta_v1.GetOptions{})
	//		podinfo1.Name = svcname
	//		podinfo1.Podname = podName
	//		podinfo1.Type = "普通"
	//		podinfo1.Image = pod.Spec.Containers[0].Image
	//		podinfo1.Image = podinfo1.Image[strings.Index(podinfo1.Image,"/")+1:]
	//		imagename := podinfo1.Image[0:strings.Index(podinfo1.Image,":")]
	//		image := models.GetImage(imagename)
	//		if(len(image.Acceptfile)!=0){
	//			podinfo1.Upload = uploadinfo{true,image.Acceptfile}
	//		}else{
	//			podinfo1.Upload = uploadinfo{false,""}
	//		}
	//		podinfo1.Createtime = pod.CreationTimestamp.Unix()
	//		podinfo1.Addr = consts.HOST
	//
	//		if svc != nil && err == nil{
	//			size := len(svc.Spec.Ports)
	//			ports := []string{}
	//			svcports := svc.Spec.Ports
	//			for i:=0 ; i < size ; i++{
	//				ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
	//			}
	//			podinfo1.Port = ports
	//		}else {
	//			podinfo1.Port = nil
	//			podStatus = "Terminating"
	//			goto KO
	//		}
	//
	//		//存储超标
	//		if pod.Status.Phase == v12.PodFailed{
	//			podStatus = pod.Status.Reason
	//			goto KO
	//		}
	//
	//		//内存超标
	//		if len(pod.Status.ContainerStatuses) != 0{
	//			for i:=0; i<len(pod.Status.ContainerStatuses);i++{
	//				cs := pod.Status.ContainerStatuses[i]
	//				if &cs.LastTerminationState != nil{
	//					if cs.LastTerminationState.Terminated != nil{
	//						podStatus = cs.LastTerminationState.Terminated.Reason
	//						goto KO
	//					}
	//				}
	//			}
	//		}
	//
	//		// PodRunning means the pod has been bound to a node and all of the containers have been started.
	//		// At least one container is still running or is in the process of being restarted.
	//		//if podStatus != string(v12.PodRunning) {
	//		//	// 汇总错误原因不为空
	//		//	if pod.Status.Reason != "" {
	//		//		podStatus = pod.Status.Reason
	//		//		goto KO
	//		//	}
	//		//
	//		//	// condition有错误信息
	//		//	for _, cond := range pod.Status.Conditions {
	//		//		if cond.Type == v12.PodReady {	// POD就绪状态
	//		//			if cond.Status != v12.ConditionTrue {	// 失败
	//		//				podStatus = cond.Reason
	//		//			}
	//		//			goto KO
	//		//		}
	//		//	}
	//		//
	//		//	// 没有ready condition, 状态未知
	//		//	podStatus = "Unknown"
	//		//}
	//	KO:
	//		podinfo1.Status = podStatus
	//		switch podinfo1.Status {
	//		case "Running": podinfo1.Status = "运行中"
	//		case "Pending": podinfo1.Status = "创建中"
	//		case "Terminating" : podinfo1.Status = "关闭中"
	//		case "Unschedulable" : podinfo1.Status = "无法分配"
	//		case "Unkown" : podinfo1.Status = "未知错误"
	//		case "Failed" : podinfo1.Status = "容器出错"
	//		case "ContainersNotReady" : podinfo1.Status = "容器未准备好"
	//		case "Succeeded" : podinfo1.Status = "容器成功终止"
	//		case "Evicted" : podinfo1.Status = "容器被驱逐"
	//		case "Error" : podinfo1.Status = "错误"
	//		}
	//		pods = append(pods, podinfo1)
	//	}
	//}
	return pods
}

func getNormalPods(username string,Type bool)[]podinfo {
	pods := []podinfo{}
	podItf := util.GetPodItf(username,Type)
	svcItf := util.GetSvcItf(username,Type)
	var podinfo1 podinfo
	if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
		for _, pod := range podList.Items {
			podName := pod.Name
			podStatus := string(pod.Status.Phase)
			idx := strings.LastIndex(podName,username)
			svcname := podName[:idx] + username
			svc,err := svcItf.Get(svcname,meta_v1.GetOptions{})
			podinfo1.Name = svcname
			podinfo1.Podname = podName
			flag ,_ := models.IsClassPod(podName)
			if(flag){
				cp := models.GetClassPod(podName)
				user := models.GetUserInfoById(cp.Userid)
				podinfo1.Owner.Username = user.Username
				podinfo1.Owner.Name = user.Name
			}
			podinfo1.Type = "普通"
			podinfo1.Image = pod.Spec.Containers[0].Image
			podinfo1.Image = podinfo1.Image[strings.Index(podinfo1.Image,"/")+1:]
			imagename := podinfo1.Image[0:strings.Index(podinfo1.Image,":")]
			image := models.GetImage(imagename)
			if(len(image.Acceptfile)!=0){
				podinfo1.Upload = uploadinfo{true,image.Acceptfile}
			}else{
				podinfo1.Upload = uploadinfo{false,""}
			}
			podinfo1.Createtime = pod.CreationTimestamp.Unix()
			podinfo1.Addr = consts.HOST

			if svc != nil && err == nil{
				size := len(svc.Spec.Ports)
				ports := []string{}
				svcports := svc.Spec.Ports
				for i:=0 ; i < size ; i++{
					ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
				}
				podinfo1.Port = ports
			}else {
				podinfo1.Port = nil
				podStatus = "Terminating"
				goto KO
			}

			//存储超标
			if pod.Status.Phase == v12.PodFailed{
				podStatus = pod.Status.Reason
				goto KO
			}

			//内存超标
			if len(pod.Status.ContainerStatuses) != 0{
				for i:=0; i<len(pod.Status.ContainerStatuses);i++{
					cs := pod.Status.ContainerStatuses[i]
					if &cs.LastTerminationState != nil{
						if cs.LastTerminationState.Terminated != nil{
							podStatus = cs.LastTerminationState.Terminated.Reason
							goto KO
						}
					}
				}
			}
		KO:
			podinfo1.Status = podStatus
			switch podinfo1.Status {
			case "Running": podinfo1.Status = "运行中"
			case "Pending": podinfo1.Status = "创建中"
			case "Terminating" : podinfo1.Status = "关闭中"
			case "Unschedulable" : podinfo1.Status = "无法分配"
			case "Unkown" : podinfo1.Status = "未知错误"
			case "Failed" : podinfo1.Status = "容器出错"
			case "ContainersNotReady" : podinfo1.Status = "容器未准备好"
			case "Succeeded" : podinfo1.Status = "容器成功终止"
			case "Evicted" : podinfo1.Status = "容器被驱逐"
			case "Error" : podinfo1.Status = "错误"
			}
			pods = append(pods, podinfo1)
		}
	}
	return pods
}

func ListPods(c *gin.Context){
	var pods []podinfo
	var quota map[string]resourceinfo
	var data struct{
		Username 	string						`json:"username"`
		Pods		[]podinfo 					`json:"pods"`
		Quota		map[string]resourceinfo		`json:"quota"`
	}
	claims := checkToken(c)
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username

		username2 := c.Param("username")
		if username2 != ""{
			username = username2
		}
		data.Username = username
		quota = GetResourceUsage(username)
		if(code != consts.SUCCESS){
			code = consts.SUCCESS
		}
		pods = append(pods, getNormalPodsFromDB(username,consts.CLASS)...)
		pods = append(pods,getSharePods(username)...)
		pods = append(pods,getClassPods(username)...)
		data.Pods = pods
		data.Quota = quota
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return

}

func joinShareMysql(podname string,container string,username string,dbname string,password string){
	code = consts.SUCCESS
	pod := kubectl.Pod(podname,"share",container)
	pod.CopyToPod("/root/shells/mysql.sh","/root")
	pod.CopyToPod("/root/shells/mysql1.sh","/root")
	passwd := models.GetSharePassword(container)						//容器本身mysql的root密码
	err = pod.Exec([]string{"bash","/root/mysql.sh",dbname,passwd,password})
	if err != nil{
		logger.Error("共享mysql容器 " +dbname+ " 失败",zap.String("reason","非法输入"),zap.String("user",username),zap.String("type","创建"))
		code = consts.ERROR_DB_EXEC
	}
}

func JoinSharepod(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		var info struct{
			Podname		string		`json:"podname"`
			Image		string		`json:"image"`
			DBName 		string 		`json:"dbname"`
			Password	string 		`json:"password"`
		}
		_ = c.BindJSON(&info)
		data = false
		if models.CheckExist(info.DBName,info.Podname){
			data = false
			//已存在
			logger.Error("共享容器 " +info.DBName+ " 失败",zap.String("reason","数据库已存在"),zap.String("user",username),zap.String("type","创建"))
			code = consts.ERROR_DB_EXIST
			goto LAST
		}
		podname:= info.Podname
		container := GetContainerName(podname)

		if(strings.Contains(info.Image,"mysql")){
			joinShareMysql(podname,container,username,info.DBName,info.Password)
			if(code != consts.SUCCESS){
				goto LAST
			}
		}else if(strings.Contains(info.Image,"sql-server")){

		}

		data = true
		//添加数据库信息
		models.AddToDatabase(podname,username,info.DBName,info.Password)
		logger.Info("共享容器成功",zap.String("user",username),zap.String("dbname",info.DBName+"_"+username),zap.String("type","创建"))
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func deleteShareMysql(podname string,container string,dbname string,username string){
	pod := kubectl.Pod(podname,"share",container)
	passwd := models.GetSharePassword(container)
	err := pod.Exec([]string{"bash","/root/mysql1.sh",dbname,passwd})
	if err != nil{
		logger.Error("取消共享mysql失败",zap.String("reason","脚本执行失败"),zap.String("user",username),zap.String("dbname",dbname),zap.String("type","删除"))
		code = consts.ERROR_DB_DELETE
	}
}

//用户取消共享某容器
func RemoveSharepod(c *gin.Context){
	claims := checkToken(c)
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var db struct{
			Dbname		string		`json:"dbname"`
			Image		string		`json:"image"`
			Podname		string		`json:"podname"`
		}
		_ = c.BindJSON(&db)
		podname:= db.Podname
		username := claims.Username

		username2 := c.Param("username")
		if username2 != "" {
			username = username2
		}
		container := GetContainerName(podname)
		dbname := db.Dbname
		if models.CheckExist(dbname,podname)==false{
			//不存在
			logger.Error("取消共享容器失败",zap.String("reason","数据库不存在"),zap.String("user",username),zap.String("dbname",db.Dbname),zap.String("type","删除"))
			code = consts.ERROR_DB_NOT_EXIST
			goto LAST
		}
		if(strings.Contains(db.Image,"mysql")){
			deleteShareMysql(podname,container,db.Dbname,username)
			if(code != consts.SUCCESS){
				goto LAST
			}
		}else if(strings.Contains(db.Image,"sql-server")){

		}
		models.DeleteSharepodInfoFromDB(db.Dbname)
		logger.Info("取消共享容器成功",zap.String("user",username),zap.String("dbname",db.Dbname),zap.String("type","删除"))
	}
LAST:
	ListPods(c)
	return
}

func GetPodLog(c *gin.Context){
	var data string
	claims := checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		username := claims.Username
		username2 := c.Param("username")
		if username2 != ""{
			username = username2
		}
		podItf := util.GetPodItf(username,consts.NORMAL)
		var forminfo form
		_ = c.BindJSON(&forminfo)
		name := forminfo.Podname
		linenum := forminfo.Linenum
		if linenum < 0 {
			linenum = 0
		}

		deployItf := util.GetPodItf(username,consts.NORMAL)

		if _,err = deployItf.Get(name,meta_v1.GetOptions{});err != nil{
			logger.Error("读取容器日志失败",zap.String("reason","容器不存在"),zap.String("type","读日志"))
			code = consts.ERROR_POD_NOT_EXIST
			goto LAST
		}

		req := podItf.GetLogs(name,&v12.PodLogOptions{})
		podLogs, _ := req.Stream()

		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)

		data = buf.String()

		switch linenum {
		case 0:
			goto LAST
		default:
			datatmp := ""
			if data == ""{
				goto LAST
			}
			data = data[:len(data)-1]
			for ;linenum > 0;linenum -- {
				idx := strings.LastIndex(data,"\n")
				if idx < 0{
					datatmp = data + datatmp
					break;
				}
				datatmp = data[idx:] + datatmp
				data = data[:idx]
			}
			data = datatmp
			if data[0] == '\n'{
				data = data[1:]
			}
			data = data + "\n"
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func IsNameExist(c *gin.Context){
	code = consts.SUCCESS
	var data bool
	claims := checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		username := claims.Username
		var forminfo form
		_ = c.BindJSON(&forminfo)
		name := forminfo.Name

		if username == "admin" {
			username = "share"
		}

		name = name + "-" + username

		data = models.IsPodExist(name)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func Ssh(c *gin.Context){
	var (
		wsConn *ws.WsConnection
		restConf *rest.Config
		sshReq *rest.Request
		podName string
		containerName string
		executor remotecommand.Executor
		handler *util.StreamHandler
		token string
		data bool
	)
	podName = c.Query("podname")
	podName1,err := base64.StdEncoding.DecodeString(podName)
	if err != nil{
		logger.Error("podname解码失败",zap.String("type","ssh"))
		code = consts.ERROR_SSH_DECODE
		goto LAST
	}
	podName = string(podName1)
	token = c.Query("token")
	code = consts.SUCCESS
	if token == ""{
		code = consts.ERROR_AUTH_TOKEN
		logger.Error("token为空",zap.String("type","ssh"))
		goto LAST
	}else {
		_,err := util.ParseToken(token)
		if err != nil{
			code = consts.ERROR_AUTH_TOKEN
			logger.Error("token无效",zap.String("type","ssh"))
			goto LAST
		}else {
			//username := claims.Username
			//username2 := c.Query("username")
			//if username2 != ""{
			//	username = username2
			//}
			containerName = GetContainerName(podName)

			var ns string

			flag,class:=models.IsClassPod(podName)
			if(flag){
				ns = class
			}else{
				ns = containerName[strings.LastIndex(containerName,"-")+1:]
			}

			// 得到websocket长连接
			if wsConn, err = ws.InitWebsocket(c.Writer, c.Request); err != nil {
				return
			}

			sshReq = util.GetClientset().CoreV1().RESTClient().Post().
				Resource("pods").
				Name(podName).
				Namespace(ns).
				SubResource("exec").
				VersionedParams(&v12.PodExecOptions{
					Container:	containerName,
					Command: 	[]string{"bash"},
					Stdin: true,
					Stdout: true,
					Stderr: true,
					TTY: true,
				},scheme.ParameterCodec)

			restConf = util.GetRestConfig()

			// 创建到容器的连接
			if executor, err = remotecommand.NewSPDYExecutor(restConf, "POST", sshReq.URL()); err != nil {
				goto LAST
			}

			// 配置与容器之间的数据流处理回调
			handler = &util.StreamHandler{ WsConn: wsConn, ResizeEvent: make(chan remotecommand.TerminalSize)}
			if err = executor.Stream(remotecommand.StreamOptions{
				Stdin:             handler,
				Stdout:            handler,
				Stderr:            handler,
				TerminalSizeQueue: handler,
				Tty:               true,
			}); err != nil {
				goto LAST
			}
			data = true
			return
		}

	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	wsConn.WsClose()
	return
}

func GetPodUsage(c *gin.Context){
	var data struct{
		UsedCPU		int64		`json:"usedCPU"`
		UsedMem		int64		`json:"usedMem"`
		UsedStorage	int64		`json:"usedStorage"`
		TotCPU		int64		`json:"totCPU"`
		TotMem		int64		`json:"totMem"`
		TotStorage	int64		`json:"totStorage"`
	}
	claims := checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		var info struct{
			Podname		string		`json:"podname"`
			Image		string		`json:"image"`
		}
		c.BindJSON(&info)
		username := claims.Username
		username2 := c.Param("username")
		if username2 != ""{
			username = username2
		}
		podItf := util.GetPodItf(username,consts.CLASS)
		pod,_ := podItf.Get(info.Podname,meta_v1.GetOptions{})
		if(strings.Contains(info.Image,":")){
			info.Image = info.Image[0:strings.Index(info.Image,":")]
		}
		image := models.GetImage(info.Image)
		var podmetrics util.PodMetrics
		_ = util.GetPodMetric(util.GetClientset(),&podmetrics,username,info.Podname)
		container := podmetrics.Containers[0]
		cpu := resource.MustParse(container.Usage.CPU)
		mem := resource.MustParse(container.Usage.Memory)
		storage,err := util.ExecuteRemoteCommand(pod,[]string{"du","-sh",image.Dupath})
		if(err != nil){
			println("查看容量出错")
			goto LAST
		}
		storage = storage[0:strings.Index(storage,"\t")]
		storage = strings.ReplaceAll(storage,"K","k")
		sto := resource.MustParse(storage)
		data.UsedCPU = cpu.MilliValue()
		data.UsedMem = mem.ScaledValue(resource.Mega)
		data.UsedStorage = sto.ScaledValue(resource.Mega)
		ctner := pod.Spec.Containers[0]
		tcpu := ctner.Resources.Requests.Cpu().MilliValue()
		tmem := ctner.Resources.Requests.Memory().ScaledValue(resource.Mega)
		tsto := ctner.Resources.Requests.StorageEphemeral().ScaledValue(resource.Mega)
		data.TotCPU = tcpu
		data.TotMem = tmem
		data.TotStorage = tsto
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func getSimpleClassPod(classen string)[]podinfo{
	var pods []podinfo

	var podinfo1 podinfo

	podlist := models.GetClassPodsByClassen(classen)
	for i := range podlist{
		pod := podlist[i]
		user := models.GetUserInfoById(pod.Userid)
		podinfo1.Name = GetContainerName(pod.Podname)
		podinfo1.Podname = pod.Podname
		podinfo1.Type = "普通"
		podinfo1.Owner.Username = user.Username
		podinfo1.Owner.Name = user.Name
		podinfo1.Image = pod.Image
		podinfo1.Createtime = pod.Createtime
		podinfo1.Addr = pod.Addr
		ports := strings.Split(pod.Port+",",",")
		podinfo1.Port = ports[0:len(ports)-1]
		pods = append(pods, podinfo1)
	}
	return pods
}

func getSimplePodInfo(username string)[]podinfo{
	var pods []podinfo

	var podinfo1 podinfo

	user := models.GetUserInfo(username)
	podlist := models.GetUserPods(user.ID)
	for i := range podlist{
		pod := podlist[i]
		podinfo1.Name = pod.Name
		podinfo1.Podname = pod.Podname
		podinfo1.Type = "普通"
		podinfo1.Owner.Username = username
		podinfo1.Owner.Name = user.Name
		podinfo1.Image = pod.Image
		podinfo1.Image = podinfo1.Image[strings.Index(podinfo1.Image,"/")+1:]
		imagename := podinfo1.Image[0:strings.Index(podinfo1.Image,":")]
		image := models.GetImage(imagename)
		if(len(image.Acceptfile)!=0){
			podinfo1.Upload = uploadinfo{true,image.Acceptfile}
		}else{
			podinfo1.Upload = uploadinfo{false,""}
		}
		podinfo1.Createtime = pod.Createtime
		podinfo1.Addr = pod.Addr
		ports := strings.Split(pod.Port+",",",")
		podinfo1.Port = ports[0:len(ports)-1]

		pods = append(pods, podinfo1)
	}
	return pods
}

func ListAllPods(c *gin.Context){
	_ = checkToken(c)
	var data []podinfo
	if code != consts.SUCCESS {
		goto LAST
	}else{
		users := models.GetAllUsers()
		for i := range users{
			data = append(data,getSimplePodInfo(users[i].Username)...)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func upload(c *gin.Context,username string)(bool,string){
	//获取表单文件
	var data bool
	code = consts.SUCCESS
	formFile,header,err := c.Request.FormFile("file")
	if err != nil {
		//获取表单文件出错
		data = false
		return data,""
	}
	defer formFile.Close()

	//创建保存文件
	//cmd := exec.Command("mkdir",consts.FILEPATH+username+"/")
	//cmd.Output()
	os.Mkdir(consts.FILEPATH+username,0777)
	destFile,err := os.Create(consts.FILEPATH+username+"/"+header.Filename)
	if err!=nil{
		//创建文件出错
		code = consts.ERROR_FILE_CREATE
		data = false
		return data,""
	}
	defer destFile.Close()

	//读取表单文件，写入保存文件
	_,err = io.Copy(destFile,formFile)
	if err != nil{
		//写入文件出错
		code = consts.ERROR_FILE_WRITE
		data = false
		return data,""
	}
	data = true
	return data,header.Filename
}

func UploadFile(c *gin.Context){
	claims := checkToken(c)
	var (
		data bool
		filename string
	)
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username

		username2 := c.Param("username")
		if username2 != "" {
			username = username2
		}
		data,filename = upload(c,username)
		podname := c.PostForm("pod")
		image := c.PostForm("image")
		if(!data){
			//上传失败
			logger.Error("上传文件失败",zap.String("reason",consts.GetMsg(code)),zap.String("user",username),zap.String("type","上传"))
			goto LAST
		}else{
			data = false
			containerName := GetContainerName(podname)
			_,err := util.GetPodItf(username,consts.NORMAL).Get(podname,meta_v1.GetOptions{})
			if err != nil{
				code = consts.ERROR_POD_NOT_EXIST
				goto LAST
			}
			pod := kubectl.Pod(podname,username,containerName)
			dest := consts.FILEPATH+username+"/"
			if(strings.Contains(image,"tomcat")){
				err = pod.CopyToPod(dest+filename,"/usr/local/tomcat/webapps/"+filename)
				if err != nil{
					code = consts.ERROR_FILE_COPY
					goto LAST
				}
			}else{
				r,_ := zip.OpenReader(dest+filename)
				defer r.Close()

				for _,file:= range r.File{
					rc,_ := file.Open()
					fname := dest + file.Name
					os.MkdirAll(getDir(fname,"/"),0755)
					w,_ := os.Create(fname)
					io.Copy(w,rc)
					w.Close()
					rc.Close()
				}
				dir := dest + getDir(filename,".")
				err = pod.CopyToPod("/root/shells/mv.sh","/root")
				if err != nil{
					code = consts.ERROR_FILE_SHELL_COPY
					goto LAST
				}
				err = pod.CopyToPod("/root/shells/dotnet.sh","/root")
				if err != nil{
					code = consts.ERROR_FILE_SHELL_COPY
					goto LAST
				}
				err = pod.CopyToPod("/root/shells/rm.sh","/root")
				if err != nil{
					code = consts.ERROR_FILE_SHELL_COPY
					goto LAST
				}
				if(strings.Contains(image,":")){
					image = image[0:strings.Index(image,":")]
				}
				img := models.GetImage(image)
				err = pod.Exec([]string{"bash","/root/rm.sh",img.Dupath})
				if err != nil{
					code = consts.ERROR_FILE_SHELL_RUN
					goto LAST
				}
				err = pod.CopyToPod(dir,img.Dupath)
				if err != nil{
					code = consts.ERROR_FILE_COPY
					goto LAST
				}
				if(strings.Contains(image,"aspnet")){
					go pod.Exec([]string{"bash","/root/dotnet.sh",img.Dupath})
				}else{
					err = pod.Exec([]string{"bash","/root/mv.sh",img.Dupath})
				}
				if err != nil{
					code = consts.ERROR_FILE_SHELL_RUN
					goto LAST
				}
			}
			os.RemoveAll(consts.FILEPATH+username+"/")
			data = true
			logger.Info("上传文件成功",zap.String("user",username),zap.String("file",filename),zap.String("image",image),zap.String("type","上传"))
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func getDir(path string,split string) string{
	return path[0:strings.LastIndex(path, split)]
}

func SyncPodsToDB(c *gin.Context){
	users := models.GetAllUsers()
	for i := range users{
		username := users[i].Username
		podItf := util.GetPodItf(username,consts.CLASS)
		svcItf := util.GetSvcItf(username,consts.CLASS)
		if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
			for _, pod := range podList.Items {
				podName := pod.Name
				idx := strings.LastIndex(podName,username)
				svcname := podName[:idx] + username
				svc,err := svcItf.Get(svcname,meta_v1.GetOptions{})
				image := pod.Spec.Containers[0].Image
				image = image[strings.Index(image,"/")+1:]
				createtime := pod.CreationTimestamp.Unix()
				addr := consts.HOST
				var ports []string

				if svc != nil && err == nil{
					size := len(svc.Spec.Ports)
					ports = []string{}
					svcports := svc.Spec.Ports
					for i:=0 ; i < size ; i++{
						ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
					}
				}
				models.AddPodInfo(podName,svcname,image,addr,ports,createtime,username)
			}
		}
	}
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",true)
}

func SyncClassPodsToDB(c *gin.Context){
	classes := models.GetAllClasses()
	for i := range classes{
		class := classes[i]
		podItf := util.GetPodItf(class.Classen,consts.CLASS)
		svcItf := util.GetSvcItf(class.Classen,consts.CLASS)
		if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
			for _, pod := range podList.Items {
				podname := pod.Name
				container := GetContainerName(podname)
				tmp := container[0:strings.LastIndex(container,"-")]
				username := tmp[strings.LastIndex(tmp,"-")+1:]
				image := pod.Spec.Containers[0].Image
				image = image[strings.Index(image,"/")+1:]
				svc,_ := svcItf.Get(GetContainerName(podname),meta_v1.GetOptions{})
				size := len(svc.Spec.Ports)
				ports := []string{}
				svcports := svc.Spec.Ports
				for i:=0 ; i < size ; i++{
					ports = append(ports,fmt.Sprint(svcports[i].Port) + ":" + fmt.Sprint(svcports[i].NodePort))
				}
				models.CreateClassPod(username,podname,class.Classen,image,consts.HOST,ports,pod.GetCreationTimestamp().Unix())
			}
		}
	}
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",true)
}