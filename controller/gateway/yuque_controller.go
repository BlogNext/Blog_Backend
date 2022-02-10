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
	"log"
)

type YuqueController struct {
	controller.BaseController
}

func (c *YuqueController) WebHook() {
	data, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	log.Println("语雀同步过来的数据", string(data))

	//数据解码
	yuqueWebhookData := new(response.ResponseDocDetailSerializer)
	err = json.Unmarshal(data, yuqueWebhookData)
	if err != nil {
		panic(err)
	}

	yuqueConfig, err := config.GetConfig("yuque")

	if err != nil {
		panic(fmt.Sprintf("语雀配置失败"))
	}
	yuqueInfo := yuqueConfig.GetStringMap("yuque")

	yuque.SyncData(yuqueWebhookData, yuqueInfo["token"].(string))

	help.Gin200SuccessResponse(c.Ctx, "WebHook触发完成", nil)

	return
}
