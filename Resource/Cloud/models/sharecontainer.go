package models

type Sharecontainer struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Containername	string		`json:"containername"`
	Password		string 		`json:"password"`
}

func AddShareContainer(containername string,password string){
	var share Sharecontainer
	share.Containername = containername
	share.Password = password
	db.Save(&share)
}

func GetSharePassword(containername string)string{
	var share Sharecontainer
	db.Table("sharecontainer").Select("*").Where("containername = ?",containername).First(&share)
	return share.Password
}

func DeleteShareContainer(containername string){
	db.Table("sharecontainer").Where("containername = ?",containername).Delete(Sharecontainer{})
}
