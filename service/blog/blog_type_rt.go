package blog

import (
	"github.com/blog_backend/common-lib/arithmetic/lru"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/common-lib/db/mysql/my_db_proxy"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
	"strings"
	"time"
)

//博客类型前台服务
type BlogTypeRtService struct {
}

//通过博客ids获取blogTypeEntity
func (s *BlogTypeRtService) getListByids(ids []uint64) (result map[uint64]*blog.BlogTypeEntity) {
	db := mysql.GetNewDB(false)
	tableName := model.BlogTypeModel{}.TableName()
	selectFelid := []string{"id", "yuque_name", "created_at", "updated_at"}
	rows, err := db.Table(tableName).
		Select(strings.Join(selectFelid, ", ")).Where("id IN (?)", ids).Rows()

	if err != nil {
		return nil
	}

	result = make(map[uint64]*blog.BlogTypeEntity)

	for rows.Next() {
		var id uint64
		var yuqueName string
		var createdAt uint64
		var updatedAt uint64

		rows.Scan(&id, &yuqueName, &createdAt, &updatedAt)

		blogTypeEntity := new(blog.BlogTypeEntity)
		blogTypeEntity.ID = id
		blogTypeEntity.Title = yuqueName
		blogTypeEntity.CreatedAt = createdAt
		blogTypeEntity.UpdatedAt = updatedAt

		result[id] = blogTypeEntity

	}

	return
}

//获取分类列表
func (s *BlogTypeRtService) GetList(perPage, page int) (result *entity.ListResponseEntity) {

	cacheKey := lru.GenerateCacheKey(perPage, page)

	resultCache, ok := BlgLruUnsafety.Get(cacheKey)

	if ok {
		//有缓存,直接返回缓存对象
		result = resultCache.(*entity.ListResponseEntity)
		return result
	}

	myDBProxy := my_db_proxy.NewMyDBProxyByTable(model.BlogTypeModel.TableName())

	//返回结果
	result = new(entity.ListResponseEntity)

	//没有缓存，获取数据集
	myDBProxy.ExecProxy(func(db *gorm.DB) {

		rows, _ := db.Select("id, yuque_name, yuque_id, created_at, updated_at").
			Limit(perPage).Offset((page - 1) * perPage).Rows()

		blogTypeModelList := make([]map[string]interface{}, 0)

		for rows.Next() {

			var id int64
			var yuqueName string
			var yuqueId int64
			var createdAt int64
			var updatedAt int64
			rows.Scan(&id, &yuqueName, &yuqueId, &createdAt, &updatedAt)

			item := make(map[string]interface{})

			item["id"] = id
			item["yuque_name"] = yuqueName
			item["yuque_id"] = yuqueId
			item["create_time"] = createdAt
			item["update_time"] = updatedAt

			blogTypeModelList = append(blogTypeModelList, item)
		}

		result.SetList(blogTypeModelList)
	})

	//没有缓存的情况下，继续计算count值，然后设置count
	myDBProxy.ExecProxy(func(db *gorm.DB) {
		var count int64
		db.Count(&count)
		result.SetCount(count)
		result.SetPerPage(perPage)
	})

	BlgLruUnsafety.Add(cacheKey, result, 5*time.Minute)

	return
}
