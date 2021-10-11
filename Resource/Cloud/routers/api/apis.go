package api

import (
	"Cloud/pkg/consts"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/apps/v1beta1"
	v12 "k8s.io/api/core/v1"
	"strings"
)

const (
	gob_path = consts.PATH + "gob/nodes.info"
)

var(
	deployment = v1beta1.Deployment{}

	service = v12.Service{}

	YAML []byte
	JSON []byte

	err error
	code int
)

//func init(){
//	logger = util.GetLogger()
//}

func GetContainerName(podname string)string {
	return podname[0:strings.LastIndex(podname[0:strings.LastIndex(podname,"-")],"-")]
}

func TestDownload(c *gin.Context){
	code = consts.SUCCESS
	var info struct{
		Filepath	string	`json:"filepath"`
	}
	c.Bind(&info)
	filename := info.Filepath[strings.LastIndex(info.Filepath,"/")+1:]
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type","application/octet-stream")
	c.File(info.Filepath)
}

func TestDownload2(c *gin.Context){
	code = consts.SUCCESS
	filepath := "/root/beihangLogin"
	filename := filepath[strings.LastIndex(filepath,"/")+1:]
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type","application/octet-stream")
	c.File(filepath)
}

func Download(c *gin.Context){
	_ = checkToken(c)
	if code != consts.SUCCESS {
		goto LAST
	}else{
		filepath := c.Param("filepath")
		filename := filepath[strings.LastIndex(filepath,"/")+1:]
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Writer.Header().Add("Content-Type","application/octet-stream")
		c.File(filepath)
	}
LAST:
	return
}

func Download2(c *gin.Context){
	filepath := c.Param("filepath")
	println(filepath)
	filename := filepath[strings.LastIndex(filepath,"/")+1:]
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type","application/octet-stream")
	c.File(filepath)
}