package router

import (
	"Cloud/router/api"
	"Cloud/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RouterInit(){

	r := gin.Default()
	r.POST("/test",test1)
	r.OPTIONS("/test",Options)

	r.POST("/api/applyContainer",api.ApplyContainer)
	r.OPTIONS("/api/applyContainer",Options)
	r.POST("/api/createContainer",api.CreateContainer)
	r.OPTIONS("/api/createContainer",Options)
	r.POST("/api/deleteContainer",api.DeleteContainer)
	r.OPTIONS("/api/deleteContainer",Options)
	r.POST("/api/listAContainer",api.ListContainerByAdmin)
	r.OPTIONS("/api/listAContainer",Options)
	r.POST("/api/listUContainer",api.ListContainerByUser)
	r.OPTIONS("/api/listUContainer",Options)
	r.POST("/api/updateContainer",api.UpdateContainer)
	r.OPTIONS("/api/updateContainer",Options)

	r.POST("/api/addProject",api.AddProject)
	r.OPTIONS("/api/addProject",Options)
	r.POST("/api/updateProject",api.UpdateProject)
	r.OPTIONS("/api/updateProject",Options)
	r.POST("/api/podLog",api.PodLog)
	r.OPTIONS("/api/podLog",Options)
	r.POST("/api/deleteProject",api.DeleteProject)
	r.OPTIONS("/api/deleteProject",Options)
	r.POST("/api/addMember",api.AddMember)
	r.OPTIONS("/api/addMember",Options)
	r.POST("/api/deleteMember",api.DeleteMember)
	r.OPTIONS("/api/deleteMember",Options)
	r.POST("/api/listSProject",api.ListProjectByStudent)
	r.OPTIONS("/api/listSProject",Options)
	r.POST("/api/listTProject",api.ListProjectByTeacher)
	r.OPTIONS("/api/listTProject",Options)
	r.POST("/api/listMember",api.ListMember)
	r.OPTIONS("/api/listMember",Options)
	r.POST("/api/getExp",api.GetExps)
	r.OPTIONS("/api/getExp",Options)
	r.POST("/api/getCourse",api.GetCourses)
	r.OPTIONS("/api/getCourse",Options)

	r.POST("/api/nodeMetric",api.NodeMonitor)
	r.OPTIONS("/api/nodeMetric",Options)
	//r.POST("/api/podMetric",api.PodMonitor)
	//r.OPTIONS("/api/podMetric",Options)



	r.GET("/ping",func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		fmt.Println(c.Request.Body)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}


func Options(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization,Content-Type")
	w.Header().Set("Access-Control-Allow-Methods","GET,POST,DELETE")
}


func test1(c *gin.Context){
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	util.ReturnSuccess(c,"ok",nil)
}

