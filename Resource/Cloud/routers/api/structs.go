package api

type auth struct {	//用户验证
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type Node struct{			//节点某一时间点信息
	Node		string		`json:"name"`
	Usedmem		int64		`json:"usedMem"`
	Totmem		int64		`json:"totalMem"`

	UsedCPU		int64		`json:"usedCPU"`
	TotCPU		int64		`json:"totalCPU"`

	UsedSto		int64		`json:"usedStorage"`
	TotSto		int64		`json:"totalStorage"`
}

type ChartData struct {		//节点监控的表单描点信息
	Usedmem		int64		`json:"y1"`
	UsedCPU		int64		`json:"y2"`
	UsedSto		int64		`json:"y3,omitempty"`
	Timestamp	int64		`json:"x"`
}

type Nodeinfo struct {		//节点监控汇总信息
	Node		string		`json:"name"`
	Totmem		int64		`json:"totalMem"`
	TotCPU		int64		`json:"totalCPU"`
	TotSto		int64		`json:"totalSto,omitempty"`
	ChartData	[]ChartData	`json:"chartData"`
}

type Nodeinfos []Nodeinfo

type form struct {					//创建容器的表单
	Name 		string 				`json:"name"`
	Podname		string				`json:"podname"`
	Linenum		int	   				`json:"linenum"`
	Image		[]string			`json:"image"`
	Env			[]Children			`json:"env"`
	Port 		[]int				`json:"port"`
	Limit		map[string]int64	`json:"limit"`
	Url			string				`json:"url"`
}

type deleteForm struct {	//删除容器的表单
	Name 		[]string 	`json:"podname"`
}

type extrainfo struct {
	Dbname		string		`json:"dbname"`
	Password	string		`json:"password"`
}

type uploadinfo struct {
	Uploadable		bool		`json:"uploadable"`
	Accept			string		`json:"accept"`
}

type owner struct {
	Username 		string		`json:"username"`
	Name			string		`json:"name"`
}

type podinfo struct{		//容器基本信息
	Type		string		`json:"type"`
	Name 		string		`json:"name"`
	Podname 	string 		`json:"podname"`
	Image 		string 		`json:"image"`
	Status 		string 		`json:"status"`
	Addr		string 		`json:"address"`
	Port 		[]string 	`json:"port"`
	Createtime	int64  		`json:"createTime"`
	Extrainfo	extrainfo	`json:"extraInfo"`
	Upload		uploadinfo	`json:"upload"`
	Owner		owner		`json:"owner,omitempty"`
}

type repository struct {	//镜像种类
	Names []string `json:"repositories"`
}

type svcversion struct{	//镜像的版本号
	Name	string 		`json:"name"`
	Tags	[]string	`json:"tags"`
}

type Children struct {	//label-value结构
	Value 	string `json:"value"`
	Label 	string `json:"label"`
	Comment string `json:"comment,omitempty"`
}

type repoinfo struct {	//镜像综合信息(种类、版本号等的综合)
	Childs	[]Children 	`json:"children"`
	Value	string 		`json:"value"`
	Label   string 		`json:"label"`
}

type resourceinfo struct{	//资源描述信息
	Requestcpu	int64		`json:"requestCPU"`
	Requestmem	int64		`json:"requestMem"`
	Requeststo	int64		`json:"requestStorage"`
}

type Useruse struct{		//用户使用信息
	Username 	string			`json:"username"`
	Name		string			`json:"name"`
	Authority	string			`json:"authority"`
	Total		resourceinfo	`json:"total"`
	Used		resourceinfo	`json:"used"`
}