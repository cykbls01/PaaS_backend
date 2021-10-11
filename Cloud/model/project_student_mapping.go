package model


type Project_Student_Mapping struct {
	ProjectId 		int         `json:"project_id"`
	StudentId	    string		`json:"student_id"`
	Type		    int		    `json:"type"`
}

func AddRelation(projectid,Type int,studentid string){
	var relation Project_Student_Mapping
	relation.ProjectId = projectid
	relation.StudentId = studentid
	relation.Type = Type
	db.Table("project_student_mapping").Save(&relation)
}

func DeleteRelation(projectid int,studentid string){
	var relation Project_Student_Mapping
	relation.ProjectId = projectid
	relation.StudentId = studentid
	db.Table("project_student_mapping").Delete(&relation)
}

func GetRelationByProjectid(projectid int)[]*Project_Student_Mapping{
	var relation []*Project_Student_Mapping
	db.Table("project_student_mapping").Select("*").Where("project_id = ?",projectid).Find(&relation)
	return relation
}

func GetRelationByStudentid(studentid string)[]*Project_Student_Mapping{
	var relation []*Project_Student_Mapping
	db.Table("project_student_mapping").Select("*").Where("student_id = ?",studentid).Find(&relation)
	return relation
}

func GetRelation(projectid int,studentid string)*Project_Student_Mapping{
	var relation Project_Student_Mapping
	relation.ProjectId = projectid
	relation.StudentId = studentid
	db.Table("project_student_mapping").Select("*").Where(&relation).First(&relation)
	return &relation
}

func GetMembers(proid int)[]*User{
	relations := GetRelationByProjectid(proid)
	var users []*User
	for _,v := range relations{
		user := GetUser(v.StudentId)
		users = append(users, user)
	}
	return users
}

