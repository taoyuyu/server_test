package routers

import (
	"server_test/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/getInfo", &controllers.UserController{})
}
