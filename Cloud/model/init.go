package model
import (
	"Cloud/util"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)
var db *gorm.DB

func MysqlInit(){
	var err error

	sec,err := util.Cfg.GetSection("database")
	var dbType,dbName,user,password,host string
	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()

	db,err = gorm.Open(dbType,fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil{
		log.Println(err)
	}

	if err != nil {
		panic("连接数据库失败")
	}
	db.SingularTable(true)
}


