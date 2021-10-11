package api

import (
	"Cloud/k8s"
	"Cloud/model"
	"Cloud/util"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"regexp"
	"strings"
)

func NodeMonitor(c *gin.Context){
	// 登录参数验证
	//var id string
	//var isLogin bool
	//if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
	//	return
	//}

	var nodes k8s.NodeMetricsList
	var data []k8s.NodeInfo

	_ = k8s.GetNodeMetric(&nodes)
	nodeItf := k8s.GetNodeItf()

	for _,item := range nodes.Items{

		var nodeinfo k8s.NodeInfo

		k8sNode,_ := nodeItf.Get(context.TODO(),item.Metadata.Name,meta_v1.GetOptions{})
		nodeinfo.Node = item.Metadata.Name
		node := model.GetNode(nodeinfo.Node)

		util.CreateClient(node.Host,22,node.UserName,node.Password)
		tmpstr := util.RunShell("df -h /home/buaa | grep /home/buaa")
		reg := regexp.MustCompile(`\s+`)
		strs := reg.Split(tmpstr,-1)
		sto := resource.MustParse(strs[2])
		mem := resource.MustParse(item.Usage.Memory)
		cpu := resource.MustParse(item.Usage.CPU)
		nodeinfo.Usedmem = mem.ScaledValue(resource.Mega)
		nodeinfo.UsedCPU = cpu.MilliValue()
		nodeinfo.UsedSto = sto.ScaledValue(resource.Mega)
		totSto := resource.MustParse(strs[1])

		nodeinfo.Totmem = k8sNode.Status.Allocatable.Memory().ScaledValue(resource.Mega)
		nodeinfo.TotCPU = k8sNode.Status.Allocatable.Cpu().MilliValue()
		nodeinfo.TotSto = totSto.ScaledValue(resource.Mega)
		data = append(data,nodeinfo)
	}
	util.ReturnSuccess(c,"获取成功",data)
}

func PodLog(c *gin.Context){
	//登录参数验证
	var id string
	var isLogin bool
	if id,isLogin = util.GetUserIdFromAuthInfo(c); !isLogin{
		return
	}
	id = id


	// 解析参数
	type Json struct {
		PodName string `json:"podName"`
		Namespace string `json:"namespace"`
	}
	var JsonData Json
	err := c.BindJSON(&JsonData)
	log.Println(err)

	podItf := k8s.GetPodItf(JsonData.Namespace,true)
	//name := JsonData.PodName
	req := podItf.GetLogs(JsonData.PodName,&v12.PodLogOptions{})
	podLogs, _ := req.Stream(context.TODO())

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	data := buf.String()
	str := strings.Split(data,"\n")


	util.ReturnSuccess(c,"成功获取日志",str)
	return
}
