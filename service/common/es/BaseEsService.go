package es

import (
	"context"
	"github.com/blog_backend/common-lib/config"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)


const (
	BLOG_INDEX = "xiaochen_blog_next_blog"
)

type BaseEsService struct {
	//es连接的客户端
	Client *elastic.Client
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
	s.Client = client

	return s
}

//更新一个文档
func (b *BaseEsService) UpdateDoc(index, doc_id string, doc interface{}) (*elastic.UpdateResponse, error) {
	res, err := b.Client.Update().Index(index).Id(doc_id).Doc(doc).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

//删除一个文档
func (b *BaseEsService) DeleteDoc(index, doc_id string) (*elastic.DeleteResponse, error) {
	res, err := b.Client.Delete().Index(index).Id(doc_id).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

//获取一个文档
func (b *BaseEsService) GetDoc(index, doc_id string) (*elastic.GetResult, error) {
	res, err := b.Client.Get().Index(index).Id(doc_id).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}
