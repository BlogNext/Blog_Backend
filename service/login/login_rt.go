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
	claims := &front.LoginEntity{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "ly",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		UserFrontEntity: user.UserFrontEntity{
			BaseEntity: entity.BaseEntity{
				DocID:     "",
				ID:        uint64(userYuQueModel.ID),
				CreatedAt: uint64(userYuQueModel.CreatedAt),
				UpdatedAt: uint64(userYuQueModel.UpdatedAt),
			},
			UserId: uint64(userYuQueModel.UserId),
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
