package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindData(c *gin.Context,value interface{})(err error){
	err = c.BindJSON(value)
    Log(c,value)
	if err != nil{
		Log(c,err)
	}
	return err
}

func ReturnSuccess(c *gin.Context,msg string,data interface{}){
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK,gin.H{
			"code" : 1001,
			"msg" :  msg,
			"data" : data,
	})
}
func ReturnError(c *gin.Context,msg string){
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK,gin.H{
		"code" : 2001,
		"msg" :  msg,
		"data" : nil,
	})
}