package blog

import (
	"context"
	"encoding/json"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/service/es"
	"github.com/olivere/elastic/v7"
	"log"
)

type BlogEsRtService struct {
	es.BaseEsService
}

//keyword 搜索的关键字
//per_page 每页多少条
//page 第几页

func (b *BlogEsRtService) SearchBlog(keyword string, per_page, page int) (result *entity.ListResponseEntity) {

	search_result, err := b.BaseEsService.SearchDoc(func() es.Commend {
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

	if err != nil {
		return nil
	}

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
