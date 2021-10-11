package models

type Classtask struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Classid			int			`json:"classid"`
	Starttime		int64		`json:"starttime"`
	Endtime			int64		`json:"endtime"`
	Taskname		string		`json:"taskname"`
	Taskgoal		string		`json:"taskgoal"`
	Taskstep		string		`json:"taskstep"`
	Filepath		string		`json:"filepath"`
}

func IsTaskExistInClass(classid int,taskname string)bool{
	var ct Classtask
	db.Table("classtask").Select("*").Where("classid = ? and taskname = ?",classid,taskname).First(&ct)
	return ct.ID != 0
}

func AddClassTask(classid int,starttime int64, endtime int64,taskname string,taskgoal string,taskstep string,
	filepath string){
	var ct Classtask
	db.Table("classtask").Select("*").Where("classid = ? and taskname = ?",classid,taskname).First(&ct)
	ct.Classid = classid
	ct.Starttime = starttime
	ct.Endtime = endtime
	ct.Taskname = taskname
	ct.Taskgoal = taskgoal
	ct.Taskstep = taskstep
	ct.Filepath = filepath
	db.Table("classtask").Save(&ct)
}

func GetTasksByClassid(classid int)[]Classtask{
	var cts []Classtask
	db.Table("classtask").Select("*").Where("classid = ?",classid).Find(&cts)
	return cts
}

func GetTaskByid(taskid int)Classtask{
	var ct Classtask
	db.Table("classtask").Select("*").Where("id = ?",taskid).First(&ct)
	return ct
}

func Get

func DeleteTaskByid(taskid int){
	db.Table("classtask").Where("id = ?",taskid).Delete(Classtask{})
}

func DeleteTasksInClass(classid int){
	db.Table("classtask").Where("classid = ?",classid).Delete(Classtask{})
}
