package oauth

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/controller"
	"github.com/blog_backend/controller/oauth/service"
	validate_request "github.com/blog_backend/controller/oauth/validate-request"
	"github.com/blog_backend/help"
	"github.com/blog_backend/model"
	mode_oauth "github.com/blog_backend/model/oauth"
	"log"
	"net/url"
	"time"
)

type OauthController struct {
	controller.BaseController
}

func (o *OauthController) Test() {
	now_time := time.Now()
	client_id, _ := service.GenerateClientAppId(now_time)
	client_secret, _ := service.GenerateClientAppSecret(now_time)
	data := make(map[string]string)
	data["client_id"] = client_id
	data["client_secret"] = client_secret
	help.Gin200SuccessResponse(o.Ctx, "test", data)
}

//注册接入客户
func (o *OauthController) RegisterClinet() {
	now_time := time.Now()
	ClientId, _ := service.GenerateClientAppId(now_time)
	ClientSecret, _ := service.GenerateClientAppSecret(now_time)
	client_name := o.Ctx.PostForm("client_name")
	redirect_url := o.Ctx.PostForm("redirect_url")
	client := &mode_oauth.OauthClient{
		ClientName:      client_name,
		ClientAppId:     ClientId,
		ClientAppSecret: ClientSecret,
		RedirectUrl:     redirect_url,
		Year:            now_time.Year(),
		Month:           now_time.Month(),
		BaseModel: model.BaseModel{
			CreatedAt: now_time.Unix(),
			UpdatedAt: now_time.Unix(),
		},
	}

	db := mysql.GetDefaultDBConnect()

	log.Println(fmt.Sprintln("主键为空返回true", db.NewRecord(client)))
	db.Create(client)
	log.Println(fmt.Sprintln("创建后返回false", db.NewRecord(client)))

	help.Gin200SuccessResponse(o.Ctx, "", nil)

	return
}

//创建授权码
func (o *OauthController) CreatePreauthCode() {

	//验证请求参数
	var preauth_code_request validate_request.PreauthCodeRequest
	o.ShouldBindWith(&preauth_code_request, binding.Query)

	log.Println("打印一下绑定的参数")
	log.Println(preauth_code_request.ClientAppId)
	log.Println(preauth_code_request.RedirectUrl)
	//获取预授权码
	redirect_url, err := url.Parse(preauth_code_request.RedirectUrl)
	if err != nil {
		panic(err)
	}

	log.Println("路径", redirect_url.String())

	client_and_user := &service.ClientAndUser{
		ID: 0,
		AcceptClient: service.AcceptClient{
			Id:          0,
			ClientAppId: preauth_code_request.ClientAppId,
			RedirectUrl: *redirect_url,
		},
		OauthUser: service.OauthUser{
			Id: 0,
		},
	}

	preauth_code, err := service.GetPreauthCode(client_and_user, nil...)

	if err != nil {
		panic(err)
	}

	help.Gin200SuccessResponse(o.Ctx, "成功获取预授权码", preauth_code)
}

//预授权码模式颁发token
func (o *OauthController) PreauthCodeExchangeAuthorizerAccessToken() {
	var preauth_code_change_token validate_request.PreauthCodeChangeTokenRequest
	o.ShouldBindWith(&preauth_code_change_token, binding.Query)

	client_and_user := &service.ClientAndUser{
		ID: 0,
		AcceptClient: service.AcceptClient{
			ClientAppId:     preauth_code_change_token.ClientAppId,
			ClientAppSecret: preauth_code_change_token.ClientAppSecret,
		},
		OauthUser: service.OauthUser{},
	}
	//获取token
	token, err := service.GetAccessToken(preauth_code_change_token.PreauthCode, client_and_user, nil...)

	if err != nil {
		panic(err)
	}

	help.Gin200SuccessResponse(o.Ctx, "换取access_token成功", token)
}

//简化模式颁发token，通过用户的账号密码以及客户的账号密码，颁发token
func (o *OauthController) ImplicitGrantExchangeAuthorizerAccessToken() {
	//验证请求参数
	var preauth_code_request validate_request.PreauthCodeRequest
	o.ShouldBindWith(&preauth_code_request, binding.Query)

	redirect_url, err := url.Parse(preauth_code_request.RedirectUrl)
	if err != nil {
		panic(err)
	}

	//获取预授权码
	client_and_user := &service.ClientAndUser{
		ID: 0,
		AcceptClient: service.AcceptClient{
			Id:          0,
			ClientAppId: preauth_code_request.ClientAppId,
			RedirectUrl: *redirect_url,
		},
		OauthUser: service.OauthUser{
			Id: 0,
		},
	}

	preauth_code, err := service.GetPreauthCode(client_and_user, nil...)

	if err != nil {
		panic(err)
	}

	//获取token
	token, err := service.GetAccessToken(preauth_code.PreauthCode, client_and_user, nil...)

	if err != nil {
		panic(err)
	}

	help.Gin200SuccessResponse(o.Ctx, "换取access_token成功", token)
}

//refresh_token刷新token
func (o *OauthController) RefreshToken() {
	//验证请求参数
	var refresh_token validate_request.RefreshTokenRequest
	o.ShouldBindWith(&refresh_token, binding.Query)
	//获取token
	token, err := service.FlushTokenByRefreshToken(refresh_token.RefreshToken, nil...)

	if err != nil {
		panic(err)
	}

	help.Gin200SuccessResponse(o.Ctx, "refresh_token刷新token成功", token)

}

//校验AccessToken，以后这个方法可以写成rpc服务，毕竟这不是对外的服务，
func (o *OauthController) vilidateAccessToken() {
	//验证请求参数
	var validate_token validate_request.ValidateTokenRequest
	o.ShouldBindWith(&validate_token, binding.Query)
	//获取token
	_, err := service.ValidateAccessToken(validate_token.Token)

	if err != nil {
		panic(err)
	}
	help.Gin200SuccessResponse(o.Ctx, "恭喜你，验证token成功", nil)
}
