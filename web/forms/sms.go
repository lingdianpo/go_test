package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"` //手机号码格式
	Type   string `json:"type" form:"type" binding:"required,oneOf=1 2"`
}
