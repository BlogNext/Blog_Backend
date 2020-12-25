package front

import (
	log_entity "github.com/blog_backend/entity/login/front"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/login"
	"github.com/gin-gonic/gin"
)

//前端登录中间件
func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		log_token := c.GetHeader("x-access-token")
		service := new(login.LoginRtService)
		login_entity := new(log_entity.LoginEntity)

		ok := service.IsLogin(log_token, login_entity)
		if !ok {
			help.Gin200SuccessResponse(c, "请先登录", nil)
			panic("请登录")
		}
		c.Next()
	}
}
