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
	"log"
	"strings"
)

//搜索等级
const (
	MysqlSearchLevel = "mysql"
	EsSearchLevel    = "es"
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
	blogModel := new(model.BlogModel)
	queryResult := db.First(blogModel, id)

	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到博客:%d", id))
	}

	result := make([]*blog.BlogEntity, 1)
	result[0] = ChangeToBlogEntity(blogModel)

	return result[0]
}

//获取博客列表，用于私人空间
func (s *BlogRtService) GetListByPerson(perPage, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetDefaultDBConnect()

	blogTableName := model.BlogModel{}.TableName()

	//博客需要的字段
	blogFelid := []string{"id", "user_id", "blog_type_id", "cover_plan_id", "title", "abstract", "browse_total", "created_at", "updated_at"}

	for index, felid := range blogFelid {
		blogFelid[index] = fmt.Sprintf("%s.%s", blogTableName, felid)
	}

	var count int64
	//sql
	db = db.Table(blogTableName)

	//私密博客过滤
	db = db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_0)

	db.Count(&count)

	rows, err := db.Select(strings.Join(blogFelid, ", ")).Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Rows()

	if err != nil {
		return nil
	}

	queryResult := make([]*blog.BlogListEntity, 0)

	coverPlanIds := make([]uint64, 0)
	blogTypeIds := make([]uint64, 0)
	userIdIds := make([]uint, 0)

	for rows.Next() {
		var id uint64
		var userId uint
		var blogTypeId uint64
		var coverPlanId uint64
		var title string
		var abstract string
		var browseTotal uint
		var createdAt uint64
		var updatedAt uint64
		rows.Scan(&id, &userId, &blogTypeId, &coverPlanId, &title, &abstract, &browseTotal, &createdAt, &updatedAt)

		//博客实体
		blogEntity := new(blog.BlogListEntity)
		blogEntity.ID = id
		blogEntity.UserId = uint64(userId)
		blogEntity.BlogTypeId = blogTypeId
		blogEntity.CoverPlanId = coverPlanId
		blogEntity.Title = title
		blogEntity.Abstract = abstract
		blogEntity.BrowseTotal = browseTotal
		blogEntity.CreatedAt = createdAt
		blogEntity.UpdatedAt = updatedAt

		coverPlanIds = append(coverPlanIds, coverPlanId)
		blogTypeIds = append(blogTypeIds, blogTypeId)
		userIdIds = append(userIdIds, userId)

		queryResult = append(queryResult, blogEntity)
	}

	//填充信息
	PaddingAttachemtInfoByBlogListEntity(coverPlanIds, queryResult) //填充附件信息
	PaddingBlogTypeInfoByBlogListEntity(blogTypeIds, queryResult)   //博客类型实体
	PaddingUserInfoByBlogListEntity(userIdIds, queryResult)         //填充用户信息

	//构建结果返回
	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(perPage)
	result.SetList(queryResult)

	return
}

//列表页
func (s *BlogRtService) GetList(filter map[string]string, perPage, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetDefaultDBConnect()

	blogTableName := model.BlogModel{}.TableName()

	//博客需要的字段
	blogFelid := []string{"id", "user_id", "blog_type_id", "cover_plan_id", "title", "abstract", "browse_total", "created_at", "updated_at"}

	for index, felid := range blogFelid {
		blogFelid[index] = fmt.Sprintf("%s.%s", blogTableName, felid)
	}

	var count int64
	//sql
	db = db.Table(blogTableName)

	db = db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	//过滤分类id过滤
	if filter["blog_type_id"] != "" {
		db = db.Where("blog_type_id = ?", filter["blog_type_id"])
	}

	db.Count(&count)

	rows, err := db.Select(strings.Join(blogFelid, ", ")).Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Rows()

	if err != nil {
		return nil
	}

	queryResult := make([]*blog.BlogListEntity, 0)

	coverPlanIds := make([]uint64, 0)
	blogTypeIds := make([]uint64, 0)
	userIdIds := make([]uint, 0)

	for rows.Next() {
		var id uint64
		var userId uint
		var blogTypeId uint64
		var coverPlanId uint64
		var title string
		var abstract string
		var browseTotal uint
		var createdAt uint64
		var updatedAt uint64
		rows.Scan(&id, &userId, &blogTypeId, &coverPlanId, &title, &abstract, &browseTotal, &createdAt, &updatedAt)

		//博客实体
		blogEntity := new(blog.BlogListEntity)
		blogEntity.ID = id
		blogEntity.UserId = uint64(userId)
		blogEntity.BlogTypeId = blogTypeId
		blogEntity.CoverPlanId = coverPlanId
		blogEntity.Title = title
		blogEntity.Abstract = abstract
		blogEntity.BrowseTotal = browseTotal
		blogEntity.CreatedAt = createdAt
		blogEntity.UpdatedAt = updatedAt

		coverPlanIds = append(coverPlanIds, coverPlanId)
		blogTypeIds = append(blogTypeIds, blogTypeId)
		userIdIds = append(userIdIds, userId)

		queryResult = append(queryResult, blogEntity)
	}

	//填充信息
	PaddingAttachemtInfoByBlogListEntity(coverPlanIds, queryResult) //填充附件信息
	PaddingBlogTypeInfoByBlogListEntity(blogTypeIds, queryResult)   //博客类型实体
	PaddingUserInfoByBlogListEntity(userIdIds, queryResult)         //填充用户信息

	//构建结果返回
	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(perPage)
	result.SetList(queryResult)

	return
}

