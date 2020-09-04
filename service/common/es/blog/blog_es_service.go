package blog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/service/common/es"
	"github.com/olivere/elastic/v7"
	"log"
)

type BlogEsService struct {
	*es.BaseEsService
}

func NewBlogEsService() (*BlogEsService, error) {

	log.Println("创建blog_es_service")

	base_es_service, err := es.NewBaseEsService()

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

//更新文档
//keyword 搜索的关键字
//per_page 每页多少条
//page 第几页
func (b *BlogEsService) SearchBlog(keyword string, per_page, page int) (result *entity.ListResponseEntity) {

	search_serivce := b.Client.Search().Index(es.BLOG_INDEX).From((page - 1) * per_page).Size(per_page).Pretty(true)

	multi_match_query := elastic.NewMultiMatchQuery(keyword, []string{"title", "abstract", "content"}...)

	search_response, err := search_serivce.Query(multi_match_query).Do(context.Background())

	if err != nil {
		log.Println("searchBlog：es搜索错误：", err)
		result = nil
		return result
	}

	//构建结果返回
	result = new(entity.ListResponseEntity)
	result.SetPerPage(per_page)

	var blog_entity_list []*blog.BlogEntity

	if search_response.TotalHits() > 0 {

		result.SetCount(search_response.TotalHits())

		blog_entity_list = make([]*blog.BlogEntity, len(search_response.Hits.Hits))

		for index, hit := range search_response.Hits.Hits {

			t := new(blog.BlogEntity)
			err := json.Unmarshal(hit.Source, t)
			if err != nil {
				panic("es解析失败：" + err.Error())
			}

			blog_entity_list[index] = t
		}
	}

	result.SetList(blog_entity_list)
	return
}
