package controllers

import (
	"encoding/json"
	"log"
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
	data := getData(userName)
	//data := searchInDB(userName)
	c.Ctx.WriteString(data)
	return
}

func getData(userName string) string {
	//检查缓存
	result, err := redis_client.Get(userName)
	if err != nil {
		log.Println(err)
		log.Println("数据查询")
		//查询数据库
		user, _ := models.FindUserBasicByUserName(userName)
		questions, _ := models.FindQuestionnaireByUserBasicID(user.UserBasicID)

		//打包json
		bs, err1 := json.Marshal(questions)
		if err1 != nil {
			log.Println(err1)
		}
		result = string(bs)
		//写入缓存
		err2 := redis_client.Set(userName, result, 3)
		if err2 != nil {
			log.Println(err2)
		}
	}
	return result
}

func searchInDB(userName string) string {
	user, _ := models.FindUserBasicByUserName(userName)
	questions, _ := models.FindQuestionnaireByUserBasicID(user.UserBasicID)

	bs, _ := json.Marshal(questions)
	return string(bs)
}
