package api

import (
	"Cloud/pkg/consts"
	"Cloud/pkg/util"
	"bufio"
	"encoding/gob"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wzyonggege/logger"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

func GetNodeInfo(c *gin.Context){
	_ = checkToken(c)
	var data Nodeinfos
	if code != consts.SUCCESS{
		goto LAST
	}else{
		err = Read(&data)
		if err != nil{
			goto LAST
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}


func Read(nodes *Nodeinfos)error{
	file,err := os.Open(gob_path)
	if err != nil{
		logger.Error("打开gob数据出错,文件不存在")
		return err
	}
	dec := gob.NewDecoder(file)
	err2 := dec.Decode(nodes)
	if err2 != nil{
		logger.Error("gob数据编码出错")
		return err2
	}
	return nil
}

func Write(nodes *Nodeinfos) error {
	file, err := os.Create(gob_path)
	if err != nil {
		logger.Error("打开gob数据出错")
		return err
	}
	enc := gob.NewEncoder(file)
	err2 := enc.Encode(nodes)
	if err2 != nil{
		logger.Error("gob数据译码出错")
		return err2
	}
	return nil
}

type Sortnode []Node

func (s Sortnode) Less(i,j int) bool{
	return s[i].Node < s[j].Node
}

func (s Sortnode) Len() int{
	return len(s)
}

func (s Sortnode) Swap(i, j int) {s[i],s[j]=s[j],s[i]}

func GetTotalResource(){
	var data Sortnode
	if code != consts.SUCCESS{
		goto LAST
	}else{
		var nodes util.NodeMetricsList

		_ = util.GetNodeMetric(util.GetClientset(),&nodes)
		nodeItf := util.GetNodeItf()

		for _,item := range nodes.Items{

			var nodeinfo Node

			k8sNode,_ := nodeItf.Get(item.Metadata.Name,meta_v1.GetOptions{})
			nodeinfo.Node = item.Metadata.Name
			cmd := exec.Command("bash","/root/shells/df.sh",nodeinfo.Node)
			output,_ := cmd.Output()
			tmpstr := string(output)
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

		sort.Sort(data)
		var tmpnodes Nodeinfos
		err = Read(&tmpnodes)
		if len(tmpnodes) == 0{
			for _,item:= range data{
				var nodeinfo Nodeinfo
				var chartdata ChartData
				nodeinfo.Node = item.Node
				nodeinfo.TotCPU = item.TotCPU
				nodeinfo.Totmem = item.Totmem
				nodeinfo.TotSto = item.TotSto
				chartdata.UsedCPU = item.UsedCPU
				chartdata.Usedmem = item.Usedmem
				chartdata.UsedSto = item.UsedSto
				chartdata.Timestamp = time.Now().Unix() * 1000
				nodeinfo.ChartData = append(nodeinfo.ChartData,chartdata)

				tmpnodes = append(tmpnodes,nodeinfo)
			}
		}else{
			for i,item:= range data{
				var chartdata ChartData
				chartdata.UsedCPU = item.UsedCPU
				chartdata.Usedmem = item.Usedmem
				chartdata.UsedSto = item.UsedSto
				chartdata.Timestamp = time.Now().Unix() * 1000
				tmpnodes[i].ChartData = append(tmpnodes[i].ChartData,chartdata)
				if len(tmpnodes[i].ChartData) > 60{
					tmpnodes[i].ChartData = tmpnodes[i].ChartData[1:]
				}
			}
		}
		err = Write(&tmpnodes)
		if err != nil{
			return
		}
	}
LAST:
	return
}

type INFO struct{
	Level		string		`json:"level"`
	Time		string		`json:"time"`
	Msg			string		`json:"msg"`
	User		string		`json:"user"`
	Reason		string		`json:"reason,omitempty"`
	Type		string		`json:"type"`
}

func GetLog(c *gin.Context){
	_ = checkToken(c)
	var data []INFO
	if code != consts.SUCCESS{
		goto LAST
	}else{
		dir := consts.PATH + "logs/"
		files,_ := ioutil.ReadDir(dir)
		for _,f := range files{
			fname := f.Name()
			fi,_ := os.Open(dir+fname)
			br := bufio.NewReader(fi)
			for{
				a,_,c := br.ReadLine()
				if c == io.EOF{
					break;
				}
				var tinfo INFO
				json.Unmarshal(a,&tinfo)
				tinfo.Time = tinfo.Time[0:strings.Index(tinfo.Time,".")]
				tinfo.Time = strings.ReplaceAll(tinfo.Time,"T"," ")
				if tinfo.Level == "error"{
					tinfo.Msg = tinfo.Msg + ":" + tinfo.Reason
					tinfo.Reason = ""
				}
				data = append(data, tinfo)
			}
		}
	}
LAST:
	c.Set("code",code)
	c.Set("msg",consts.GetMsg(code))
	c.Set("data",data)
	return
}