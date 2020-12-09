package blog

import (
	"errors"
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
	"strings"
)

//搜索等级
const (
	MYSQL_SEARCH_LEVEL = "mysql"
	ES_SEARCH_LEVEL    = "es"
)

//博客前台服务
type BlogRtService struct {
}

//浏览量自增
func (s *BlogRtService) IncBrowse(id uint) {
	db := mysql.GetDefaultDBConnect()
	db.Model(model.BlogModel{}).Where("id = ?", id).UpdateColumn("browse_total", gorm.Expr("browse_total + ?", 1))
}

//博客详情
func (s *BlogRtService) Detail(id uint) *blog.BlogEntity {
	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	query_result := db.First(blog_model, id)

	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到博客:%d", id))
	}

	result := make([]*blog.BlogEntity, 1)
	result[0] = ChangeToBlogEntity(blog_model)

	return result[0]
}

//列表页
func (s *BlogRtService) GetList(filter map[string]string, per_page, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetDefaultDBConnect()

	blog_table_name := model.BlogModel{}.TableName()

	//博客需要的字段
	blog_felid := []string{"id", "user_id", "blog_type_id", "cover_plan_id", "title", "abstract", "browse_total", "created_at", "updated_at"}

	for index, felid := range blog_felid {
		blog_felid[index] = fmt.Sprintf("%s.%s", blog_table_name, felid)
	}

	var count int64
	//sql
	db = db.Table(blog_table_name)

	//过滤分类id过滤
	if filter["blog_type_id"] != "" {
		db = db.Where("blog_type_id = ?", filter["blog_type_id"])
	}

	db.Count(&count)

	rows, err := db.Select(strings.Join(blog_felid, ", ")).Order("created_at DESC").Limit(per_page).Offset((page - 1) * per_page).Rows()

	if err != nil {
		return nil
	}

	query_result := make([]*blog.BlogListEntity, 0)

	cover_plan_ids := make([]uint64, 0)
	blog_type_ids := make([]uint64, 0)
	user_id_ids := make([]uint, 0)

	for rows.Next() {
		var id uint64
		var user_id uint
		var blog_type_id uint64
		var cover_plan_id uint64
		var title string
		var abstract string
		var browse_total uint
		var created_at uint64
		var updated_at uint64
		rows.Scan(&id, &user_id, &blog_type_id, &cover_plan_id, &title, &abstract, &browse_total, &created_at, &updated_at)

		//博客实体
		blog_entity := new(blog.BlogListEntity)
		blog_entity.ID = id
		blog_entity.UserId = uint64(user_id)
		blog_entity.BlogTypeId = blog_type_id
		blog_entity.CoverPlanId = cover_plan_id
		blog_entity.Title = title
		blog_entity.Abstract = abstract
		blog_entity.BrowseTotal = browse_total
		blog_entity.CreatedAt = created_at
		blog_entity.UpdatedAt = updated_at

		cover_plan_ids = append(cover_plan_ids, cover_plan_id)
		blog_type_ids = append(blog_type_ids, blog_type_id)
		user_id_ids = append(user_id_ids, user_id)

		query_result = append(query_result, blog_entity)
	}

	//填充信息
	PaddingAttachemtInfoByBlogListEntity(cover_plan_ids, query_result) //填充附件信息
	PaddingBlogTypeInfoByBlogListEntity(blog_type_ids, query_result)   //博客类型实体
	PaddingUserInfoByBlogListEntity(user_id_ids, query_result)         //填充用户信息

	//构建结果返回
	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(per_page)
	result.SetList(query_result)

	return
}

