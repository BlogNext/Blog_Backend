package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/blog_backend/exception"
	"net/url"
	"time"
)

const (
	//jwt生成的类型
	PREAUTH_CODE = iota
	ACCESS_TOKEN
	REFRESH_TOKEN
)

//jwt的秘钥
var mySigningKey []byte

//token映射表
var TokenMap []string

func init() {

	if len(mySigningKey) <= 0 {
		mySigningKey = []byte("AllYourBase")
	}

	if TokenMap == nil {
		TokenMap = make([]string, 3)
		TokenMap[PREAUTH_CODE] = "code"
		TokenMap[ACCESS_TOKEN] = "access_token"
		TokenMap[REFRESH_TOKEN] = "refresh_token"
	}
}

/**
jwt 生成和验签
*/

//JWT生成
func generateJWT(client_and_user *ClientAndUser, jwt_type int) (string, error) {
	//jwt公共属性设置
	token := new(Token)
	token.StandardClaims.Issuer = "ly"               //发布者
	token.UserId = client_and_user.OauthUser.Id      //授权用户的id
	token.ClientId = client_and_user.AcceptClient.Id //客户id
	token.ClientAppId = client_and_user.AcceptClient.ClientAppId
	token.RedirectUrl = client_and_user.RedirectUrl.String()
	switch jwt_type {
	case PREAUTH_CODE:
		token.StandardClaims.Subject = TokenMap[PREAUTH_CODE]
		token.StandardClaims.ExpiresAt = time.Now().Add(3 * time.Minute).Unix() //code 3分钟过期时间
		break
	case ACCESS_TOKEN:
		token.StandardClaims.Subject = TokenMap[ACCESS_TOKEN]
		token.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix() //token 2小时过期时间
		break
	case REFRESH_TOKEN:
		token.StandardClaims.Subject = TokenMap[REFRESH_TOKEN]
		//refresh_token没有过期时间
	default:
		panic("非法的jwt_type")
	}
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return jwt_token.SignedString(mySigningKey)
}

//验签jwt
func validateJWT(jwt_token string) (*Token, error) {
	result_token := new(Token)

	token, err := jwt.ParseWithClaims(jwt_token, result_token, func(token *jwt.Token) (i interface{}, err error) {
		return mySigningKey, nil
	})

	if err != nil {
		//验证不成功的时候,打印一下日志
		if ve, ok := err.(*jwt.ValidationError); ok {
			//断言是否为jwt的错误
			oa_exception := new(exception.BaseException)
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				//那个2的几次方按位与的算法
				fmt.Println("token过期或者还没到生效时间")
				oa_exception.SetErrorCode(exception.TOKEN_EXPIRED)
				oa_exception.SetErrorMsg("token过期或者还没到生效时间")
			} else {
				fmt.Println("非法的token错误信息: ", err)
				oa_exception.SetErrorCode(exception.TOKEN_ILLEGALITY)
				oa_exception.SetErrorMsg(fmt.Sprintln("非法的token错误信息: ", err))
			}

			return nil, oa_exception
		}

		fmt.Println("非法的token错误信息: ", err)
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return result_token, nil
}

/**
token结构体,客户和用户的中间产物
*/
type Token struct {
	//jwt传输定义一下
	//必填必填字段
	ClientId    int64  `json:"client_id,omitempty"` //客户的appid
	ClientAppId string `json:"client_app_id,omitempty"`
	RedirectUrl string `json:"redirect_url,omitempty"` //授权后重定向的地址,绝对路径
	UserId      int64  `json:"user_id,omitempty"`      //用户id
	//非必填
	PreauthCode  string `json:"preauth_code,omitempty"`  //预授权码,jwt_token
	AccessToken  string `json:"access_token,omitempty"`  //token,jwt_token
	RefreshToken string `json:"refresh_token,omitempty"` //refresh_token,jwt_token
	jwt.StandardClaims
}

//生成预授权码
func GetPreauthCode(client_and_user *ClientAndUser, account_is_valid ...ValidateAccountIsValid) (*Token, error) {

	preauth_code, err := generateJWT(client_and_user, PREAUTH_CODE)

	if err != nil {
		//生成token失败
		return nil, err
	}

	token := new(Token)
	token.PreauthCode = preauth_code
	token.RedirectUrl = client_and_user.AcceptClient.RedirectUrl.String()
	//验证账号是否有效
	err = validateAccountIsValid(token, account_is_valid...)
	if err != nil {
		return nil, err
	}

	return token, nil
}

//验证预授权码
func ValidatePreauthCode(preauth_code string) (*Token, error) {
	return validateJWT(preauth_code)
}

//获取AccessToken
func GetAccessToken(preauth_code string, client_and_user *ClientAndUser, account_is_valid ...ValidateAccountIsValid) (*Token, error) {
	//验证预授权码JWT
	_, err := ValidatePreauthCode(preauth_code)

	if err != nil {
		return nil, err
	}

	//有效的预授权码,客户、用户的账号有效
	//生成accessToken和refreshtoken
	access_token, access_err := generateJWT(client_and_user, ACCESS_TOKEN)
	refresh_token, refresh_err := generateJWT(client_and_user, REFRESH_TOKEN)

	if access_err != nil || refresh_err != nil {
		return nil, errors.New(fmt.Sprintf("生成token失败: accrss_err=%s \n refresh_err=%s", access_err, refresh_err))
	}

	token := new(Token)
	token.AccessToken = access_token
	token.RefreshToken = refresh_token

	//验证账号是否有效
	err = validateAccountIsValid(token, account_is_valid...)
	if err != nil {
		return nil, err
	}

	return token, nil
}

//验证access_token
func ValidateAccessToken(access_token string) (*Token, error) {
	return validateJWT(access_token)
}

//通过refresh_token刷新access_token
//由控制器决定是否更新数据库中的refresh_token
func FlushTokenByRefreshToken(refresh_token string, account_is_valid ...ValidateAccountIsValid) (*Token, error) {
	rf_token, err := validateJWT(refresh_token)

	if err != nil {
		//refresh_token解析失败
		return nil, err
	}

	redirect_url, err := url.Parse(rf_token.RedirectUrl)
	if err != nil {
		return nil, err
	}

	client_and_user := ClientAndUser{
		AcceptClient: AcceptClient{
			Id:          rf_token.ClientId,
			RedirectUrl: *redirect_url,
		},
		OauthUser: OauthUser{
			Id: rf_token.UserId,
		},
	}

	token, err := GetAccessToken(refresh_token, &client_and_user)

	if err != nil {
		return nil, err
	}

	//验证账号是否有效
	err = validateAccountIsValid(token, account_is_valid...)
	if err != nil {
		return nil, err
	}

	return token, nil
}

//通过回调验证账户是否有效
func validateAccountIsValid(token *Token, account_is_valid ...ValidateAccountIsValid) error {
	//验证账号是否有效
	if account_is_valid != nil {
		for _, fuc := range account_is_valid {
			err := fuc(token)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
