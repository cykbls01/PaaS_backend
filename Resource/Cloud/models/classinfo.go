package models

import "time"

type Classinfo struct {
	Classid      	int    		`gorm:"primary_key" json:"classid"`
	Userid       	int    		`json:"userid"`
	Classen      	string		`json:"classen"`
	Classch      	string 		`json:"classch"`
	Classprofile 	string		`json:"classprofile"`
	Endtime			int64		`json:"endtime"`
}

func IsClassEnded(classid int)bool{
	var ci Classinfo
	db.Table("classinfo").Select("*").Where("classid = ?",classid).First(&ci)
	now := time.Now().Unix()
	return now > ci.Endtime
}

func AddClassInfo(username string,classen string,classch string,classprofile string,endtime int64){
	user := GetUserInfo(username)
	var classinfo Classinfo
	db.Table("classinfo").Select("*").Where("classen = ?",classen).First(&classinfo)
	classinfo.Classen = classen
	classinfo.Classch = classch
	classinfo.Classprofile = classprofile
	classinfo.Userid = user.ID
	classinfo.Endtime = endtime
	db.Save(&classinfo)
}

func GetClassInfo(username string)[]Classinfo{
	user := GetUserInfo(username)
	var classinfos []Classinfo
	db.Table("classinfo").Select("*").Where("userid = ?",user.ID).Find(&classinfos)
	return classinfos
}

func GetAllClasses()[]Classinfo{
	var cis []Classinfo
	db.Table("classinfo").Select("*").Order("userid").Find(&cis)
	return cis
}

func GetClassByID(id int)Classinfo{
	var classinfo Classinfo
	db.Table("classinfo").Select("*").Where("classid = ?",id).First(&classinfo)
	return classinfo
}

func DeleteClassInfo(id int){
	db.Table("classinfo").Where("classid = ? ",id).Delete(Classinfo{})
}

func IsClassExist(classen string)bool{
	var classinfo Classinfo
	db.Table("classinfo").Select("*").Where("classen = ?",classen).First(&classinfo)
	return classinfo.Classid != 0
}