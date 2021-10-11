package models

type Podinfo struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Podname			string		`json:"podname"`
	Name			string		`json:"name"`
	Image			string		`json:"image"`
	Addr			string		`json:"addr"`
	Port			string		`json:"port"`
	Createtime		int64		`json:"createtime"`
	Userid			int			`json:"userid"`
}

func IsPodExist(name string)bool{
	var podinfo Podinfo
	db.Table("podinfo").Select("*").Where("name = ?",name).First(&podinfo)
	return podinfo.ID!=0
}

func GetUserPods(userid int)[]Podinfo{
	var podinfos []Podinfo
	db.Table("podinfo").Select("*").Where("userid = ?",userid).Find(&podinfos)
	return podinfos
}

func PodTerminating(podname string){
	var pi Podinfo
	db.Table("podinfo").Select("*").Where("podname = ?",podname).First(&pi)
	pi.Port = ""
	db.Table("podinfo").Save(pi)
}

func AddPodInfo(podname string,name string,image string, addr string,port []string,createtime int64,username string){
	var pi Podinfo
	db.Table("podinfo").Select("*").Where("podname = ?",podname).First(&pi)
	pi.Userid = GetUserInfo(username).ID
	pi.Podname = podname
	pi.Name = name
	pi.Image = image
	pi.Addr = addr
	pi.Createtime = createtime
	var ports string
	ports = port[0]
	if(len(port)>1){
		for i:=1;i<len(port);i++{
			ports = ports + "," + port[i]
		}
	}
	pi.Port = ports
	db.Table("podinfo").Save(&pi)
}

func DeletePod(podname string){
	db.Table("podinfo").Where("podname = ?",podname).Delete(Podinfo{})
}

func DeleteUserPods(userid int){
	db.Table("podinfo").Where("userid = ?",userid).Delete(Podinfo{})
}
