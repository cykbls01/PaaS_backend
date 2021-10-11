package model



type Project struct {
	Id 				int 		`gorm:"primary_key" json:"id"`
	SenderId	    string		`json:"sender_id"`
	Name		    string	    `json:"name"`
	ExpId		    int		    `json:"exp_id"`
	ExpName		    string		`json:"exp_name"`
	Status			int		    `json:"status"`
	CourseId		int		    `json:"course_id"`
	CourseName		string		`json:"course_name"`
}

func AddProject(userid,name,expName,courseName string,expId,courseId,status int){
	var project Project
	project.SenderId = userid
	project.Name = name
	project.CourseId = courseId
	project.ExpId = expId
	project.CourseName = courseName
	project.ExpName = expName
	project.Status = status
	db.Table("project").Save(&project)
}

func DeleteProject(projectid int){
	var project Project
	project.Id = projectid
	db.Table("project").Delete(&project)
}

func UpdateProject(project Project,name string){
	db.Table("project").Model(&project).Update("name",name)
}

func GetProject(projectid int)*Project{
	var project Project
	project.Id = projectid
	db.Table("project").Select("*").Where(&project).First(&project)
	return &project
}

func GetProjectByStudentId(studentid string)[]*Project{
	var project []*Project
	db.Table("project").Select("*").Where("sender_id = ?",studentid).Find(&project)
	return project
}

func GetProjectByCourseId(courseid int)[]*Project{
	var project []*Project
	db.Table("project").Select("*").Where("course_id = ?",courseid).Find(&project)
	return project
}

func GetProjectByAssistantId(id string)[]*Project{
	var project []*Project
	relations := GetCourseByAssistant(id)
	for _,v := range relations{
		project = append(project,GetProjectByCourseId(v.CourseId)...)
	}
	return project
}

func GetProjectByTeacherId(id string)[]*Project{
	var project []*Project
	courses := GetCourseByTeacher(id)
	for _,v := range courses{
		project = append(project,GetProjectByCourseId(v.Id)...)
	}
	return project
}
