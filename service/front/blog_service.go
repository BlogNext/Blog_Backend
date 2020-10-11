package front

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/backend"
	"github.com/blog_backend/service/common/es/blog"
	"log"
)

type BlogService struct {
}

//searchLevel 搜索等级
//keyword 搜索的关键字
//per_page 每页多少条
//page 第几页
func (b *BlogService) SearchBlog(searchLevel string, keyword string, per_page, page int) (result *entity.ListResponseEntity) {

	switch searchLevel {
	case MYSQL_SEARCH_LEVEL:
		result = b.SearchBlogMysqlLevel(keyword, per_page, page)
	case ES_SEARCH_LEVEL:
		//es搜索
		blog_s := new(blog.BlogEsService)
		result = blog_s.SearchBlog(keyword, per_page, page)

		if result == nil {
			//降级为mysql搜索
			result = b.SearchBlogMysqlLevel(keyword, per_page, page)
		}
	}

	return
}

//mysql等级搜索博客
func (b *BlogService) SearchBlogMysqlLevel(keyword string, per_page, page int) (result *entity.ListResponseEntity) {

	log.Println("进入mysql搜索")

	var blog_model_list []*model.BlogModel
	var count int64

	db := mysql.GetDefaultDBConnect()
	db = db.Table(model.BlogModel{}.TableName()).Where("content like ? OR title like ?", "%"+keyword+"%", "%"+keyword+"%")
	db.Count(&count)
	db.Limit(per_page).Offset((page - 1) * per_page).Find(&blog_model_list)

	log.Println("总数:", count, "数据:", blog_model_list, "数据长度:", len(blog_model_list))

	//转化为传输层的对象
	bs := new(backend.BlogService)
	blog_entity_list := bs.ChangeToBlogEntityFormList(blog_model_list)

	//构建结果返回
	result = new(entity.ListResponseEntity)

	result.SetCount(count)
	result.SetPerPage(per_page)
	result.SetList(blog_entity_list)

	return result
}
