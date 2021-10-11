package consts

const (
	MYSQL_PORT = 3306
	TOMCAT_PORT = 8080
	NGINX_PORT = 80
)

var (
	port map[string]int
)

func init(){
	port = make(map[string]int)
	port["mysql"] = MYSQL_PORT
	port["tomcat"] = TOMCAT_PORT
	port["nginx"] = NGINX_PORT
}

func Getport(svc string)int{
	return port[svc]
}

