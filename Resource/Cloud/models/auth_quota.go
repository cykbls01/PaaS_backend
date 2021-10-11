package models

type Auth_quota struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Userid	 		int			`json:"userid"`
	Totcpu			int64		`json:"totcpu"`
	Totmem			int64		`json:"totmem"`
	Totstorage		int64		`json:"totstorage"`
	Usedcpu			int64		`json:"usedcpu"`
	Usedmem			int64		`json:"usedmem"`
	Usedstorage		int64		`json:"usedstorage"`
}

func GetUserQuota(userid int)Auth_quota{
	var aq Auth_quota
	db.Table("auth_quota").Select("*").Where("userid = ?",userid).First(&aq)
	return aq
}

func UpdateUserQuota(aq Auth_quota){
	var tmpaq Auth_quota
	db.Table("auth_quota").Select("*").Where("userid = ?",aq.Userid).First(&tmpaq)
	tmpaq.Totstorage = aq.Totstorage
	tmpaq.Totmem = aq.Totmem
	tmpaq.Totcpu = aq.Totcpu
	tmpaq.Usedstorage = aq.Usedstorage
	tmpaq.Usedmem = aq.Usedmem
	tmpaq.Usedcpu = aq.Usedcpu
	db.Table("auth_quota").Save(&tmpaq)
}

func CreateUserQuota(userid int){
	var aq Auth_quota
	db.Table("auth_quota").Select("*").Where("userid = ?",aq.Userid).First(&aq)
	aq.Userid = userid
	aq.Totcpu = 2000
	aq.Totmem = 4000
	aq.Totstorage = 100000
	db.Table("auth_quota").Save(&aq)
}

func DeleteUserQuota(userid int){
	db.Table("auth_quota").Where("userid = ?",userid).Delete(Auth_quota{})
}