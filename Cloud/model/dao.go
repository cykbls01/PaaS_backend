package model

import "fmt"

type User struct {
	Id 				string 		`gorm:"primary_key" json:"id"`
	Name	        string		`json:"name"`
	Role            int		    `json:"role"`
	Department      int		    `json:"department_id"`
}

type Experiment struct {
	Id 				int 		`gorm:"primary_key" json:"id"`
	Name	        string		`json:"name"`
	CourseId	    int 		`json:"course_id"`
	Deadline        string      `json:"deadline"`
}

type Course struct {
	Id 				int 		`gorm:"primary_key" json:"id"`
	Name	        string		`json:"name"`
	TeacherId	    string 		`json:"teacher_id"`
}

type Assistant struct {
	Id 				int 		`gorm:"primary_key" json:"id"`
	StudentId	    string		`json:"student_id"`
	CourseId	    int 		`json:"course_id"`
}

type CourseStudentMapping struct {
	StudentId	    string		`json:"student_id"`
	CourseId	    int 		`json:"course_id"`
}

func GetUser(id string)*User{
	var user User
	user.Id = id
	db.Table("user").Select("*").Where(&user).First(&user)
	return &user
}

func ExistUser(id string)bool{
	var user User
	db.Table("user").Select("*").Where("id = ?",id).Find(&user)
	return user != User{}
}

func GetUserByDepartment(id int)[]*User{
	var users []*User
	db.Table("user").Select("*").Where("department_id = ?",id).Find(&users)
	return users
}

func GetUsers()[]*User{
	var users []*User
	db.Table("user").Select("*").Find(&users)
	return users
}

func GetExperiment(id int)*Experiment{
	var exp Experiment
	db.Table("experiment").Select("*").Where("id = ?",id).First(&exp)
	return &exp
}

func GetExperimentByCourse(id int)[]*Experiment{
	var exps []*Experiment
	db.Table("experiment").Select("*").Where("course_id = ?",id).Find(&exps)
	return exps
}

func GetCourse(id int)*Course{
	var course Course
	db.Table("course").Select("*").Where("id = ?",id).First(&course)
	return &course
}

func GetCourseByAssistant(id string)[]*Assistant{
	var assistants []*Assistant
	db.Table("assistant").Select("*").Where("student_id = ?",id).Find(&assistants)
	return assistants
}

func GetCourseByTeacher(id string)[]*Course{
	var courses []*Course
	db.Table("course").Select("*").Where("teacher_id = ?",id).Find(&courses)
	return courses
}

func GetCourseByStudent(id string)[]*Course{
	var relation []*CourseStudentMapping
	var courses []*Course
	db.Table("course_student_mapping").Select("*").Where("student_id = ?", id).Find(&relation)
	for _,v := range relation{
		courses  = append(courses,GetCourse(v.CourseId))
	}
	return courses
}

func GetCourseUserRelation(id string,courseId int)int{
	var assistant Assistant
	var relation CourseStudentMapping

    fmt.Println(id)
	fmt.Println(courseId)
	db.Table("assistant").Select("*").Where("student_id = ? AND course_id = ?", id,courseId).Find(&assistant)
	db.Table("course_student_mapping").Select("*").Where("student_id = ? AND course_id = ?", id,courseId).Find(&relation)
    course := GetCourse(courseId)
    fmt.Println(relation)


	if (assistant != Assistant{}){
		return 2
	}else if(relation != CourseStudentMapping{}){
		return 1
	}else if id == course.TeacherId{
		return 2
	}else{
		return 0
	}
}
