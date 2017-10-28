package main

import (
	_ "server_test/models"
	_ "server_test/redis_client"
	_ "server_test/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
	beego.beego.BConfig.WebConfig.StaticDir["/static"] = "static"
}
