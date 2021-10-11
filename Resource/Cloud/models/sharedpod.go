package models

import "time"

type Sharedpod struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Userid	 		int			`json:"userid"`
	Podname			string		`json:"podname"`
	Dbname			string		`json:"dbname"`
	Password		string 		`json:"password"`
	Createtime		int64		`json:"createtime" sql:"type:bigint"`
}

func CheckExist(dbname string,podname string)bool{
	var sharedpod Sharedpod
	db.Table("sharedpod").Select("*").Where("dbname = ? and podname = ?",dbname,podname).First(&sharedpod)
	if sharedpod.ID>0{
		return true
	}
	return false
}

func AddToDatabase(podname string,username string,database string,password string){
	var pod Sharedpod
	user := GetUserInfo(username)
	pod.Podname = podname
	pod.Userid = user.ID
	pod.Dbname = database
	pod.Password = password
	pod.Createtime = time.Now().Unix()
	db.Save(&pod)
}

func GetPodIfExist(username string)[]Sharedpod{
	user := GetUserInfo(username)
	var pods []Sharedpod
	db.Table("sharedpod").Select("*").Where("userid = ?",user.ID).Find(&pods)
	return pods
}

func DeleteSharepodInfoFromDB(dbname string){
	db.Table("sharedpod").Where("dbname = ?",dbname).Delete(Sharedpod{})
}

func DeleteSharepod(podname string){
	db.Table("sharedpod").Where("podname = ?",podname).Delete(Sharedpod{})
}