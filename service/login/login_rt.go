package login

import (
	"errors"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/login/front"
	"github.com/blog_backend/entity/user"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"strings"
	"time"
)

var mySigningKey []byte

func init() {
	mySigningKey = []byte("xiaochen123")
}

//前端登录
type LoginRtService struct {
}

//是否登录
//login_token jwt token
//login_entity 登录实体，如果登录成功，会赋值信息
//return true登录，false未登录
func (u *LoginRtService) IsLogin(login_token string, login_entity *front.LoginEntity) bool {

	if login_token == "" {
		return false
	}
	
	token, err := jwt.ParseWithClaims(login_token, login_entity, func(token *jwt.Token) (i interface{}, err error) {
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
func (u *LoginRtService) LoginByYuque(login, password string) (login_token string) {
	db := mysql.GetDefaultDBConnect()
	model := new(model.UserYuQueModel)
	query_result := db.Where("login = ?", login).First(model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(exception.NewException(exception.VALIDATE_ERR, "未找到用户login:"+login))
	}

	if strings.Compare(password, "xiaochen123") != 0 {
		panic("密码不正确")
	}

	//生成jwt
	claims := &front.LoginEntity{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "ly",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		UserFrontEntity: user.UserFrontEntity{
			BaseEntity: entity.BaseEntity{
				DocID:     "",
				ID:        uint64(model.ID),
				CreatedAt: uint64(model.CreatedAt),
				UpdatedAt: uint64(model.UpdatedAt),
			},
			UserId: uint64(model.UserId),
			Login:  model.Login,
			Name:   model.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	login_token, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	return login_token
}
