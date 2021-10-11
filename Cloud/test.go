package main

import (
   "Cloud/k8s"
   "Cloud/model"
   "Cloud/router"
   "Cloud/util"
)

func main(){

   k8s.K8sInit()
   model.MysqlInit()
   util.RedisInit()
   router.RouterInit()

}