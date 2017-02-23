package controllers

import (
	"encoding/json"
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
	result, err := redis_client.Get(userName)
	if err != nil {
		//查询数据库
		user, _ := models.FindUserBasicByUserName(userName)
		questions, _ := models.FindQuestionnaireByUserBasicID(user.UserBasicID)
		fmt.Println(questions)

		//打包json
		bs, err1 := json.Marshal(questions)
		if err1 != nil {
			fmt.Println(err1)
		}
		result = string(bs)
		//写入缓存
		err2 := redis_client.Set(userName, result, 5)
		if err2 != nil {
			fmt.Println(err2)
		}
	}
	c.Data["response"] = result
	c.ServeJSON()
	return
}
