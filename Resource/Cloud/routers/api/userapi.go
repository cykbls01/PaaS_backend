package api

import (
	"Cloud/models"
	"Cloud/pkg/consts"
	"Cloud/pkg/util"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/wzyonggege/logger"
	"go.uber.org/zap"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

func GetAuth(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := consts.INVALID_PARAMS
	var isExist int
	if ok {
		isExist = models.CheckAuth(username, password)
		if isExist > 0 {
			token, err := util.GenerateToken(username, password,isExist)
			if err != nil {
				code = consts.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				data["currentAuthority"] = consts.GetMsg(isExist)

				code = consts.SUCCESS
			}

		} else {
			code = consts.ERROR_AUTH
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" :  consts.GetMsg(code),
		"data" : data,
	})
}

func CreateUser(c *gin.Context){
	var user struct{
		Username	string		`json:"username"`
		Password	string		`json:"password"`
		Authority	int			`json:"authority"`
	}

	err := c.Bind(&user)
	if err != nil{
		c.String(http.StatusOK,"格式错误")
		return
	}
	models.CreateUser(user.Username,user.Password,user.Username,user.Authority)
	c.String(http.StatusOK,"创建成功")
}

func BatchCreateUser(c *gin.Context){
	claim := checkToken(c)
	var data bool
	if(code != consts.SUCCESS){
		goto LAST
	}else{
		type singleuser struct{
			Username		string		`json:"username"`
			Name			string		`json:"name"`
		}
		var user struct{
			Users		[]singleuser	`json:"users"`
			Type		string			`json:"type"`
		}
		err := c.Bind(&user)
		if err != nil{
			logger.Error("批量创建用户出错",zap.String("reason","格式出错"),zap.String("type","创建用户"),zap.String("user",claim.Username))
			code = consts.ERROR_USER_CREATE
			data = false
			goto LAST
		}
		for i := range user.Users{
			if(models.IsUserExist(user.Users[i].Username)){
				continue
			}
			var t int
			if user.Type == "student"{
				t = consts.AUTH_STUDENT
				models.CreateUser(user.Users[i].Username,"abcde-12345",user.Users[i].Name,t)
			}else if user.Type == "teacher"{
				t = consts.AUTH_TEACHER
				models.CreateUser(user.Users[i].Username,"12345-abcde",user.Users[i].Name,t)
			}else if user.Type == "admin"{
				t = consts.AUTH_ADMIN
				models.CreateUser(user.Users[i].Username,"12345-abcde",user.Users[i].Name,t)
			}
			util.CreateQuota(user.Users[i].Username)
			models.CreateUserQuota(models.GetUserInfo(user.Users[i].Username).ID)
		}
		logger.Info("批量创建用户成功",zap.String("type","创建用户"),zap.String("user",claim.Username))
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func UpdateUserQuota(userid int){
	var aq models.Auth_quota
	user := models.GetUserInfoById(userid)
	println(user.Username)
	data := GetResourceUsage(user.Username)
	aq.Userid = userid
	aq.Totcpu = data["total"].Requestcpu
	aq.Totmem = data["total"].Requestmem
	aq.Totstorage = data["total"].Requeststo
	aq.Usedcpu = data["used"].Requestcpu
	aq.Usedmem = data["used"].Requestmem
	aq.Usedstorage = data["used"].Requeststo
	println(aq.Totcpu)
	models.UpdateUserQuota(aq)
}

func DeleteUser(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if(code != consts.SUCCESS){
		goto LAST
	}else{
		var info struct{
			Username		string		`json:"username"`
		}
		c.BindJSON(&info)
		userid := models.GetUserInfo(info.Username).ID
		models.DeleteUserPods(userid)
		models.DeleteUserQuota(userid)
		DeleteUserClassPods(info.Username)
		models.DeleteUserClasspods(userid)
		models.DeleteUser(info.Username)
		nsItf := util.GetNSItf()
		nsItf.Delete(info.Username,&meta_v1.DeleteOptions{})
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetUserInfo(c *gin.Context){
	data := make(map[string]string)
	claims:=checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		username := claims.Username
		info := models.GetUserInfo(username)
		data["name"] = info.Name
		data["email"] = info.Email
		data["profile"] = info.Profile
		data["authority"] = consts.GetMsg(info.Authority)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func CheckPassword(c *gin.Context){
	claims := checkToken(c)
	var flag bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		var userinfo struct{
			Password	string		`json:"password"`
		}
		_ = c.Bind(&userinfo)
		flag = models.CheckPassword(username,userinfo.Password)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",flag)
	return
}

func UpdateUserInfo(c *gin.Context){
	claims := checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username
		var userinfo struct{
			Email 		string 		`json:"email"`
			Profile 	string 		`json:"profile"`
			Oldpasswd	string		`json:"oldPassword"`
			Newpasswd	string 		`json:"newPassword"`
		}
		_ = c.Bind(&userinfo)
		if userinfo.Oldpasswd != "" {
			flag := models.CheckPassword(username, userinfo.Oldpasswd)
			if !flag {
				data = false
				code = consts.ERROR
				goto LAST
			}
			models.UpdatePassword(username, userinfo.Newpasswd)
		}else{
			models.UpdateInfo(username, userinfo.Email, userinfo.Profile)
		}
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func checkToken(c *gin.Context)*util.Claims{
	code = consts.SUCCESS
	var claims *util.Claims
	token := c.Request.Header.Get("Authorization")
	if token == ""{
		code = consts.ERROR_AUTH_TOKEN
		logger.Error("token为空",zap.String("type","鉴权"))
		return nil
	}else {
		token = token[7:]
		claims,err = util.ParseToken(token)
		if err != nil{
			code = consts.ERROR_AUTH_TOKEN
			logger.Error("token无效",zap.String("type","鉴权"))
			return nil
		}
	}
	return claims
}

func GetResourceUsage(username string)map[string]resourceinfo{
	code = consts.SUCCESS
	data := make(map[string]resourceinfo)
	rqItf := util.GetClientset().CoreV1().ResourceQuotas(username)
	rq,err := rqItf.Get("resource-"+username,meta_v1.GetOptions{})
	if err != nil{
		//logger.Error("获取命名空间资源使用情况失败")
		code = consts.ERROR
		return nil
	}

	hardrcpu := rq.Status.Hard[v12.ResourceRequestsCPU]
	hardrmem := rq.Status.Hard[v12.ResourceRequestsMemory]
	hardrsto := rq.Status.Hard[v12.ResourceRequestsEphemeralStorage]

	usedrcpu := rq.Status.Used[v12.ResourceRequestsCPU]
	usedrmem := rq.Status.Used[v12.ResourceRequestsMemory]
	usedrsto := rq.Status.Used[v12.ResourceRequestsEphemeralStorage]

	var rshard,rsused resourceinfo
	rshard.Requestcpu = hardrcpu.MilliValue()
	rshard.Requestmem = hardrmem.ScaledValue(resource.Mega)
	rshard.Requeststo = hardrsto.ScaledValue(resource.Mega)

	rsused.Requestcpu = usedrcpu.MilliValue()
	rsused.Requestmem = usedrmem.ScaledValue(resource.Mega)
	rsused.Requeststo = usedrsto.ScaledValue(resource.Mega)

	if pvcList, err := util.GetPVCItf(username,consts.NORMAL).List(meta_v1.ListOptions{}); err == nil {
		for _, pvc := range pvcList.Items {
			tpvc := &pvc
			quantity := new(resource.Quantity)
			*quantity = tpvc.Spec.Resources.Requests["storage"]
			size := quantity.ScaledValue(resource.Mega)
			rsused.Requeststo += size
		}
	}

	data["total"] = rshard
	data["used"] = rsused

	return data
}

func GetUserUsage(c *gin.Context){
	_ = checkToken(c)
	var data []Useruse
	if code != consts.SUCCESS{
		goto LAST
	}else{
		//nsItf := util.GetNSItf()
		//nslist , _ :=nsItf.List(meta_v1.ListOptions{})
		users := models.GetAllUsers()
		//for _,item := range nslist.Items{
		for i := range users{
			//ns := item.GetName()
			ns := users[i].Username
			if ns == "default" || strings.Index(ns,"kube") >= 0{
				continue
			}
			var useruse Useruse
			useruse.Username = ns
			user := models.GetUserInfo(ns)
			useruse.Name = user.Name
			go util.CreateQuota(ns)
			switch user.Authority {
			case 0: useruse.Authority = "class"
			case 1: useruse.Authority = "student"
			case 2: useruse.Authority = "teacher"
			case 3: useruse.Authority = "admin"
			}
			tmpuser := models.GetUserInfo(ns)
			tmp := models.GetUserQuota(tmpuser.ID)
			useruse.Total.Requestcpu = tmp.Totcpu
			useruse.Total.Requestmem = tmp.Totmem
			useruse.Total.Requeststo = tmp.Totstorage
			useruse.Used.Requestcpu = tmp.Usedcpu
			useruse.Used.Requestmem = tmp.Usedmem
			data = append(data,useruse)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func SetUserAuth(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		data = true
		var info struct{
			Username	string		`json:"username"`
			Authority	string		`json:"authority"`
		}
		_ = c.Bind(&info)
		var authority int
		switch info.Authority {
		case "student" : authority = consts.AUTH_STUDENT
		case "teacher" : authority = consts.AUTH_TEACHER
		}
		models.UpdateAuthority(info.Username,authority)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func UpdateAllUserQuota(c *gin.Context){
	users := models.GetAllUsers()
	for i := range users{
		UpdateUserQuota(models.GetUserInfo(users[i].Username).ID)
	}
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",true)
}

func GetUserRS(c *gin.Context){
	code = consts.SUCCESS
	claims := checkToken(c)
	var data map[string]resourceinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		username := claims.Username

		username2 := c.Param("username")
		if username2 != "" {
			username = username2
		}

		data = GetResourceUsage(username)

	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func EditUserRS(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var infos struct{
			Username 	string 		`json:"username"`
			Requestcpu	int64		`json:"requestCPU"`
			Requestmem	int64		`json:"requestMem"`
			Requeststo	int64		`json:"requestStorage"`
		}
		_ = c.Bind(&infos)
		RQItf := util.GetRQItf(infos.Username,consts.NORMAL)
		RQ,err := RQItf.Get("resource-" + infos.Username,meta_v1.GetOptions{})
		if err != nil{
			code = consts.ERROR_RESOURCEQUOTA_GET
			logger.Error(consts.GetMsg(code),zap.String("user",infos.Username),zap.String("type","用户管理"))
			goto LAST
		}
		RQ.Spec = v12.ResourceQuotaSpec{
			Hard: v12.ResourceList{
				v12.ResourceLimitsCPU : resource.MustParse(fmt.Sprint(int64(float64(infos.Requestcpu) * 1.2)) + "m"),
				v12.ResourceRequestsCPU: resource.MustParse(fmt.Sprint(infos.Requestcpu) + "m"),
				v12.ResourceLimitsMemory : resource.MustParse(fmt.Sprint(int64(float64(infos.Requestmem) * 1.2)) + "M"),
				v12.ResourceRequestsMemory : resource.MustParse(fmt.Sprint(infos.Requestmem) + "M"),
				v12.ResourceLimitsEphemeralStorage: resource.MustParse(fmt.Sprint(int64(float64(infos.Requeststo) * 1.2)) + "M"),
				v12.ResourceRequestsEphemeralStorage: resource.MustParse(fmt.Sprint(infos.Requeststo) + "M"),
			},
		}
		_,err = RQItf.Update(RQ)
		if err != nil{
			code = consts.ERROR_RESOURCEQUOTA_UPDATE
			logger.Error(consts.GetMsg(code),zap.String("user",infos.Username),zap.String("type","用户管理"))
			goto LAST
		}
		UpdateUserQuota(models.GetUserInfo(infos.Username).ID)
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func ResetUser(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Username	string		`json:"username"`
		}
		_ = c.BindJSON(&info)
		models.UpdatePassword(info.Username,"123456")
		data = true
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}