//排序方向
//per_page 每页多少条
func (s *BlogRtService) GetListBySort(sortDimension string, perPage int) (result *entity.ListResponseEntity) {

	db := mysql.GetDefaultDBConnect()
	var blogModelList []*model.BlogModel
	//查询
	switch sortDimension {
	case "browse_total":
		db = db.Table(model.BlogModel{}.TableName())
		db = db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
		db.Order("browse_total DESC").Limit(perPage).Find(&blogModelList)
	case "created_at":
		db = db.Table(model.BlogModel{}.TableName())
		db = db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
		db.Order("created_at DESC").Limit(perPage).Find(&blogModelList)
	default:
		exception.NewException(exception.VALIDATE_ERR, "非法的sort_dimension")
	}

	//转化为传输层的对象
	blogSortEntityList := ChangeBlogSortEntityByList(blogModelList)

	//构建结果返回
	result = new(entity.ListResponseEntity)
	filterList := []help.Filter{
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
	result.SetFilter(filterList)
	result.SetCount(int64(perPage))
	result.SetPerPage(perPage)
	result.SetList(blogSortEntityList)

	return result
}

//searchLevel 搜索等级
//keyword 搜索的关键字
//per_page 每页多少条
//page 第几页
func (s *BlogRtService) SearchBlog(searchLevel string, keyword string, perPage, page int) (result *entity.ListResponseEntity) {

	switch searchLevel {
	case MysqlSearchLevel:
		result = s.SearchBlogMysqlLevel(keyword, perPage, page)
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
func (s *BlogRtService) SearchBlogMysqlLevel(keyword string, perPage, page int) (result *entity.ListResponseEntity) {

	var blogModelList []*model.BlogModel
	//var count int64

	//如果存在缓存，先从缓冲中取
	lruCacheList, ok := GetBlogByLru(keyword)
	if ok {
		log.Println("lru中取数据")
		//构建结果返回
		result = new(entity.ListResponseEntity)
		result.SetCount(0)
		result.SetPerPage(perPage)
		result.SetList(lruCacheList.([]*blog.BlogListEntity))
		return result
	}

	//缓存没有，数据库取
	db := mysql.GetDefaultDBConnect()
	db = db.Table(model.BlogModel{}.TableName())

	db = db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	if keyword != "" {
		db = db.Where("content like ? OR title like ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	//db.Count(&count)
	db.Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&blogModelList)

	//转化为传输层的对象
	list := ChangeToBlogListEntityList(blogModelList)

	//构建结果返回
	result = new(entity.ListResponseEntity)

	result.SetCount(0)
	result.SetPerPage(perPage)
	result.SetList(list)

	if list != nil {

		//加入lru缓存
		AddBlogToLru(keyword, list)
	}

	log.Println("数据库去数据")

	return result
}

//blogInfo模块统计展示
func (s *BlogRtService) GetStat() (result *entity.ListResponseEntity) {

	response := make(map[string]uint, 3)

	db := mysql.GetDefaultDBConnect()

	blogTableName := model.BlogModel{}.TableName()
	db = db.Table(blogTableName)
	db = db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)

	//文章总数
	{
		var count int64
		//sql
		countDb := db.Select("id")
		countDb.Count(&count)
		response["blog_total"] = uint(count)
	}

	//最新一篇文章的时间、最新和最早的文章时间戳
	{
		//最新的博客时间
		var lastCreateAt uint
		lastCreateAtRow := db.Select("created_at").Order("id DESC").Limit(1).Row()
		lastCreateAtRow.Scan(&lastCreateAt)
		response["last_create_at"] = lastCreateAt

		//最老的博客时间
		var firstCreateAt uint
		firstCreateAtDb := mysql.GetDefaultDBConnect()
		firstCreateAtRow := firstCreateAtDb.Table(blogTableName).Select("created_at").
			Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1).
			Order("id ASC").Limit(1).Row()
		firstCreateAtRow.Scan(&firstCreateAt)
		response["diff_time"] = lastCreateAt - firstCreateAt
	}

	//构建结果返回
	result = new(entity.ListResponseEntity)

	result.SetList(response)

	return result
}
