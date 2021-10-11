package routers


import (
	"Cloud/middleware/jwt"
	"Cloud/pkg/consts"
	"Cloud/pkg/setting"
	"Cloud/routers/api"
	"github.com/gin-gonic/gin"
)

func InitROuter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.POST("/api/login", api.GetAuth)

	r.GET("/api/ssh", api.Ssh)

	r.POST("/createuser", api.CreateUser)

	r.POST("/api/editUser",jwt.JWT(consts.AUTH_ADMIN),api.EditUserRS)

	r.GET("/api/getQuota/:username",jwt.JWT(consts.AUTH_ADMIN),api.GetUserRS)

	r.GET("/api/getQuota",jwt.JWT(consts.AUTH_STUDENT),api.GetUserRS)

	r.GET("/api/userList",jwt.JWT(consts.AUTH_ADMIN),api.GetUserUsage)

	r.GET("/api/podlist/:username",jwt.JWT(consts.AUTH_ADMIN), api.ListPods)

	r.GET("/api/nodeInfo",jwt.JWT(consts.AUTH_ADMIN),api.GetNodeInfo)

	r.POST("/api/delete/:username",jwt.JWT(consts.AUTH_ADMIN), api.Delete)

	r.POST("/api/podlog/:username",jwt.JWT(consts.AUTH_ADMIN), api.GetPodLog)

	//r.Use(jwt.JWT(consts.AUTH_TEACHER))
	//{
	//
	//}

	r.POST("/api/podlog",jwt.JWT(consts.AUTH_STUDENT), api.GetPodLog)

	r.GET("/api/podlist",jwt.JWT(consts.AUTH_STUDENT), api.ListPods)

	r.POST("/api/create",jwt.JWT(consts.AUTH_STUDENT), api.Create)

	r.POST("/api/delete",jwt.JWT(consts.AUTH_STUDENT), api.Delete)

	r.GET("/api/images",jwt.JWT(consts.AUTH_STUDENT), api.GetImage)

	r.POST("/api/validateName",jwt.JWT(consts.AUTH_STUDENT), api.IsNameExist)

	r.GET("/api/currentUser",jwt.JWT(consts.AUTH_STUDENT), api.GetUserInfo)

	r.POST("/api/validateOldPassword",jwt.JWT(consts.AUTH_STUDENT), api.CheckPassword)

	r.POST("/api/updateUserInfo",jwt.JWT(consts.AUTH_STUDENT), api.UpdateUserInfo)

	r.POST("/api/upload",jwt.JWT(consts.AUTH_STUDENT),api.UploadFile)

	r.POST("/api/upload/:username",jwt.JWT(consts.AUTH_STUDENT),api.UploadFile)

	r.OPTIONS("/api/upload",jwt.JWT(consts.AUTH_STUDENT),api.UploadFile)

	r.POST("/api/resetPassword",jwt.JWT(consts.AUTH_ADMIN),api.ResetUser)

	r.POST("/api/batchcreateuser",jwt.JWT(consts.AUTH_ADMIN),api.BatchCreateUser)

	r.POST("/api/deleteuser",jwt.JWT(consts.AUTH_ADMIN),api.DeleteUser)

	r.GET("/api/getLogs",jwt.JWT(consts.AUTH_STUDENT),api.GetLog)



	r.POST("/api/createsharepod",jwt.JWT(consts.AUTH_TEACHER),api.JoinSharepod)

	r.POST("/api/deletesharepod",jwt.JWT(consts.AUTH_TEACHER),api.RemoveSharepod)

	r.POST("/api/deletesharepod/:username",jwt.JWT(consts.AUTH_TEACHER),api.RemoveSharepod)



	r.POST("/api/addimage",jwt.JWT(consts.AUTH_ADMIN),api.AddNewImage)

	r.POST("/api/getimageinfo",jwt.JWT(consts.AUTH_STUDENT),api.GetImageInfo)

	r.GET("/api/getallimageinfo",jwt.JWT(consts.AUTH_ADMIN),api.GetAllImage)


	r.POST("/api/updateimageinfo",jwt.JWT(consts.AUTH_ADMIN),api.UpdateImage)

