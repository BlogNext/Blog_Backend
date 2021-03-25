package blog

import (
	"encoding/json"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"strings"
	"time"
)

//博客类型前台服务
type BlogTypeRtService struct {
}

//通过博客ids获取blogTypeEntity
func (s *BlogTypeRtService) getListByids(ids []uint64) (result map[uint64]*blog.BlogTypeEntity) {
	db := mysql.GetDefaultDBConnect()
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
	db := mysql.GetDefaultDBConnect()

	var count int64
	db = db.Model(&model.BlogTypeModel{})

	db.DryRun = true
	statement := db.Count(&count).Statement
	db.DryRun = false
	cacheKey := "blog_type_list" + db.Dialector.Explain(statement.SQL.String(),statement.Vars...)
	//如果存在缓存，先从缓冲中取
	lruCacheList, ok := BlgLruUnsafety.Get(cacheKey)
	if ok {
		//有缓存
		//构建结果返回
		result = new(entity.ListResponseEntity)
		json.Unmarshal(lruCacheList.([]uint8),result)
		return result
	}

	db.Count(&count)

	rows, _ := db.Select("id, yuque_name, yuque_id, created_at, updated_at").
		Limit(perPage).Offset((page - 1) * perPage).Rows()

	defer rows.Close()

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

	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(perPage)
	result.SetList(blogTypeModelList)



	//加入缓存
	jsonCache,_ := json.Marshal(result)
	BlgLruUnsafety.Add(cacheKey, jsonCache, 5*time.Minute)

	return
}
