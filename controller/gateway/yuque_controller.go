package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/controller"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/yuque"
	"io/ioutil"
)

type YuqueController struct {
	controller.BaseController
}

func (c *YuqueController) WebHook() {
	data, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	//数据解码
	yuque_webhook_data := new(response.ResponseDocDetailSerializer)
	err = json.Unmarshal(data, yuque_webhook_data)
	if err != nil {
		panic(err)
	}

	yuque_config, err := config.GetConfig("yuque")

	if err != nil {
		panic(fmt.Sprintf("语雀配置失败"))
	}
	yuque_info := yuque_config.GetStringMap("yuque")

	yuque.SyncData(yuque_webhook_data, yuque_info["token"].(string))

	help.Gin200SuccessResponse(c.Ctx, "WebHook触发完成", nil)

	return
}
