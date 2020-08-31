package es

import (
	"github.com/blog_backend/common-lib/config"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

type BaseEsService struct {
	client *elastic.Client
}

func NewBaseEsService(host string, username string, password string) *BaseEsService {

	if len(host) <= 0 || len(username) <= 0 || len(password) <= 0 {
		es_config, _ := config.GetConfig("es")
		es_default_info := es_config.GetStringMap("default")
		host = es_default_info["host"].(string)
		username = es_default_info["username"].(string)
		password = es_default_info["password"].(string)
	}

	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetBasicAuth(username, password),
		elastic.SetTraceLog(log.New(os.Stdout, "blog_next", 0)), //跟踪请求和响应细节
	)

	if err != nil {
		panic(err)
	}

	s := new(BaseEsService)
	s.client = client

	return s
}
