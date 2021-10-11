package model

import (
	"Cloud/util"
	"fmt"
)

type Apply_Record struct {
	Id 				int 		`gorm:"primary_key" json:"id"`
	SenderId	    string		`json:"sender_id"`
	ApplyTime		string		`json:"apply_time"`
	FinishTime		string		`json:"finish_time"`
	Status		    int		    `json:"status"`
	VmName			string		`json:"vm_name"`
	OperateType		int		    `json:"operate_type"`
	Detail		    string		`json:"detail"`
	ApplyMsg		string		`json:"apply_msg"`
	ReplyMsg		string		`json:"reply_msg"`
	DueTime         string      `json:"due_time"`
}
func GetAllApply(){
	var ar []Apply_Record
	db.Table("apply_record").Select("*").Find(&ar)
	fmt.Println(ar[0].Id)
}

func GetApplyByName(id,name string)*Apply_Record{
	var apply Apply_Record
	db.Table("apply_record").Select("*").Where("sender_id = ? And vm_name = ?",id,name).First(&apply)
	return &apply
}

func GetApply(Id int)*Apply_Record{
	var apply Apply_Record
	db.Table("apply_record").Select("*").Where("id = ?",Id).First(&apply)
	return &apply
}

func AddApply(userId,VmName,Detail,ApplyMsg,EndTime string){
	var apply Apply_Record
	apply.SenderId = userId
	apply.Detail = Detail
	apply.VmName = VmName
	apply.ApplyMsg = ApplyMsg
	apply.DueTime = EndTime
	apply.ApplyTime = util.GetNowTime()
	apply.Status = 0
	apply.OperateType = 4
	db.Table("apply_record").Save(&apply)
}

func UpdateApply(apply *Apply_Record){
	db.Table("apply_record").Model(&apply).Update(apply)
}

func DeleteApply(){

}