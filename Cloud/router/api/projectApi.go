package api

import (
	"Cloud/model"
	"Cloud/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)


func GetExps(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

    // 封装参数
	var exps []*model.Experiment
	courses := model.GetCourseByStudent(id)
	for _,v := range courses{
		exps = append(exps,model.GetExperimentByCourse(v.Id)...)
	}

	// 删去结束的实验
	var returnExps []*model.Experiment
	t := util.StringToTime(util.GetNowTime())
	for _,v := range exps{
		a := util.StringToTime(v.Deadline)
		if !util.CompareTime(a,t){
			returnExps = append(returnExps,v)
		}
	}

	util.ReturnSuccess(c,"成功获取实验列表",returnExps)
	return
}

func GetCourses(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

    // 封装参数
    var courses []*model.Course
	courses = model.GetCourseByTeacher(id)
	assistants := model.GetCourseByAssistant(id)
	for _,v := range assistants{
		courses = append(courses,model.GetCourse(v.CourseId))
	}
	util.ReturnSuccess(c,"成功获取课程列表",courses)
	return
}

func AddMember(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		StuId []string `json:"stuId"`
		ProId string `json:"proId"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	proid,_ := strconv.Atoi(JsonData.ProId)
	project := model.GetProject(proid)
	if project.SenderId != id{
		util.ReturnError(c,"没有权限")
		return
	}

	// 添加成员
	for _,v := range JsonData.StuId{
		if model.ExistUser(v) {
			model.AddRelation(proid, 2, v)
		}
	}
	util.ReturnSuccess(c,"成功添加成员",model.GetMembers(proid))
	return
}

func DeleteMember(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		StuId string `json:"stuId"`
		ProId string `json:"proId"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	proid,_ := strconv.Atoi(JsonData.ProId)
	project := model.GetProject(proid)
	if project.SenderId != id{
		util.ReturnError(c,"没有权限")
		return
	}

	// 删除成员
	model.DeleteRelation(proid,JsonData.StuId)
	util.ReturnSuccess(c,"成功删除成员",model.GetMembers(proid))
	return
}


func ListProjectByStudent(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 组长项目/成员项目/助教项目
	var project1 []*model.Project
	var project2 []*model.Project
	var project3 []*model.Project
	project1 = model.GetProjectByStudentId(id)
	relations := model.GetRelationByStudentid(id)
	for _,v := range relations{
		project := model.GetProject(v.ProjectId)
		project2 = append(project2, project)
	}
	project3 = model.GetProjectByAssistantId(id)

	// 封装参数
	type Json struct {
		Project1 []*model.Project `json:"project1"`
		Project2 []*model.Project `json:"project2"`
		Project3 []*model.Project `json:"project3"`
	}
	var returnData Json
	returnData.Project1 = project1
	returnData.Project2 = project2
	returnData.Project3 = project3

	util.ReturnSuccess(c,"成功获取项目列表",returnData)
	return
}

func ListProjectByTeacher(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		CourseId int `json:"courseid"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	user := model.GetUser(id)
	if user.Role != 2{
		util.ReturnError(c,"没有权限")
		return
	}

	// 获取课程项目
	var project []*model.Project

	if JsonData.CourseId == -1{
		project = model.GetProjectByTeacherId(id)
		project = append(project,model.GetProjectByAssistantId(id)...)
	}else{
		project = model.GetProjectByCourseId(JsonData.CourseId)
	}
	
	util.ReturnSuccess(c,"成功获取项目列表",project)
	return
}

func ListMember(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		ProId string `json:"proId"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	proid,_ := strconv.Atoi(JsonData.ProId)
	project := model.GetProject(proid)
	if model.GetCourseUserRelation(id,project.CourseId) ==0 && model.GetRelation(proid,id).Type ==0{
		util.ReturnError(c,"没有权限")
		return
	}

	util.ReturnSuccess(c,"成功获取成员列表",model.GetMembers(proid))
	return
}

func AddProject(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		ExpId int    `json:"expId"`
		Name  string `json:"name"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	exp := model.GetExperiment(JsonData.ExpId)
	fmt.Println(exp)
	course := model.GetCourse(exp.CourseId)
	if model.GetCourseUserRelation(id,course.Id) != 1{
		util.ReturnError(c,"没有权限")
		return
	}

	model.AddProject(id,JsonData.Name,exp.Name,course.Name,exp.Id,course.Id,1)
	util.ReturnSuccess(c,"成功添加项目",model.GetProjectByStudentId(id))
	return
}

func UpdateProject(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		ProId string `json:"proId"`
		Name  string `json:"name"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	proid,_ := strconv.Atoi(JsonData.ProId)
	project := model.GetProject(proid)
	if project.SenderId != id{
		util.ReturnError(c,"没有权限")
		return
	}

	model.UpdateProject(*project,JsonData.Name)
	util.ReturnSuccess(c,"成功删除项目",model.GetProjectByStudentId(id))
	return
}

func DeleteProject(c *gin.Context){
	// 登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}

	// 解析参数
	type Json struct {
		ProId string `json:"proId"`
	}
	var JsonData Json
	err := util.BindData(c,&JsonData)
	if err != nil{
		util.ReturnError(c,"参数错误")
	}

	// 验证权限
	proid,_ := strconv.Atoi(JsonData.ProId)
	project := model.GetProject(proid)
	if project.SenderId != id{
		log.Println("123")
		util.ReturnError(c,"没有权限")
		return
	}

	model.DeleteProject(proid)
	util.ReturnSuccess(c,"成功删除项目",model.GetProjectByStudentId(id))
	return
}