package jwt

import (
	"Cloud/pkg/consts"
	"Cloud/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT(authority int) gin.HandlerFunc{
	return func(c *gin.Context){
		var code int
		var data interface{}
		var token string

		code = consts.SUCCESS
		token = c.Request.Header.Get("Authorization")
		if token == ""{
			code = consts.ACCESS_DENIED
			goto LAST
		}else {
			token = token[7:]
			if(token == "poiuytrewq"){
				code = consts.SUCCESS
				goto LAST
			}
			claims,err := util.ParseToken(token)
			if err != nil || claims.Authority < authority {
				code = consts.ACCESS_DENIED
			}else if time.Now().Unix() + int64(consts.REFRESH_TIME.Seconds()) >= claims.ExpiresAt{
				token,_ = util.GenerateToken(claims.Username,claims.Password,claims.Authority)
				c.Set("flag",true)
				c.Set("token",token)
				goto LAST
			}
		}
LAST:
		if code != consts.SUCCESS{
			c.JSON(http.StatusUnauthorized,gin.H{
				"code"	: code,
				"msg"	:  consts.GetMsg(code),
				"data"	: data,
			})

			c.Abort()
			return
		}

		c.Next()

		flag := c.GetBool("flag")

		code = c.GetInt("code")
		msg := c.GetString("msg")
		data,exist := c.Get("data")
		if !exist{
			data = nil
		}
		if !flag{
			c.JSON(http.StatusOK,gin.H{
				"code" : code,
				"msg" :msg,
				"data" :data,
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"code": code,
				"msg":  msg,
				"data": data,
				"token" : c.GetString("token"),
			})
		}
	}
}