//排序方向
//per_page 每页多少条
func (s *BlogRtService) GetListBySort(sort_dimension string, per_page int) (result *entity.ListResponseEntity) {

	db := mysql.GetDefaultDBConnect()
	var blog_model_list []*model.BlogModel
	//查询
	switch sort_dimension {
	case "browse_total":
		db = db.Table(model.BlogModel{}.TableName())
		db.Order("browse_total DESC").Limit(per_page).Find(&blog_model_list)
	case "created_at":
		db = db.Table(model.BlogModel{}.TableName())
		db.Order("created_at DESC").Limit(per_page).Find(&blog_model_list)
	default:
		exception.NewException(exception.VALIDATE_ERR, "非法的sort_dimension")
	}

	//转化为传输层的对象
	blog_sort_entity_list := ChangeBlogSortEntityByList(blog_model_list)

	//构建结果返回
	result = new(entity.ListResponseEntity)
	filter_list := []help.Filter{
		help.Filter{
			Label: "排序维度",
			Field: "sort_dimension",
			Options: []help.Option{
				help.Option{
					Label: "浏览量",
					Value: "browse_total",
				},
				help.Option{
					Label: "创建时间",
					Value: "created_at",
				},
			},
		},
	}
	result.SetFilter(filter_list)
	result.SetCount(int64(per_page))
	result.SetPerPage(per_page)
	result.SetList(blog_sort_entity_list)

	return result
}

//searchLevel 搜索等级
//keyword 搜索的关键字
//per_page 每页多少条
//page 第几页
func (s *BlogRtService) SearchBlog(searchLevel string, keyword string, per_page, page int) (result *entity.ListResponseEntity) {

	switch searchLevel {
	case MYSQL_SEARCH_LEVEL:
		result = s.SearchBlogMysqlLevel(keyword, per_page, page)
		//case ES_SEARCH_LEVEL:
		//	//es搜索
		//	blog_s := new(BlogEsRtService)
		//	result = blog_s.SearchBlog(keyword, per_page, page)
		//
		//	if result == nil {
		//		//降级为mysql搜索
		//		result = s.SearchBlogMysqlLevel(keyword, per_page, page)
		//	}
	}

	return
}

//mysql等级搜索博客
func (s *BlogRtService) SearchBlogMysqlLevel(keyword string, per_page, page int) (result *entity.ListResponseEntity) {

	var blog_model_list []*model.BlogModel
	var count int64

	db := mysql.GetDefaultDBConnect()
	db = db.Table(model.BlogModel{}.TableName())
	if keyword != "" {
		db = db.Where("content like ? OR title like ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	db.Count(&count)
	db.Order("created_at DESC").Limit(per_page).Offset((page - 1) * per_page).Find(&blog_model_list)

	//转化为传输层的对象
	blog_list_entity_list := ChangeToBlogListEntityList(blog_model_list)

	//构建结果返回
	result = new(entity.ListResponseEntity)

	result.SetCount(count)
	result.SetPerPage(per_page)
	result.SetList(blog_list_entity_list)

	return result
}

//blogInfo模块统计展示
func (s *BlogRtService) GetStat() (result *entity.ListResponseEntity) {

	response := make(map[string]uint, 3)

	db := mysql.GetDefaultDBConnect()

	blog_table_name := model.BlogModel{}.TableName()
	db = db.Table(blog_table_name)

	//文章总数
	{
		var count int64
		//sql
		count_db := db.Select("id")
		count_db.Count(&count)
		response["blog_total"] = uint(count)
	}

	//最新一篇文章的时间、最新和最早的文章时间戳
	{
		//最新的博客时间
		var last_create_at uint
		last_create_at_row := db.Select("created_at").Order("id DESC").Limit(1).Row()
		last_create_at_row.Scan(&last_create_at)
		response["last_create_at"] = last_create_at

		//最老的博客时间
		var first_create_at uint
		first_create_at_db := mysql.GetDefaultDBConnect()
		first_create_at_row := first_create_at_db.Table(blog_table_name).Select("created_at").
			Order("id ASC").Limit(1).Row()
		first_create_at_row.Scan(&first_create_at)
		response["diff_time"] = last_create_at - first_create_at
	}

	//构建结果返回
	result = new(entity.ListResponseEntity)

	result.SetList(response)

	return result
}
