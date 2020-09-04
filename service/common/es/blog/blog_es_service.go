package blog

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/service/common/es"
	"github.com/olivere/elastic/v7"
	"log"
)

type BlogEsService struct {
	es.BaseEsService
}

//删除blog文档
func (b *BlogEsService) DeleteDoc(blog_doc *blog.BlogEntity) *elastic.DeleteResponse {
	//构建一个命令
	commend := es.BuildDeleteDocCommend(es.BLOG_INDEX, blog_doc.DocID)
	//设置命令
	b.SetExecCommend(commend)
	//运行命令
	resule, err := es.RunCommend(b)

	if err != nil {
		return nil
	}

	return resule.(*elastic.DeleteResponse)
}

//添加一个doc,返回文档在es中的唯一标识
func (b *BlogEsService) AddDoc(blog_doc *blog.BlogEntity) *elastic.IndexResponse {

	if blog_doc == nil {
		panic(errors.New("blog_doc为空"))
	}

	//构建一个命令
	commend := es.BuildAddDocCommend(es.BLOG_INDEX, blog_doc.DocID)
	//设置命令
	b.SetExecCommend(commend)
	//运行命令
	resule, err := es.RunCommend(b)

	if err != nil {
		return nil
	}

	return resule.(*elastic.IndexResponse)
}

//更新一个文档的内容
func (b *BlogEsService) UpdateDoc(blog_doc *blog.BlogEntity) *elastic.UpdateResponse {

	if blog_doc == nil {
		panic(errors.New("blog_doc为空"))
	}

	//构建一个命令
	commend := es.BuildUpdateDocCommend(es.BLOG_INDEX, blog_doc.DocID, blog_doc)
	//设置命令
	b.SetExecCommend(commend)
	//运行命令
	resule, err := es.RunCommend(b)

	if err != nil {
		return nil
	}

	return resule.(*elastic.UpdateResponse)
}

//keyword 搜索的关键字
//per_page 每页多少条
//page 第几页

func (b *BlogEsService) SearchBlog(keyword string, per_page, page int) (result *entity.ListResponseEntity) {
	//构建一个命令
	commend := es.BuildCallbackCommend(func() es.Commend {
		return func(client *elastic.Client) (i interface{}, err error) {

			search_serivce := client.Search().Index(es.BLOG_INDEX).From((page - 1) * per_page).Size(per_page).Pretty(true)

			multi_match_query := elastic.NewMultiMatchQuery(keyword, []string{"title", "abstract", "content"}...)

			search_result, err := search_serivce.Query(multi_match_query).Do(context.Background())

			if err != nil {
				log.Println("searchBlog：es搜索错误：", err)
				return nil, err
			}

			return search_result, nil
		}
	})

	//设置命令
	b.SetExecCommend(commend)

	//运行命令
	commend_result, err := es.RunCommend(b)

	if err != nil {
		return nil
	}

	search_result := commend_result.(*elastic.SearchResult)

	//构建结果返回
	result = new(entity.ListResponseEntity)
	result.SetPerPage(per_page)

	var blog_entity_list []*blog.BlogEntity

	if search_result.TotalHits() > 0 {

		result.SetCount(search_result.TotalHits())

		blog_entity_list = make([]*blog.BlogEntity, len(search_result.Hits.Hits))

		for index, hit := range search_result.Hits.Hits {

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
