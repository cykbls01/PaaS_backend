package models

type Class_auth struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Userid	 		int			`json:"userid"`
	Classid			int			`json:"classid"`
}

func CountClassMember(classid int)int{
	var cnt int
	db.Table("class_auth").Where("classid = ?",classid).Count(&cnt)
	return cnt
}

func IsUserExistInClass(userid int,classid int)bool{
	var ca Class_auth
	println(userid)
	println(classid)
	db.Table("class_auth").Select("*").Where("userid = ? and classid = ?",userid,classid).First(&ca)
	println(ca.ID)
	return ca.ID != 0
}

func AddUsertoClass(userid int,classid int){
	var ca Class_auth
	ca.Classid = classid
	ca.Userid = userid
	db.Table("class_auth").Save(&ca)
}

func GetUsersinClass(classid int)[]Class_auth{
	var cas []Class_auth
	db.Table("class_auth").Select("*").Where("classid = ?",classid).Find(&cas)
	return cas
}

func GetUsersClass(userid int)[]Class_auth{
	var cas []Class_auth
	db.Table("class_auth").Select("*").Where("userid = ?",userid).Find(&cas)
	return cas
}

func DeleteClassAuths(classid int){
	db.Table("class_auth").Where("classid = ?").Delete(Class_auth{})
}

func DeleteUserfromClass(userid int,classid int){
	db.Table("class_auth").Where("userid = ? and classid = ?",userid,classid).Delete(Class_auth{})
}