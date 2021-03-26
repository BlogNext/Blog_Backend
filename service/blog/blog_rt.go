package blog

import (
	"encoding/json"
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
	"time"
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
func (s *BlogRtService) IncBrowse(id uint64) {
	db := mysql.GetNewDB(false)
	db.Model(model.BlogModel{}).Where("id = ?", id).UpdateColumn("browse_total", gorm.Expr("browse_total + ?", 1))
}

//博客详情
func (s *BlogRtService) Detail(id uint64) *blog.BlogEntity {
	db := mysql.GetNewDB(false)
	blogModel := new(model.BlogModel)
	db = db.Where("id = ?", id)
	db.First(blogModel, id)

	find := errors.Is(db.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到博客:%d", id))
	}

	result := make([]*blog.BlogEntity, 1)
	result[0] = ChangeToBlogEntity(blogModel)

	return result[0]
}

//获取博客列表，用于私人空间
func (s *BlogRtService) GetListByPerson(perPage, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetNewDB(false)

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
	userIdIds := make([]uint64, 0)

	for rows.Next() {
		var id uint64
		var userId uint64
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
		blogEntity.UserId = userId
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
	db := mysql.GetNewDB(false)
	//获取执行的sql，不执行sql
	dbDryRun := mysql.GetNewDB(true)

	blogTableName := model.BlogModel{}.TableName()

	//博客需要的字段
	blogFelid := []string{"id", "user_id", "blog_type_id", "cover_plan_id", "title", "abstract", "browse_total", "created_at", "updated_at"}

	for index, felid := range blogFelid {
		blogFelid[index] = fmt.Sprintf("%s.%s", blogTableName, felid)
	}

	var count int64
	//sql
	db.Table(blogTableName)
	dbDryRun.Table(blogTableName)

	db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	dbDryRun.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	//过滤分类id过滤
	if filter["blog_type_id"] != "" {
		db.Where("blog_type_id = ?", filter["blog_type_id"])
		dbDryRun.Where("blog_type_id = ?", filter["blog_type_id"])
	}

	//缓存
	statement := dbDryRun.Count(&count).Statement
	cacheKey := dbDryRun.Dialector.Explain(statement.SQL.String(), statement.Vars...)
	cacheKey = "list_" + cacheKey

	//如果存在缓存，先从缓冲中取
	lruCacheList, ok := BlgLruUnsafety.Get(cacheKey)
	if ok {
		//有缓存
		//构建结果返回
		result = new(entity.ListResponseEntity)
		json.Unmarshal(lruCacheList.([]uint8), result)
		return result
	}

	//真正执行sql
	db.Count(&count)

	rows, err := db.Select(strings.Join(blogFelid, ", ")).Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Rows()

	if err != nil {
		return nil
	}

	queryResult := make([]*blog.BlogListEntity, 0)

	coverPlanIds := make([]uint64, 0)
	blogTypeIds := make([]uint64, 0)
	userIdIds := make([]uint64, 0)

	for rows.Next() {
		var id uint64
		var userId uint64
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

	//加入缓存
	jsonCache, _ := json.Marshal(result)
	BlgLruUnsafety.Add(cacheKey, jsonCache, 5*time.Second)

	return
}

//排序方向
//per_page 每页多少条
func (s *BlogRtService) GetListBySort(sortDimension string, perPage int) (result *entity.ListResponseEntity) {

	db := mysql.GetNewDB(false)
	//获取执行的sql，不执行sql
	dbDryRun := mysql.GetNewDB(true)

	var blogModelList []*model.BlogModel
	var cacheKey string
	//查询
	switch sortDimension {
	case "browse_total":
	case "created_at":
	default:
		exception.NewException(exception.VALIDATE_ERR, "非法的sort_dimension")
	}

	db.Table(model.BlogModel{}.TableName())
	dbDryRun.Table(model.BlogModel{}.TableName())
	db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	dbDryRun.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)

	//排序字段
	orderBySql := fmt.Sprintf("%s DESC", sortDimension)
	statement := dbDryRun.Order(orderBySql).Limit(perPage).Find(&blogModelList).Statement
	cacheKey = db.Dialector.Explain(statement.SQL.String(), statement.Vars...)
	cacheKey = "sort_" + cacheKey

	//如果存在缓存，先从缓冲中取
	lruCacheList, ok := BlgLruUnsafety.Get(cacheKey)
	if ok {
		//有缓存
		//构建结果返回
		result = new(entity.ListResponseEntity)
		json.Unmarshal(lruCacheList.([]uint8), result)
		return result
	}

	//获取数据
	db.Order(orderBySql).Limit(perPage).Find(&blogModelList)

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

	//加入缓存
	jsonCache, _ := json.Marshal(result)
	BlgLruUnsafety.Add(cacheKey, jsonCache, 5*time.Second)

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

	//缓存没有，数据库取
	db := mysql.GetNewDB(false)
	dbDryRun := mysql.GetNewDB(true)

	db.Table(model.BlogModel{}.TableName())
	dbDryRun.Table(model.BlogModel{}.TableName())
	db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	dbDryRun.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)
	if keyword != "" {
		db.Where("content like ? OR title like ?", "%"+keyword+"%", "%"+keyword+"%")
		dbDryRun.Where("content like ? OR title like ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var blogModelList []*model.BlogModel

	//先看看有没有缓存
	statement := dbDryRun.Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&blogModelList).Statement
	cacheKey := dbDryRun.Dialector.Explain(statement.SQL.String(), statement.Vars...)
	cacheKey = "search_" + cacheKey
	//如果存在缓存，先从缓冲中取
	lruCacheList, ok := BlgLruUnsafety.Get(cacheKey)
	if ok {
		//有缓存
		//构建结果返回
		result = new(entity.ListResponseEntity)
		json.Unmarshal(lruCacheList.([]uint8), result)
		return result
	}

	//没有缓存
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
		jsonCache, _ := json.Marshal(result)
		BlgLruUnsafety.Add(cacheKey, jsonCache, 5*time.Second)
	}

	return result
}

//blogInfo模块统计展示
func (s *BlogRtService) GetStat() (result *entity.ListResponseEntity) {

	response := make(map[string]uint, 3)

	db := mysql.GetNewDB(false)

	cacheKey := "stat_cache"
	lruCacheList, ok := BlgLruUnsafety.Get(cacheKey)
	if ok {
		result = new(entity.ListResponseEntity)
		json.Unmarshal(lruCacheList.([]uint8), result)
		return result
	}

	blogTableName := model.BlogModel{}.TableName()

	db.Table(blogTableName)
	db.Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1)

	//文章总数
	{
		var count int64
		//sql
		db.Select("id")
		db.Count(&count)
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
		firstCreateAtDb := mysql.GetNewDB(false)
		firstCreateAtRow := firstCreateAtDb.Table(blogTableName).Select("created_at").
			Where("yuque_public = ?", model.BLOG_MODEL_YUQUE_PUBLIC_1).
			Order("id ASC").Limit(1).Row()
		firstCreateAtRow.Scan(&firstCreateAt)
		response["diff_time"] = lastCreateAt - firstCreateAt
	}

	//构建结果返回
	result = new(entity.ListResponseEntity)

	result.SetList(response)

	jsonCache, _ := json.Marshal(result)
	BlgLruUnsafety.Add(cacheKey, jsonCache, 5*time.Second)

	return result
}