//
	r.POST("/api/setuserauthority",jwt.JWT(consts.AUTH_ADMIN),api.SetUserAuth)


	//r.GET("/api/shareimages",jwt.JWT(consts.AUTH_ADMIN),api.GetShareRepos)

	//r.POST("/api/createsharecontainer",jwt.JWT(consts.AUTH_ADMIN),api.CreateSharePod)

	r.GET("/api/listsharepod",jwt.JWT(consts.AUTH_STUDENT),api.ListSharePod)

	r.POST("/api/deletesharecontainer",jwt.JWT(consts.AUTH_ADMIN),api.DeleteSharePod)


	r.POST("/api/getpodusage",jwt.JWT(consts.AUTH_STUDENT),api.GetPodUsage)



	r.POST("/api/addclassinfo",jwt.JWT(consts.AUTH_TEACHER),api.CreateClassInfo)
	r.POST("/api/updateclassinfo",jwt.JWT(consts.AUTH_TEACHER),api.UpdateClassInfo)
	r.GET("/api/getclassinfo",jwt.JWT(consts.AUTH_TEACHER),api.GetClassInfo)
	r.POST("/api/batchcreate",jwt.JWT(consts.AUTH_TEACHER),api.BatchCreateViaExistedPod)
	r.POST("/api/getclasspods",jwt.JWT(consts.AUTH_TEACHER),api.GetClassPods)
	r.POST("/api/deleteclasspods",jwt.JWT(consts.AUTH_TEACHER),api.DeleteClassPods)
	r.POST("/api/deleteclass",jwt.JWT(consts.AUTH_TEACHER),api.DeleteClass)
	r.POST("/api/isclassexist",jwt.JWT(consts.AUTH_TEACHER),api.IsClassExist)
	r.GET("/api/listallpods",jwt.JWT(consts.AUTH_ADMIN),api.ListAllPods)

	r.POST("/api/addusertoclass",jwt.JWT(consts.AUTH_TEACHER),api.AddUsertoClass)
	r.POST("/api/getuserinclass",jwt.JWT(consts.AUTH_TEACHER),api.GetClassUserList)
	r.POST("/api/deleteusersfromclass",jwt.JWT(consts.AUTH_TEACHER),api.DeleteUserFromClass)

	//task apis for teachers
	r.POST("/api/addtask",jwt.JWT(consts.AUTH_TEACHER),api.AddTask)
	r.POST("/api/updatetask",jwt.JWT(consts.AUTH_TEACHER),api.UpdateTask)
	r.POST("/api/gettasksinclass",jwt.JWT(consts.AUTH_TEACHER),api.GetTasksInClass)
	r.POST("/api/deletetask",jwt.JWT(consts.AUTH_TEACHER),api.DeleteTask)   //
	r.POST("/api/getsubmitsintask",jwt.JWT(consts.AUTH_TEACHER),api.GetSubmitsInTask)

	//task apis for students
	r.POST("/api/submittask",jwt.JWT(consts.AUTH_STUDENT),api.SubmitTask)
	r.GET("/api/getexistingtasks",jwt.JWT(consts.AUTH_STUDENT),api.GetExistingTasks)

	//task apis for all
	r.POST("/api/gettaskdetail",jwt.JWT(consts.AUTH_STUDENT),api.GetTaskDetails)
	r.GET("/api/download/*filepath",jwt.JWT(consts.AUTH_STUDENT),api.Download)

	r.POST("/api/updatealluserquota",jwt.JWT(consts.AUTH_ADMIN),api.UpdateAllUserQuota)
	r.POST("/api/syncpodtodb",jwt.JWT(consts.AUTH_ADMIN),api.SyncPodsToDB)
	r.POST("/api/syncclasspodtodb",jwt.JWT(consts.AUTH_ADMIN),api.SyncClassPodsToDB)

	return r
}
