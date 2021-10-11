package model


type Node struct{
	Id 				int 		`gorm:"primary_key" json:"id"`
	Name		    string	    `json:"name"`
	UserName		string		`json:"user_name"`
	Password		string		`json:"password"`
	Host			string		`json:"host"`
}

func GetNode(nodeName string)Node{
	var node Node
	db.Table("node").Select("*").Where("name = ?",nodeName).Find(&node)
	return node
}