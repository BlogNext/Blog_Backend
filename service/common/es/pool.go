package es

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"sync"
)

//es连接池
var es_default_connet_pool *sync.Pool

func init() {

	if es_default_connet_pool == nil {

		es_default_connet_pool = &sync.Pool{New: func() interface{} {

			log.Println("获取es配置")
			es_config, err := config.GetConfig("es")
			if err != nil {
				panic(err)
			}

			es_default_info := es_config.GetStringMap("es")
			log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", es_default_info, es_default_info, es_default_info))

			host := es_default_info["default"].(map[string]interface{})["host"].(string)
			username := es_default_info["default"].(map[string]interface{})["username"].(string)
			password := es_default_info["default"].(map[string]interface{})["password"].(string)
			log.Println("es连接信息", "地址:"+host, "用户名:"+username, "密码："+password)

			client, err := elastic.NewClient(
				elastic.SetURL(host),
				elastic.SetBasicAuth(username, password),
				elastic.SetTraceLog(log.New(os.Stdout, "blog_next", 0)), //跟踪请求和响应细节
			)

			if err != nil {
				panic(err)
			}

			return client

		}}
	}
}
