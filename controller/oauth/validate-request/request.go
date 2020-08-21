package validate_request

import (
	"github.com/go-playground/validator/v10"
	"github.com/blog_backend/exception"
	"log"
	"net/url"
	"regexp"
	"strings"
)

//采用gin的自定义验证器，动态的验证

//验证请求参数
type PreauthCodeRequest struct {
	ClientAppId   string `form:"client_app_id" binding:"required"`
	RedirectUrl   string `form:"redirect_url" binding:"redirect_url"`
	UserLoginType int    `form:"user_login_type"`
	Email         string `form:"email" binding:"my_email"`
	Phone         string `form:"phone" binding:"my_phone"`
	Password      string `form:"password"`
}

func (pcr *PreauthCodeRequest) GetError(err validator.ValidationErrors) exception.MyException {
	err_msg := "参数验证失败"
	for _, val := range err {
		if val.Field() == "ClientAppId" {
			switch val.Tag() {
			case "required":
				err_msg = "client_app_id 必填"
				break
			}
		}

		if val.Field() == "RedirectUrl" {
			switch val.Tag() {
			case "redirect_url":
				err_msg = "redirect_url地址格式非法"
				break
			}
		}

		if val.Field() == "Email" {
			switch val.Tag() {
			case "my_email":
				err_msg = "email 格式不对"
				break
			}
		}

		if val.Field() == "Phone" {
			switch val.Tag() {
			case "my_phone":
				err_msg = "手机号格式不对"
				break
			}
		}
	}

	return exception.NewException(exception.VALIDATE_ERR, err_msg)
}

//预授权码，换取token
type PreauthCodeChangeTokenRequest struct {
	ClientAppId     string `form:"client_app_id" binding:"required"`
	ClientAppSecret string `form:"client_app_secret" binding:"required"`
	PreauthCode     string `form:"preauth_code" binding:"required"`
}

//通过refreshToken刷新Token
type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" binding:"required"`
}

//验证token请求
type ValidateTokenRequest struct {
	Token string `form:"token" binding:"required"`
}

//自定义预授权码的验证
func MyPhone(fl validator.FieldLevel) bool {
	log.Println("验证手机号")
	if fl.GetTag() == "phone" {
		regular := `^(1[3|4|5|8][0-9]\d{8})$`
		reg := regexp.MustCompile(regular)
		phone := fl.Field().Interface().(string)
		if strings.Compare(strings.Trim(phone, " "), "") == 0 {
			return true
		}
		return reg.MatchString(phone)
	}
	return true
}

func RedirectUrl(fl validator.FieldLevel) bool {
	log.Println("验证重定向地址")
	if fl.GetTag() == "redirect_url" {
		redirect_url := fl.Field().Interface().(string)
		if strings.Compare(strings.Trim(redirect_url, " "), "") == 0 {
			return true
		}
		_, err := url.Parse(redirect_url)
		if err != nil {
			return false
		}
	}
	return true
}

func MyEmail(fl validator.FieldLevel) bool {
	if fl.GetTag() == "my_email" {
		email := fl.Field().Interface().(string)
		if strings.Compare(strings.Trim(email, " "), "") == 0 {
			return true
		}

		pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
		reg := regexp.MustCompile(pattern)
		return reg.MatchString(email)
	}
	return true
}
