package blog

import (
	"context"
	"errors"
	"fmt"
	"github.com/blog_backend/service/common/es"
	"log"
)

type BlogEsService struct {
	*es.BaseEsService
}

func NewBlogEsService(host string, username string, password string) *BlogEsService {
	return &BlogEsService{
		BaseEsService: es.NewBaseEsService(host, username, password),
	}
}

//导入blog数据
func (b *BlogEsService) ImportData() {

}

//删除blog文档的所有数据
func (b *BlogEsService) DeleteData() {

}

//添加一个doc,返回文档在es中的唯一标识
func (b *BlogEsService) AddDoc(blog_doc *BlogDoc) string {

	if blog_doc == nil {
		panic(errors.New("blog_doc为空"))
	}

	doc, err := b.Client.Index().Index(es.BLOG_INDEX).BodyJson(blog_doc).Do(context.Background())
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("往索引%s添加文档成功,自动生成:ID%s,版本号是:%d", doc.Index, doc.Id, doc.Version))
	return doc.Id
}

//更新一个文档的内容
func (b *BlogEsService) UpdateDoc(doc_id string, blog_doc *BlogDoc) {

}