package blog

import (
	"context"
	"errors"
	"fmt"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/service/common/es"
	"github.com/olivere/elastic/v7"
	"log"
)

type BlogEsService struct {
	*es.BaseEsService
}

func NewBlogEsService(host string, username string, password string) (*BlogEsService, error) {

	log.Println("创建blog_es_service")

	base_es_service, err := es.NewBaseEsService(host, username, password)

	if err != nil {
		return nil, err
	}

	return &BlogEsService{
		BaseEsService: base_es_service,
	}, nil
}

//删除blog文档
func (b *BlogEsService) DeleteDoc(blog_doc *blog.BlogEntity) () {

	_, err := b.BaseEsService.DeleteDoc(es.BLOG_INDEX, blog_doc.DocID)

	if err != nil {
		panic(err)
	}

}

//添加一个doc,返回文档在es中的唯一标识
func (b *BlogEsService) AddDoc(blog_doc *blog.BlogEntity) *elastic.IndexResponse {

	if blog_doc == nil {
		panic(errors.New("blog_doc为空"))
	}

	doc, err := b.Client.Index().Index(es.BLOG_INDEX).BodyJson(blog_doc).Do(context.Background())
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("往索引%s添加文档成功,自动生成:ID%s,版本号是:%d", doc.Index, doc.Id, doc.Version))
	return doc
}

//更新一个文档的内容
func (b *BlogEsService) UpdateDoc(blog_doc *blog.BlogEntity) *elastic.UpdateResponse {

	if blog_doc == nil {
		panic(errors.New("blog_doc为空"))
	}

	doc, err := b.BaseEsService.UpdateDoc(es.BLOG_INDEX, blog_doc.DocID, blog_doc)

	if err != nil {
		panic(err)
	}

	return doc
}
