package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"server_test/models"
	"server_test/redis_client"
	"time"

	"github.com/astaxie/beego"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

type MainController struct {
	beego.Controller
}

type UserController struct {
	beego.Controller
}

type SampController struct {
	beego.Controller
}

func (c *SampController) Get() {
	c.TplName = "sample.tpl"
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *UserController) Get() {
	userName := c.GetString("username")
	if userName == `` {
		c.Ctx.WriteString("参数错误")
		return
	}
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

		//查询数据库
		user, err := models.FindUserBasicByUserName(userName)
		if err != nil {
			return `username not existed`
		}
		fmt.Println("getinfo: ", err)
		questions, _ := models.FindQuestionnaireByUserBasicID(user.UserBasicID)
		size := len(questions)
		qids := make([]int, size)
		for key, value := range questions {
			qids[key] = value.QnaireID
		}

		//打包json
		bs, err1 := json.Marshal(qids)
		if err1 != nil {
			log.Println(err1)
		}
		result = string(bs)
		//写入缓存
		err2 := redis_client.Set(userName, result, 5)
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
