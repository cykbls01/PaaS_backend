package api

import (
	"Cloud/models"
	"Cloud/pkg/consts"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wzyonggege/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func getRepos(c *gin.Context){
	code = consts.SUCCESS
	client := &http.Client{}
	resp, _ := client.Get("http://"+consts.HOST+":5000/v2/_catalog")
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	var repos repository
	var repoinfos []repoinfo
	err := json.Unmarshal(body,&repos)
	if err != nil{
		code = consts.ERROR
		goto LAST
	}
	for i:=0;i<len(repos.Names);i++{
		reponame := repos.Names[i]
		resp2,_ := client.Get("http://10.251.0.14:5000/v2/" + reponame +"/tags/list")
		defer resp2.Body.Close()
		var svcvers svcversion
		body2,_ := ioutil.ReadAll(resp2.Body)
		err := json.Unmarshal(body2,&svcvers)
		if err != nil{
			code = consts.ERROR
			goto LAST
		}
		var repo repoinfo
		for j:=0;j<len(svcvers.Tags);j++{
			var child Children
			child.Value  = svcvers.Tags[j]
			child.Label = child.Value
			repo.Childs = append(repo.Childs,child)
		}
		repo.Value = svcvers.Name
		repo.Label = repo.Value
		repoinfos = append(repoinfos,repo)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",repoinfos)
}

func getShareRepos(c *gin.Context){
	code = consts.SUCCESS
	client := &http.Client{}

	var repos repository
	var repoinfos []repoinfo

	images := models.GetAllImage()

	for i:= range images{
		if(len(images[i].Mntpath) != 0){
			repos.Names = append(repos.Names,images[i].Image)
		}
	}

	for i:=0;i<len(repos.Names);i++{
		reponame := repos.Names[i]
		resp2,_ := client.Get("http://10.251.0.14:5000/v2/" + reponame +"/tags/list")
		defer resp2.Body.Close()
		var svcvers svcversion
		body2,_ := ioutil.ReadAll(resp2.Body)
		err := json.Unmarshal(body2,&svcvers)
		if err != nil{
			code = consts.ERROR
			goto LAST
		}
		var repo repoinfo
		for j:=0;j<len(svcvers.Tags);j++{
			var child Children
			child.Value  = svcvers.Tags[j]
			child.Label = child.Value
			repo.Childs = append(repo.Childs,child)
		}
		repo.Value = svcvers.Name
		repo.Label = repo.Value
		repoinfos = append(repoinfos,repo)
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",repoinfos)
}

func GetImage(c *gin.Context){
	claims := checkToken(c)
	username := claims.Username
	if(username == "admin"){
		getShareRepos(c)
	}else{
		getRepos(c)
	}
}

func GetImageInfo(c *gin.Context){
	_ = checkToken(c)
	type rs struct{
		LowComputing		resourceinfo		`json:"lowComputing"`
		MediumComputing		resourceinfo		`json:"mediumComputing"`
		HighComputing		resourceinfo		`json:"highComputing"`
		LowMemory			resourceinfo		`json:"lowMemory"`
		MediumMemory		resourceinfo		`json:"mediumMemory"`
		HighMemory			resourceinfo		`json:"highMemory"`
		LowStorage			resourceinfo		`json:"lowStorage"`
		MediumStorage		resourceinfo		`json:"mediumStorage"`
		HighStorage			resourceinfo		`json:"highStorage"`
	}
	var data struct{
		DefaultConfiguration		rs			`json:"defaultConfiguration"`
		DefaultPort					[]int		`json:"defaultPort"`
		DefaultEnv					[]Children	`json:"defaultEnv"`
	}
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct{
			Image		string		`json:"image"`
		}
		c.BindJSON(&info)
		image := models.GetImage(info.Image)
		var resource rs
		var env []Children
		resource.LowComputing = resourceinfo{int64(image.Lowcpu) + 100 , int64(image.Lowmem), int64(image.Lowsto)}
		resource.MediumComputing = resourceinfo{int64(float64(image.Lowcpu)*1.5)+150,int64(float64(image.Lowmem)*1.5),int64(float64(image.Lowsto)*1.5)}
		resource.HighComputing = resourceinfo{int64(image.Lowcpu)*2 + 200 , int64(image.Lowmem)*2, int64(image.Lowsto)*2}
		resource.LowMemory = resourceinfo{int64(image.Lowcpu) , int64(image.Lowmem)+100, int64(image.Lowsto)}
		resource.MediumMemory = resourceinfo{int64(float64(image.Lowcpu)*1.5),int64(float64(image.Lowmem)*1.5)+150,int64(float64(image.Lowsto)*1.5)}
		resource.HighMemory = resourceinfo{int64(image.Lowcpu)*2 , int64(image.Lowmem)*2 + 200, int64(image.Lowsto)*2}
		resource.LowStorage = resourceinfo{int64(image.Lowcpu)  , int64(image.Lowmem), int64(image.Lowsto)+ 100}
		resource.MediumStorage = resourceinfo{int64(float64(image.Lowcpu)*1.5),int64(float64(image.Lowmem)*1.5),int64(float64(image.Lowsto)*1.5)+150}
		resource.HighStorage = resourceinfo{int64(image.Lowcpu)*2 , int64(image.Lowmem)*2, int64(image.Lowsto)*2 + 200}
		data.DefaultConfiguration = resource
		tmpstr := strings.Split(image.Port,",")
		var tmpPort []int
		for i := range tmpstr{
			tmp,_ := strconv.Atoi(tmpstr[i])
			tmpPort = append(tmpPort,tmp)
		}
		data.DefaultPort = tmpPort
		tmpEnv := strings.Split(image.Env,"#;")
		for i:=range tmpEnv{
			if i == len(tmpEnv) - 1{
				break
			}
			str := tmpEnv[i]
			key := str[2:strings.Index(str,"#,v:")]
			value := str[strings.Index(str,"#,v:")+4:strings.Index(str,"#,c:")]
			comment := str[strings.Index(str,"#,c:")+4:]
			env = append(env,Children{value, key,comment})
		}
		data.DefaultEnv = env
		if len(data.DefaultEnv) == 0{
			data.DefaultEnv = []Children{}
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func AddNewImage(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct {
			Url				string		`json:"url"`
			Tag				string		`json:"tag"`
			Image			string		`json:"image"`
			Lowcpu			int			`json:"lowcpu"`
			Lowmem			int			`json:"lowmem"`
			Lowsto			int			`json:"lowsto"`
			Port			[]int		`json:"port"`
			Env				[]Children	`json:"env"`
			Mntpath			string		`json:"mntpath"`
			Dupath			string		`json:"dupath"`
			Acceptfile		string		`json:"acceptfile"`
		}
		_ = c.Bind(&info)
		cmd := exec.Command("/bin/bash","/root/shells/image.sh",info.Url,info.Tag)
		_,err = cmd.Output()
		if(err != nil){
			data = false
			code = consts.ERROR_IMAGE_PUSH
			logger.Error("添加镜像失败",zap.String("type","添加镜像"),zap.String("user","admin"),zap.String("reason","请检查输入内容"))
			goto LAST
		}
		var env string
		for i := range info.Env{
			env = env + "k:" + info.Env[i].Label + "#,v:" + info.Env[i].Value+"#,c:" + info.Env[i].Comment + "#;"
		}
		if(strings.Contains(info.Image,":")){
			info.Image = info.Image[0:strings.Index(info.Image,":")]
		}
		var port string
		for i := range info.Port{
			port = port + strconv.Itoa(info.Port[i])
			if i != len(info.Port)-1{
				port = port + ","
			}
		}
		models.CreateImage(info.Image,info.Lowcpu,info.Lowmem,info.Lowsto,port,env,info.Mntpath,info.Dupath,info.Acceptfile)
		logger.Info("添加镜像成功",zap.String("image",info.Image),zap.String("type","添加镜像"),zap.String("user","admin"))
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func GetAllImage(c *gin.Context){
	_ = checkToken(c)
	type imageinfo struct{
		Image			string		`json:"image"`
		Lowcpu			int			`json:"lowcpu"`
		Lowmem			int			`json:"lowmem"`
		Lowsto			int			`json:"lowsto"`
		Port			[]int		`json:"port"`
		Env				[]Children	`json:"env"`
		Mntpath			string		`json:"mntpath"`
		Dupath			string		`json:"dupath"`
	}
	var data []imageinfo
	if code != consts.SUCCESS{
		goto LAST
	}else{
		tmp := models.GetAllImage()
		for i := range tmp{
			tmpinfo := tmp[i]
			var tinfo imageinfo
			tinfo.Image = tmpinfo.Image
			tinfo.Lowcpu = tmpinfo.Lowcpu
			tinfo.Lowmem = tmpinfo.Lowmem
			tinfo.Lowsto = tmpinfo.Lowsto
			tinfo.Dupath = tmpinfo.Dupath
			tinfo.Mntpath = tmpinfo.Mntpath
			tmpstr := strings.Split(tmpinfo.Port,",")
			var tmpPort []int
			for i := range tmpstr{
				tmp,_ := strconv.Atoi(tmpstr[i])
				tmpPort = append(tmpPort,tmp)
			}
			tinfo.Port = tmpPort
			tmpEnv := strings.Split(tmpinfo.Env,"#;")
			for i:=range tmpEnv{
				if i == len(tmpEnv) - 1{
					break
				}
				str := tmpEnv[i]
				key := str[2:strings.Index(str,"#,v:")]
				value := str[strings.Index(str,"#,v:")+4:strings.Index(str,"#,c:")]
				comment := str[strings.Index(str,"#,c:")+4:]
				tinfo.Env = append(tinfo.Env,Children{value, key,comment})
			}
			data = append(data,tinfo)
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}

func UpdateImage(c *gin.Context){
	_ = checkToken(c)
	var data bool
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var info struct {
			Image			string		`json:"image"`
			Lowcpu			int			`json:"lowcpu"`
			Lowmem			int			`json:"lowmem"`
			Lowsto			int			`json:"lowsto"`
			Port			[]int		`json:"port"`
			Env				[]Children	`json:"env"`
			Mntpath			string		`json:"mntpath"`
			Dupath			string		`json:"dupath"`
			Acceptfile		string		`json:"acceptfile"`
		}
		_ = c.Bind(&info)
		var env string
		for i := range info.Env{
			env = env + "k:" + info.Env[i].Label + "#,v:" + info.Env[i].Value+"#,c:" + info.Env[i].Comment + "#;"
		}
		if(strings.Contains(info.Image,":")){
			info.Image = info.Image[0:strings.Index(info.Image,":")]
		}
		var port string
		for i := range info.Port{
			port = port + strconv.Itoa(info.Port[i])
			if i != len(info.Port)-1{
				port = port + ","
			}
		}
		data = true
		models.CreateImage(info.Image,info.Lowcpu,info.Lowmem,info.Lowsto,port,env,info.Mntpath,info.Dupath,info.Acceptfile)
		logger.Info("修改镜像信息成功",zap.String("image",info.Image),zap.String("type","修改镜像"),zap.String("user","admin"))
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}