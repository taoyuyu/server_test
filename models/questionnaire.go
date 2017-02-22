package models

import (
	"github.com/astaxie/beego/orm"
)

type Questionnaire struct {
	QnaireID             int    `orm:"pk;column(Qnaire_id)"`
	UserBasicID          int    `orm:"column(UserBasic_id)"`
	QnaireTitle          string `orm:"column(Qnaire_title)"`
	QnaireDescription    string `orm:"column(Qnaire_description)"`
	QnaireQuesNum        int    `orm:"column(Qnaire_quesNum)"`
	QnaireCreateTime     string `orm:"column(Qnaire_createTime)"`
	QnaireLastModifyTime string `orm:"column(Qnaire_lastModifyTime)"`
	QnaireStatus         int    `orm:"column(Qnaire_status)"`
	QnaireCategory       int    `orm:"column(Qnaire_category)"`
	QnaireLogo           string `orm:"column(Qnaire_logo)"`
	QnaireAnsSheetNum    int    `orm:"column(Qnaire_ansSheetNum)"`
	QnaireType           int    `orm:"column(Qnaire_type)"`
	QnaireIsTemplate     int    `orm:"column(Qnaire_isTemplate)"`
	CreateUserName       string `orm:"column(createUserName)"`
	QuestionnaireType    int    `orm:"column(questionnaireType)"`
	SourceID             int    `orm:"column(sourceId)"`
}

func FindQuestionnaireByUserBasicID(userBasicID int) ([]Questionnaire, error) {
	o := orm.NewOrm()
	var questionnaires []Questionnaire
	_, err := o.Raw("select * from questionnaire_basic where UserBasic_id = ?", userBasicID).QueryRows(&questionnaires)
	return questionnaires, err
}
