package front

import (
	"github.com/blog_backend/entity/user"
	"github.com/dgrijalva/jwt-go"
)

//前端登录实体
type LoginEntity struct {
	jwt.StandardClaims
	UserFrontEntity user.UserFrontEntity `json:"user_front_entity"`
}

