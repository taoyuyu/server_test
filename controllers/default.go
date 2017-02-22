package controllers

import (
	"fmt"
	"server_test/models"
	"server_test/redis_client"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type UserController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *UserController) Get() {
	userName := c.GetString("username")
	//检查缓存
	result, _ := redis_client.Get(userName)
	if result != `` {
		c.Data["response"] = result
		c.ServeJSON()
		return
	}
	//查询数据库
	user, _ := models.FindUserBasicByUserName(userName)
	questions, _ := models.FindQuestionnaireByUserBasicID(user.UserBasicID)
	fmt.Println(questions)
	//写入缓存

	c.Data["response"] = questions
	c.ServeJSON()
	return
}
