package util

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Log(c*gin.Context,v ...interface{}){
	id,_ := GetUserIdFromAuthInfo(c)
	log.Println(id,c.Request.URL,v)
}
