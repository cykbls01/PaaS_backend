package models

import "time"

type Classtask_submit struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Userid	 		int			`json:"userid"`
	Taskid			int			`json:"taskid"`
	Filepath		string		`json:"filepath"`
	Submittime		int64		`json:"submittime"`
}

func IsUserSubmited(taskid ,userid int)bool{
	cs := GetTaskSubmitByUserid(userid,taskid)
	return cs.ID != 0
}

func TaskSubmit(userid int,taskid int,filepath string){
	var cs Classtask_submit
	db.Table("classtask_submit").Select("*").Where("userid = ? and taskid = ?",userid,taskid).First(&cs)
	cs.Userid = userid
	cs.Taskid = taskid
	cs.Filepath = filepath
	cs.Submittime = time.Now().Unix()
	db.Table("classtask_submit").Save(&cs)
}

func CountTaskSubmits(taskid int)int{
	var cnt int
	db.Table("classtask_submit").Where("taskid = ?",taskid).Count(&cnt)
	return cnt
}

func GetTaskSubmitByUserid(userid int,taskid int)Classtask_submit{
	var cs Classtask_submit
	db.Table("classtask_submit").Select("*").Where("userid = ? and taskid = ?",userid,taskid).Find(&cs)
	return cs
}

func DeleteCSByTaskid(taskid int){
	db.Table("classtask_submit").Where("taskid = ?",taskid).Delete(Classtask_submit{})
}
