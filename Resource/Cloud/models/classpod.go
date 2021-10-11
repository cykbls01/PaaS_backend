package models

type Classpod struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Userid	 		int			`json:"userid"`
	Podname			string		`json:"podname"`
	Class			string		`json:"class"`
	Image			string		`json:"image"`
	Addr			string		`json:"addr"`
	Port			string		`json:"port"`
	Createtime		int64		`json:"createtime"`
}

func IsClassPodExist(userid int,class string,name string)bool{
	var cp Classpod
	db.Table("classpod").Select("*").Where("userid = ? and class = ? and podname like ?",userid,class,name+"%").First(&cp)
	return cp.ID != 0
}

func IsClassPodExist2(podname string)bool{
	var cp Classpod
	db.Table("classpod").Select("*").Where("podname = ?",podname).First(&cp)
	return cp.ID != 0
}

func CreateClassPod(username string,podname string,class string,image string,addr string,port []string,createtime int64){
	user := GetUserInfo(username)
	var classpod Classpod
	db.Table("classpod").Select("*").Where("podname = ?",podname).First(&classpod)
	classpod.Userid = user.ID
	classpod.Podname = podname
	classpod.Class = class
	classpod.Image = image
	classpod.Addr = addr
	var ports string
	ports = port[0]
	if(len(port)>1){
		for i:=1;i<len(port);i++{
			ports = ports + "," + port[i]
		}
	}
	classpod.Port = ports
	classpod.Createtime = createtime
	db.Save(&classpod)
}

func IsClassPod(podname string) (bool,string){
	var classpod Classpod
	db.Table("classpod").Select("*").Where("podname = ?",podname).First(&classpod)
	return classpod.ID != 0,classpod.Class
}

func GetClassPods(username string)[]Classpod{
	user := GetUserInfo(username)
	var classpods []Classpod
	db.Table("classpod").Select("*").Where("userid = ?",user.ID).Find(&classpods)
	return classpods
}

func GetClassPod(podname string)Classpod{
	var cp Classpod
	db.Table("classpod").Select("*").Where("podname = ?",podname).First(&cp)
	return cp
}

func GetClassPodsByClassen(classen string)[]Classpod{
	var cps []Classpod
	db.Table("classpod").Select("*").Where("class = ?",classen).Find(&cps)
	return cps
}

func DeleteClassPod(podname string){
	db.Table("classpod").Where("podname = ?",podname).Delete(Classpod{})
}

func DeleteUserClasspods(userid int){
	db.Table("classpod").Where("userid = ?",userid).Delete(Classpod{})
}

func DeleteClass(class string){
	db.Table("classpod").Where("class = ?",class).Delete(Classpod{})
}