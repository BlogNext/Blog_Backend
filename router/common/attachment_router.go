package common

import (
	"github.com/blog_backend/controller"
	c_common "github.com/blog_backend/controller/common"
	"github.com/gin-gonic/gin"
)

func RegisterCommontRouter(router *gin.Engine) {

	//前后端公共的路由
	common_router := router.Group("/common")
	{
		//附件
		attachemt_router := common_router.Group("/attachment")
		{
			attachment_controller := controller.NewController(new(c_common.AttachmentController))
			attachemt_router.Any("/:action", attachment_controller)
		}

	}
}
