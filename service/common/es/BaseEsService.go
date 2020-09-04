package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
)

const (
	BLOG_INDEX = "xiaochen_blog_next_blog"
)


//es基础服务对象
type BaseEsService struct {
	//es连接的客户端
	Client *elastic.Client
}

func NewBaseEsService() (s *BaseEsService, err error) {

	//异常捕获一下
	defer func() {
		if unknown_err := recover(); unknown_err != nil {
			log.Println("链接es失败：", unknown_err)
			s = nil
			err = unknown_err.(error)
		}
	}()

	client := es_default_connet_pool.Get()
	s = new(BaseEsService)
	s.Client = client.(*elastic.Client)

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
