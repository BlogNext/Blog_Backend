package service

import (
	"fmt"
	"github.com/blog_backend/common-lib/crypto"
	"net/url"
	"strconv"
	"time"
)

//生成ClientId
func GenerateClientAppId(now_time time.Time) (string, error) {
	//设置时区
	var cstSh, err = time.LoadLocation("Asia/Shanghai") //上海
	now_time = time.Now()
	now_time.In(cstSh)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Oauth%s%d", now_time.Format("20060102"), now_time.Unix()), nil
}

func GenerateClientAppSecret(now_time time.Time) (string, error) {
	//设置时区
	var cstSh, err = time.LoadLocation("Asia/Shanghai") //上海
	now_time.In(cstSh)
	if err != nil {
		return "", err
	}
	unix_str := strconv.FormatInt(now_time.Unix(), 10)

	return crypto.SHA1(unix_str), nil
}

//接受授权的客户
type AcceptClient struct {
	Id              int64   //id
	ClientAppId     string  //客户的appid
	ClientAppSecret string  //客户的secret
	RedirectUrl     url.URL //重定向可以获取预授权码的地址
}

//授权的用户
type OauthUser struct {
	Id int64
}

//客户和用户的关系
type ClientAndUser struct {
	ID           int64
	AcceptClient //客户
	OauthUser    //用户
}

//验证账号是否有效
type ValidateAccountIsValid func(token *Token) error
