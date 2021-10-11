package models

type Images struct {
	ID 					int 		`gorm:"primary_key" json:"id"`
	Image				string		`json:"image"`
	Lowcpu				int			`json:"lowcpu"`
	Lowmem				int			`json:"lowmem"`
	Lowsto				int			`json:"lowsto"`
	Port				string		`json:"port"`
	Env					string		`json:"env"`
	Mntpath				string		`json:"mntpath"`
	Dupath				string		`json:"dupath"`
	Acceptfile			string		`json:"acceptfile"`
}

func CreateImage(image string,lowcpu int,lowmem int, lowsto int,port string,env string,mntpath string,dupath string,acceptfile string){
	var images Images
	db.Table("images").Select("*").Where("image = ?",image).First(&images)
	images.Image = image
	images.Lowcpu = lowcpu
	images.Lowmem = lowmem
	images.Lowsto = lowsto
	images.Port = port
	images.Env = env
	images.Mntpath = mntpath
	images.Dupath = dupath
	images.Acceptfile = acceptfile
	db.Save(&images)
}

func GetImage(image string) Images {
	var images Images
	db.Table("images").Select("*").Where("image = ?",image).First(&images)
	return images
}

func GetAllImage() []Images{
	var images []Images
	db.Table("images").Select("*").Find(&images)
	return images
}