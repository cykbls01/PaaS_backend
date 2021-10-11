package api

import (
	"Cloud/models"
	"Cloud/pkg/consts"
	"Cloud/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wzyonggege/logger"
	"go.uber.org/zap"
	"k8s.io/api/apps/v1beta1"
	v12 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strconv"
	"strings"
	"time"
)

func IsClassExist(c *gin.Context){
	code = consts.SUCCESS
	_ = checkToken(c)
	var data bool
	var info struct{
		Classen			string		`json:"classen"`
	}
	if code != consts.SUCCESS{
		goto LAST
	}else{
		c.BindJSON(&info)
		data = models.IsClassExist(info.Classen)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func CreateClassInfo(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classch			string				`json:"classch"`
			Classen			string				`json:"classen"`
			Classprofile	string				`json:"classprofile"`
			Endtime			int64				`json:"endtime"`
		}
		_ = c.BindJSON(&info)
		username := claims.Username
		flag := util.CreateNamespaceIfNotExist(info.Classen,consts.CLASS)
		if(!flag){
			logger.Error("创建课程失败，请检查课程英文名",zap.String("user",username),zap.String("type","创建课程"))
			data = false
			code = consts.ERROR_CLASS_EXIST
			goto LAST
		}
		models.AddClassInfo(username,info.Classen,info.Classch,info.Classprofile,info.Endtime)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func UpdateClassInfo(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classch			string				`json:"classch"`
			Classen			string				`json:"classen"`
			Classprofile	string				`json:"classprofile"`
			Endtime			int64				`json:"endtime"`
		}
		_ = c.BindJSON(&info)
		username := claims.Username
		models.AddClassInfo(username,info.Classen,info.Classch,info.Classprofile,info.Endtime)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetClassPods(c *gin.Context){
	_ = checkToken(c)
	var data []podinfo
	if code != consts.SUCCESS{
		goto LAST
	}else {
		var info struct {
			Classid		int		`json:"classid"`
		}
		_ = c.BindJSON(&info)
		class := models.GetClassByID(info.Classid)
		data = getSimpleClassPod(class.Classen)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func DeleteClassPods(c *gin.Context){
	code = consts.SUCCESS
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else {
		var info struct {
			Classid		int			`json:"classid"`
			Podname		[]string	`json:"podname"`
		}
		_ = c.Bind(&info)
		username := claims.Username
		class := models.GetClassByID(info.Classid)
		deployItf := util.GetDeployItf(class.Classen,consts.CLASS)

		svcItf := util.GetSvcItf(class.Classen,consts.CLASS)
		podItf :=util.GetPodItf(class.Classen,consts.CLASS)

		for i := 0 ; i < len(info.Podname);i++{
			podname := info.Podname[i]
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
			name := GetContainerName(podname)

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
			models.DeleteClassPod(podname)
			logger.Info("删除容器 "+name+" 成功",zap.String("user",username),zap.String("type","删除"))
		}
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	//GetClassPods(c)
	return
}

func DeleteUserClassPods(username string){
	classpods := models.GetClassPods(username)
	for i := range classpods {
		classpod := classpods[i]
		podname  := classpod.Podname
		containerName := GetContainerName(podname)
		deployItf := util.GetDeployItf(classpod.Class,consts.CLASS)
		svcItf := util.GetSvcItf(classpod.Class,consts.CLASS)
		deleteDeploy(deployItf,containerName)
		deleteSvc(svcItf,containerName)
	}
}

type cinfo struct {
	Classid      	int    		`gorm:"primary_key" json:"classid"`
	Userid       	int    		`json:"userid"`
	Owner			owner		`json:"owner"`
	Classen      	string		`json:"classen"`
	Classch      	string 		`json:"classch"`
	Classprofile 	string		`json:"classprofile"`
	Endtime			int64		`json:"endtime"`
}

func GetClassInfo(c *gin.Context){
	claims := checkToken(c)
	var data []cinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		usr := models.GetUserInfo(username)
		var classes []models.Classinfo
		if(usr.Authority == consts.AUTH_ADMIN){
			classes = models.GetAllClasses()
		}else{
			classes = models.GetClassInfo(username)
		}
		for i := range classes{
			class := classes[i]
			var tmpclass cinfo
			tmpclass.Classid = class.Classid
			tmpclass.Userid = class.Userid
			user := models.GetUserInfoById(tmpclass.Userid)
			tmpclass.Owner = owner{user.Username,user.Name}
			tmpclass.Classen = class.Classen
			tmpclass.Classch = class.Classch
			tmpclass.Classprofile = class.Classprofile
			tmpclass.Endtime = class.Endtime
			data = append(data, tmpclass)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func AddUsertoClass(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classid		int			`json:"classid"`
			Users		[]string	`json:"users"`
		}
		c.Bind(&info)
		if(models.IsClassEnded(info.Classid)){
			data = false
			code = consts.ERROR_CLASS_END
			goto LAST
		}
		for i := range info.Users{
			if(!models.IsUserExist(info.Users[i])){
				code = consts.ERROR_CLASS_USER_NOT_EXIST
				continue
			}
			user := models.GetUserInfo(info.Users[i])
			if(models.IsUserExistInClass(user.ID,info.Classid)){
				continue
			}
			models.AddUsertoClass(user.ID,info.Classid)
		}
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetClassUserList(c *gin.Context){
	_ = checkToken(c)
	type retinfo struct{
		Username	string		`json:"username"`
		Name		string		`json:"name"`
	}
	var data []retinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classid		int			`json:"classid"`
		}
		c.Bind(&info)
		cas := models.GetUsersinClass(info.Classid)
		for i := range cas{
			user := models.GetUserInfoById(cas[i].Userid)
			var tmpinfo retinfo
			tmpinfo.Name = user.Name
			tmpinfo.Username = user.Username
			data = append(data, tmpinfo)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func DeleteUserFromClass(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classid		int			`json:"classid"`
			Users		[]string	`json:"users"`
		}
		c.Bind(&info)
		for i := range info.Users{
			user := models.GetUserInfo(info.Users[i])
			models.DeleteUserfromClass(user.ID,info.Classid)
		}
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

//func BatchCreate(c *gin.Context){
//	_ = checkToken(c)
//	var data bool
//	if code != consts.SUCCESS{
//		goto LAST
//	}else{
//		var info struct{
//			Classid			int 			`json:"classid"`
//			Name			string			`json:"name"`
//			Image			[]string		`json:"image"`
//			Env				[]Children		`json:"env"`
//			Limit		map[string]int64	`json:"limit"`
//			Port 			[]int			`json:"port"`
//		}
//		c.Bind(&info)
//		imagepiece := info.Image
//		env := info.Env
//		port := info.Port
//		image := imagepiece[0] + ":" + imagepiece[1]
//		name := info.Name
//		class := models.GetClassByID(info.Classid)
//		deployItf := util.GetDeployItf(class.Classen,consts.CLASS)
//		svcItf := util.GetSvcItf(class.Classen,consts.CLASS)
//		users := models.GetUsersinClass(info.Classid)
//		for i:=range users{
//			user := models.GetUserInfoById(users[i].Userid).Username
//			tname := name + "-" + user + "-" + class.Classen
//			createDeploy1(user,deployItf,tname,image,env,info.Limit,consts.CLASS)
//			if code != consts.SUCCESS{
//				logger.Error("创建deployment "+tname+" 失败",zap.String("reason",consts.GetMsg(code)),zap.String("user",user),zap.String("type","批量创建"))
//				data = false
//				goto LAST
//			}
//			createSvc(svcItf,tname,image,port)
//			if code != consts.SUCCESS{
//				logger.Error("创建service "+tname+" 失败",zap.String("reason",consts.GetMsg(code)),zap.String("user",user),zap.String("type","批量创建"))
//				data = false
//				deleteDeploy(deployItf,name)
//				goto LAST
//			}
//			logger.Info("创建容器 "+tname+" 成功",zap.String("user",user),zap.String("type","批量创建"))
//		}
//		podItf := util.GetPodItf(class.Classen,consts.CLASS)
//		if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
//			for _, pod := range podList.Items {
//				podname := pod.Name
//				container := GetContainerName(podname)
//				tmp := container[0:strings.LastIndex(container,"-")]
//				username := tmp[strings.LastIndex(tmp,"-")+1:]
//				models.CreateClassPod(username,podname,class.Classen)
//			}
//		}
//		data = true
//	}
//LAST:
//	c.Set("code",code)
//	c.Set("msg",consts.GetMsg(code))
//	c.Set("data",data)
//	return
//}

func BatchCreateViaExistedPod(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		var info struct{
			Classid		int			`json:"classid"`
			Podname		string		`json:"podname"`
		}
		_ = c.Bind(&info)
		if(models.IsClassEnded(info.Classid)){
			data = false
			code = consts.ERROR_CLASS_END
			goto LAST
		}
		deployitf := util.GetDeployItf(username,consts.CLASS)
		deployname := GetContainerName(info.Podname)
		deploy,err := deployitf.Get(deployname,meta_v1.GetOptions{})
		if(err != nil){
			logger.Error("deployment不存在",zap.String("user",username),zap.String("type","批量创建"))
			data = false
			goto LAST
		}
		svcitf := util.GetSvcItf(username,consts.CLASS)
		svc,err := svcitf.Get(deployname,meta_v1.GetOptions{})
		if(err != nil){
			logger.Error("service不存在",zap.String("user",username),zap.String("type","批量创建"))
			data = false
			goto LAST
		}
		class := models.GetClassByID(info.Classid)
		deployitf1 := util.GetDeployItf(class.Classen,consts.CLASS)
		svcitf1 := util.GetSvcItf(class.Classen,consts.CLASS)
		users := models.GetUsersinClass(info.Classid)
		for i:=range users{
			user := models.GetUserInfoById(users[i].Userid).Username
			name := deploy.Name
			name = name[0:strings.LastIndex(name,"-")]
			name = name + "-" + user + "-" + class.Classen
			if(models.IsClassPodExist(models.GetUserInfo(user).ID,class.Classen,name)){
				continue
			}
			tmp := make(map[string]string)
			tmp["app"] = name
			tmpdeploy := v1beta1.Deployment{}
			tmpdeploy.APIVersion = "extensions/v1beta1"
			tmpdeploy.Kind = "Deployment"
			tmpdeploy.SetName(name)
			var container v12.Container
			container.Name = name
			container.Image = deploy.Spec.Template.Spec.Containers[0].Image
			container.Resources = deploy.Spec.Template.Spec.Containers[0].Resources
			container.Env = deploy.Spec.Template.Spec.Containers[0].Env
			tmpdeploy.Spec = v1beta1.DeploymentSpec{
				Replicas: deploy.Spec.Replicas,
				Template: v12.PodTemplateSpec{
					ObjectMeta: meta_v1.ObjectMeta{
						Labels: tmp,
					},
					Spec: v12.PodSpec{
						Containers: []v12.Container{container},
					},
				},
			}
			deployitf1.Create(&tmpdeploy)
			tmpsvc := v12.Service{}
			tmpsvc.APIVersion = "v1"
			tmpsvc.Kind = "Service"
			tmpsvc.SetName(name)
			svcport := []v12.ServicePort{}
			port := svc.Spec.Ports
			for i:=0 ; i < len(port) ; i++{
				svcport=append(svcport,v12.ServicePort{Port: port[i].Port,TargetPort: port[i].TargetPort,Name: port[i].Name})
			}
			tmpsvc.Spec = v12.ServiceSpec{
				Type: v12.ServiceTypeNodePort,
				Ports: svcport,
				Selector: tmp,
			}
			svcitf1.Create(&tmpsvc)
		}
		podItf := util.GetPodItf(class.Classen,consts.CLASS)
		svcItf := util.GetSvcItf(class.Classen,consts.CLASS)
		if podList, err := podItf.List(meta_v1.ListOptions{}); err == nil {
			for _, pod := range podList.Items {
				podname := pod.Name
				if(models.IsClassPodExist2(podname)){
					continue
				}
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
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func DeleteClass(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classid		int		`json:"classid"`
		}
		c.BindJSON(&info)
		classinfo := models.GetClassByID(info.Classid)
		nsItf := util.GetNSItf()
		err = nsItf.Delete(classinfo.Classen,&meta_v1.DeleteOptions{})
		if err != nil{
			logger.Error("删除命名空间 "+ classinfo.Classen + " 失败",zap.String("type","删除课程"),zap.String("user",claims.Username))
			code = consts.ERROR_DELETE_NS
			data = false
			goto LAST
		}
		models.DeleteTasksInClass(classinfo.Classid)
		models.DeleteClassAuths(classinfo.Classid)
		models.DeleteClass(classinfo.Classen)
		logger.Info("删除命名空间" + classinfo.Classen+" 成功",zap.String("type","删除课程"),zap.String("user",claims.Username))
		models.DeleteClassInfo(info.Classid)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

//task apis for teachers
func AddTask(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		classid,_ := strconv.Atoi(c.PostForm("classid"))
		starttime,_ := strconv.ParseInt(c.PostForm("starttime"),10,64)
		endtime,_ := strconv.ParseInt(c.PostForm("endtime"),10,64)
		taskname := c.PostForm("taskname")
		taskgoal := c.PostForm("taskgoal")
		taskstep := c.PostForm("taskstep")
		if(models.IsTaskExistInClass(classid,taskname)){
			data = false
			code = consts.ERROR_CLASS_TASK_EXIST
			goto LAST
		}
		models.AddClassTask(classid,starttime,endtime,taskname,taskgoal,taskstep,"")
		flag,filename := upload(c,"class-"+strconv.Itoa(classid))
		if(flag){
			filename = consts.FILEPATH+"class-"+strconv.Itoa(classid)+"/"+filename
		}
		models.AddClassTask(classid,starttime,endtime,taskname,taskgoal,taskstep,filename)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func UpdateTask(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		classid,_ := strconv.Atoi(c.PostForm("classid"))
		starttime,_ := strconv.ParseInt(c.PostForm("starttime"),10,64)
		endtime,_ := strconv.ParseInt(c.PostForm("endtime"),10,64)
		taskname := c.PostForm("taskname")
		taskgoal := c.PostForm("taskgoal")
		taskstep := c.PostForm("taskstep")
		flag,filename := upload(c,"class-"+strconv.Itoa(classid))
		if(flag){
			filename = consts.FILEPATH+"class-"+strconv.Itoa(classid)+"/"+filename
		}
		models.AddClassTask(classid,starttime,endtime,taskname,taskgoal,taskstep,filename)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetTasksInClass(c *gin.Context){
	_ = checkToken(c)
	type retinfo struct{
		Taskid    int    `json:"taskid"`
		Classid   int    `json:"classid"`
		Starttime int64  `json:"starttime"`
		Endtime   int64  `json:"endtime"`
		Taskname  string `json:"taskname"`
		Taskgoal  string `json:"taskgoal"`
		Taskstep  string `json:"taskstep"`
		Totnum    int    `json:"totnum"`
		Submitnum int    `json:"submitnum"`
	}
	var data []retinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Classid			int			`json:"classid"`
		}
		c.BindJSON(&info)
		tmpdata := models.GetTasksByClassid(info.Classid)
		for i := range tmpdata{
			task := tmpdata[i]
			totnum := models.CountClassMember(task.Classid)
			finishnum := models.CountTaskSubmits(task.ID)
			rettask := retinfo{
				Taskid:    task.ID,
				Classid:   task.Classid,
				Starttime: task.Starttime,
				Endtime:   task.Endtime,
				Taskname:  task.Taskname,
				Taskgoal:  task.Taskgoal,
				Taskstep:  task.Taskstep,
				Totnum:    totnum,
				Submitnum: finishnum,
			}
			data = append(data, rettask)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetSubmitsInTask(c *gin.Context){
	_ = checkToken(c)
	type retinfo struct {
		Username	string		`json:"username"`
		Name		string		`json:"name"`
		Taskfile	taskfile	`json:"filename"`
		Submit		bool		`json:"submit"`
	}
	var data []retinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct {
			Taskid		int		`json:"taskid"`
		}
		c.Bind(&info)
		task := models.GetTaskByid(info.Taskid)
		class := models.GetClassByID(task.Classid)
		users := models.GetUsersinClass(class.Classid)
		for i := range users{
			user := models.GetUserInfoById(users[i].Userid)
			submit := models.GetTaskSubmitByUserid(user.ID,task.ID)
			tmpret := retinfo{
				Username:   user.Username,
				Name:       user.Name,
				Taskfile: taskfile{
					Filepath:   submit.Filepath,
					Filename:   getFilename(submit.Filepath),
					Submittime: submit.Submittime,
				},
				Submit:     submit.ID != 0,
			}
			data = append(data, tmpret)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func DeleteTask(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Taskid			int			`json:"taskid"`
		}
		c.BindJSON(&info)
		task := models.GetTaskByid(info.Taskid)
		os.RemoveAll(consts.FILEPATH+"class-"+strconv.Itoa(task.Classid)+"/task-"+strconv.Itoa(task.ID))
		models.DeleteTaskByid(info.Taskid)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

//func GetTasksForStu(c *gin.Context){
//	claims := checkToken(c)
//	var data []
//	if code != consts.SUCCESS{
//		goto LAST
//	}else{
//		username := claims.Username
//		user := models.GetUserInfo(username)
//	}
//LAST:
//	c.Set("code",code)
//	c.Set("msg",consts.GetMsg(code))
//	c.Set("data",data)
//	return
//}

type taskfile struct{
	Filepath	string	`json:"filepath"`
	Filename	string	`json:"filename"`
	Submittime	int64	`json:"submittime"`
}

func getFilename(filepath string)string{
	return filepath[strings.LastIndex(filepath,"/")+1:]
}

//task apis for students
func SubmitTask(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		taskid,_ := strconv.Atoi(c.PostForm("taskid"))
		task := models.GetTaskByid(taskid)
		class := models.GetClassByID(task.Classid)
		user := models.GetUserInfo(claims.Username)
		if(time.Now().Unix()>task.Endtime){
			data = false
			code = consts.ERROR_CLASS_TASK_END
			goto LAST
		}
		flag,filename := upload(c,"class-"+strconv.Itoa(class.Classid)+"/task-"+strconv.Itoa(taskid))
		if(!flag){
			data = false
			code = consts.ERROR_CLASS_TASK_SUBMIT
			goto LAST
		}
		filename = consts.FILEPATH + "class-"+strconv.Itoa(class.Classid)+"/task-"+strconv.Itoa(taskid) + "/" + filename
		models.TaskSubmit(user.ID,taskid,filename)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetExistingTasks(c *gin.Context){
	claims := checkToken(c)
	type retinfo struct{
		Taskid		int		`json:"taskid"`
		Taskname	string	`json:"taskname"`
		Classch		string	`json:"classch"`
		Teacher		string	`json:"teacher"`
		Starttime	int64	`json:"starttime"`
		Endtime		int64	`json:"endtime"`
		Submit		bool	`json:"submit"`
	}
	var data []retinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		user := models.GetUserInfo(username)
		cas := models.GetUsersClass(user.ID)
		var tasks []models.Classtask
		for i := range cas{
			ca := cas[i]
			tasks = append(tasks, models.GetTasksByClassid(ca.Classid)...)
		}
		for i:= range tasks{
			task := tasks[i]
			class := models.GetClassByID(task.Classid)
			teacher := models.GetUserInfoById(class.Userid)
			flag := models.IsUserSubmited(task.ID,user.ID)
			tmpret := retinfo{
				Taskid:    task.ID,
				Taskname:  task.Taskname,
				Classch:   class.Classch,
				Teacher:   teacher.Name,
				Starttime: task.Starttime,
				Endtime:   task.Endtime,
				Submit:    flag,
			}
			data = append(data, tmpret)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

//task apis for all
func GetTaskDetails(c *gin.Context){
	claims := checkToken(c)
	type retinfo struct{
		Taskid 			int 		`json:"taskid"`
		Taskname		string		`json:"taskname"`
		Classch			string		`json:"classch"`
		Teacher			string		`json:"teacher"`
		Starttime		int64		`json:"starttime"`
		Endtime			int64		`json:"endtime"`
		Taskgoal		string		`json:"taskgoal"`
		Taskstep		string		`json:"taskstep"`
		Reffile			taskfile	`json:"reffile"`
		Totnum			int			`json:"totnum"`
		Submitnum		int			`json:"submitnum"`
		Isstudent		bool		`json:"isstudent"`
		Submit			bool		`json:"submit"`
		Submitfile		taskfile	`json:"submitfile"`
	}
	var data retinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Taskid 		int		`json:"taskid"`
		}
		c.Bind(&info)
		username := claims.Username
		user := models.GetUserInfo(username)
		task := models.GetTaskByid(info.Taskid)
		class := models.GetClassByID(task.Classid)
		data.Taskid = task.ID
		data.Taskname = task.Taskname
		data.Classch = class.Classch
		data.Teacher = models.GetUserInfoById(class.Userid).Name
		data.Starttime = task.Starttime
		data.Endtime = task.Endtime
		data.Taskgoal = task.Taskgoal
		data.Taskstep = task.Taskstep
		data.Reffile = taskfile{
			Filepath: task.Filepath,
			Filename: getFilename(task.Filepath),
		}
		data.Isstudent = models.IsUserExistInClass(user.ID,class.Classid)
		if(data.Isstudent){
			data.Submit = models.IsUserSubmited(task.ID,user.ID)
			if(data.Submit){
				ts := models.GetTaskSubmitByUserid(user.ID,task.ID)
				data.Submitfile = taskfile{
					Filepath: ts.Filepath,
					Filename: getFilename(ts.Filepath),
					Submittime: ts.Submittime,
				}
			}
		}else{
			data.Totnum = models.CountClassMember(class.Classid)
			data.Submitnum = models.CountTaskSubmits(task.ID)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func DealWithEndedClass(){
	classes := models.GetAllClasses()
	for i := range classes{
		class := classes[i]
		if(time.Now().Unix() > class.Endtime){
			deployItf := util.GetDeployItf(class.Classen,consts.CLASS)
			svcItf := util.GetSvcItf(class.Classen,consts.CLASS)
			pods := models.GetClassPodsByClassen(class.Classen)
			for j := range pods{
				podname := pods[j].Podname
				containername := GetContainerName(podname)
				deleteDeploy(deployItf,containername)
				deleteSvc(svcItf,containername)
				models.DeleteClassPod(podname)
			}
		}
	}
}
