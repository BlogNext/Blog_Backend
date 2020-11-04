package blog

import (
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/service/es"
	"github.com/olivere/elastic/v7"
)

type BlogEsBkService struct {
	es.BaseEsService
}

//删除blog文档
func (b *BlogEsBkService) DeleteDoc(blog_doc *blog.BlogEntity) *elastic.DeleteResponse {
	result, err := b.BaseEsService.DeleteDoc(es.BLOG_INDEX, blog_doc.DocID)
	if err != nil {
		return nil
	}
	return result
}

//添加一个doc,返回文档在es中的唯一标识
func (b *BlogEsBkService) AddDoc(blog_doc *blog.BlogEntity) *elastic.IndexResponse {
	result, err := b.BaseEsService.AddDoc(es.BLOG_INDEX, blog_doc)
	if err != nil {
		return nil
	}
	return result
}

//更新一个文档的内容
func (b *BlogEsBkService) UpdateDoc(blog_doc *blog.BlogEntity) *elastic.UpdateResponse {
	result, err := b.BaseEsService.UpdateDoc(es.BLOG_INDEX, blog_doc.DocID, blog_doc)
	if err != nil {
		return nil
	}
	return result
}
