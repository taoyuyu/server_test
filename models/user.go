package models

import (
	"github.com/astaxie/beego/orm"
)

type UserBasic struct {
	UserBasicID       int    `orm:"pk;column(UserBasic_id)"`
	UserName          string `orm:"column(Username)"`
	Password          string `orm:"column(Password)"`
	UserEmail         string `orm:"column(Useremail)"`
	UserType          int    `orm:"column(UserType)"`
	UserRegTime       string `orm:"column(User_regtime)"`
	DepartmentID      int    `orm:"column(departmentId)"`
	DepartmentName    string `orm:"column(departmentName)"`
	PostID            int    `orm:"column(postId)"`
	RoleID            int    `orm:"column(roleId)"`
	PasswordSalt      string `orm:"column(passwordsalt)"`
	ReamName          string `orm:"column(realName)"`
	State             int    `orm:"column(state)"`
	PhoneNumber       string `orm:"column(phoneNumber)"`
	HasEmailConformed int    `orm:"column(hasEmailConfirmed)"`
	HasPhoneConformed int    `orm:"column(hasPhoneConfirmed)"`
}

func FindUserBasicByUserName(userName string) (UserBasic, error) {
	o := orm.NewOrm()
	user := new(UserBasic)
	err := o.Raw("select * from user_basic where Username = ?;", userName).QueryRow(user)
	return *user, err
}
