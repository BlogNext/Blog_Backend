package login

import (
	"errors"
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/common-lib/oauth_sso/core"
	"github.com/blog_backend/common-lib/oauth_sso/oauth"
	"github.com/blog_backend/common-lib/oauth_sso/token"
	ssoUser "github.com/blog_backend/common-lib/oauth_sso/user"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/login/front"
	"github.com/blog_backend/entity/user"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

//jwt配置
var mySigningKey []byte

//oauthSSO配置
//oauthSSo服务器配置
var oauthSSOConfig map[string]string

//oauthSSo客户配置
var oauthSSOClientConfig map[string]string

func init() {
	//jwt配置初始化
	mySigningKey = []byte("xiaochen123")
	//oauthSSO配置初始化
	if oauthSSOConfig == nil {
		config.LoadConfig("server", "config", "yaml")

		//运行服务器
		serverConfig, err := config.GetConfig("server")
		if err != nil {
			log.Fatalln(err,"配置信息加载失败")
		}

		oauthSSOMapConfig := serverConfig.GetStringMap("oauthSSO")
		oauthSSOConfig = make(map[string]string)
		oauthSSOConfig["scheme"] = oauthSSOMapConfig["scheme"].(string)
		oauthSSOConfig["host"] = oauthSSOMapConfig["host"].(string)
		core.SetOauthSSoSchemeConfig(oauthSSOConfig["scheme"])
		core.SetOauthSSoHostConfig(oauthSSOConfig["host"])

		//客户配置
		oauthSSOClientConfig = make(map[string]string)
		oauthSSOClientConfig["clientId"] = oauthSSOMapConfig["client_id"].(string)
		oauthSSOClientConfig["clientSecret"] = oauthSSOMapConfig["client_secret"].(string)


	}

}

//前端登录
type LoginRtService struct {
}

//是否登录
//loginToken jwt token
//loginEntity 登录实体，如果登录成功，会赋值信息
//return true登录，false未登录
func (u *LoginRtService) IsLogin(loginToken string, loginEntity *front.LoginEntity) bool {

	if loginToken == "" {
		return false
	}

	token, err := jwt.ParseWithClaims(loginToken, loginEntity, func(token *jwt.Token) (i interface{}, err error) {
		return mySigningKey, nil
	})

	if err != nil {
		panic(err)
	}

	if _, ok := token.Claims.(*front.LoginEntity); ok && token.Valid {
		return true
	}

	return false
}

//语雀登录
//login语雀的login
//password登录密码
func (u *LoginRtService) LoginByYuque(login, password string) (loginToken string) {
	db := mysql.GetNewDB(false)
	userYuQueModel := new(model.UserYuQueModel)
	queryResult := db.Where("login = ?", login).First(userYuQueModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(exception.NewException(exception.VALIDATE_ERR, "未找到用户login:"+login))
	}

	if strings.Compare(password, "xiaochen123") != 0 {
		panic(exception.NewException(exception.VALIDATE_ERR, "密码不正确"))
	}

	//生成jwt
	return u.buildToken(userYuQueModel)
}

//jwt生成登录Token
func (u *LoginRtService) buildToken(userYuQueModel *model.UserYuQueModel) string {

	//生成jwt
	claims := &front.LoginEntity{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "ly",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		UserFrontEntity: user.UserFrontEntity{
			BaseEntity: entity.BaseEntity{
				DocID:     "",
				ID:        userYuQueModel.ID,
				CreatedAt: uint64(userYuQueModel.CreatedAt),
				UpdatedAt: uint64(userYuQueModel.UpdatedAt),
			},
			UserId: userYuQueModel.UserId,
			Login:  userYuQueModel.Login,
			Name:   userYuQueModel.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	loginToken, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	return loginToken
}

//sso登录
//返回登录的token
func (u *LoginRtService) BlogNextPreCode(request *front.BlogNextPreCodeRequest) string {
	manage := oauth.NewManage(oauthSSOClientConfig["clientId"], oauthSSOClientConfig["clientSecret"])
	//预授权码换取accessToken
	pacatr := new(oauth.PreAuthCodeAccessTokenResponse)
	err := manage.PreAuthCodeAccessToken(request.PreCode, pacatr)
	if err != nil {
		panic(exception.NewException(exception.VALIDATE_ERR, err.Error()))
	}

	//accessToken换取用户信息，实现登录
	tokenManage := token.NewTokenManage(pacatr.RefreshToken, oauthSSOClientConfig["clientId"], oauthSSOClientConfig["clientSecret"])
	userManage := ssoUser.NewManage(tokenManage)
	uir := new(ssoUser.UserInfoResponse)
	err = userManage.UserInfo(uir)
	if err != nil {
		//获取用户信息失败
		panic(exception.NewException(exception.VALIDATE_ERR, err.Error()))
	}

	db := mysql.GetNewDB(false)
	userYuQueModel := new(model.UserYuQueModel)
	queryResult := db.Where("user_id = ?", uir.Id).First(userYuQueModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(exception.NewException(exception.VALIDATE_ERR, fmt.Sprintf("未找到用户user_id:%d", uir.Id)))
	}

	return u.buildToken(userYuQueModel)
}
