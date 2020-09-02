package es

import (
	"context"
	"fmt"
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

func NewBaseEsService(host string, username string, password string) (s *BaseEsService, err error) {

	//异常捕获一下
	defer func() {
		if unknown_err := recover(); unknown_err != nil {
			log.Println("链接es失败：", unknown_err)
			s = nil
			err = unknown_err.(error)
		}
	}()

	if len(host) <= 0 || len(username) <= 0 || len(password) <= 0 {
		es_config, err := config.GetConfig("es")
		if err != nil {
			panic(err)
		}

		log.Println("获取es配置")

		es_default_info := es_config.GetStringMap("es")
		log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", es_default_info, es_default_info, es_default_info))
		host = es_default_info["default"].(map[string]interface{})["host"].(string)
		username = es_default_info["default"].(map[string]interface{})["username"].(string)
		password = es_default_info["default"].(map[string]interface{})["password"].(string)
	}

	log.Println("es连接信息", "地址:"+host, "用户名:"+username, "密码："+password)

	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetBasicAuth(username, password),
		elastic.SetTraceLog(log.New(os.Stdout, "blog_next", 0)), //跟踪请求和响应细节
	)

	if err != nil {
		panic(err)
	}

	s = new(BaseEsService)
	s.Client = client

	return s, nil
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
