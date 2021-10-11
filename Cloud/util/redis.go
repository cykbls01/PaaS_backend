package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"strings"
)
var RedisConn1 redis.Conn

var pool *redis.Pool

func RedisInit(){
	sec,_ := Cfg.GetSection("redis")
	var host string
	host = sec.Key("RedisHost").String()

	pool = &redis.Pool{
		MaxIdle:     30,
		MaxActive:   1024,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if _, err = c.Do("PING"); err != nil {
				c.Close()
				fmt.Println(err)
				return nil, err
			}
			return c, err
		},
	}
	fmt.Println("Success!")
}

func GetUserIdFromAuthInfo(c *gin.Context) (id string,isLogin bool){
	token := c.Request.Header.Get("Authorization")
	RedisConn :=pool.Get()
	defer RedisConn.Close()
	id,_= redis.String(RedisConn.Do("GET",token))
	if len(id) > 0{
		isLogin = true
	}else{
		isLogin = false
		ReturnError(c,"没有登录")
	}
	id = strings.Trim(id, "\"")
	//id = "17373273"
	//isLogin = true
	return
}

func TestRedis(){
	conn := pool.Get()
	_,e := redis.String(conn.Do("GET","123"))
	fmt.Println(e)
}