package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/common"
	"github.com/thoas/go-funk"
	"log"
	"strings"
	"time"
)

//博客
type BlogService struct {
}

//添加博客
func (s *BlogService) AddBlog(blog_type_id, cover_plan_id int64, title, abstract, content string) {
	db := mysql.GetDefaultDBConnect()

	blog_model := new(model.BlogModel)
	blog_model.BlogTypeId = blog_type_id
	blog_model.Title = title
	blog_model.Abstract = abstract
	blog_model.Content = content
	blog_model.CoverPlanId = cover_plan_id
	blog_model.CreateTime = time.Now().Unix()
	blog_model.UpdateTime = time.Now().Unix()

	db.Create(blog_model)

	if db.NewRecord(*blog_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("保存失败:%s", db.Error.Error())))
	}
}

//更新博客
func (s *BlogService) UpdateBlog(id, blog_type_id, cover_plan_id int64, title, abstract, content string) {
	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	db.Where("id = ?", id).First(blog_model)

	if db.NewRecord(*blog_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("找不到记录:%d", id)))
	}

	blog_model.BlogTypeId = blog_type_id
	blog_model.Title = title
	blog_model.Abstract = abstract
	blog_model.Content = content
	blog_model.UpdateTime = time.Now().Unix()
	blog_model.CoverPlanId = cover_plan_id

	db.Save(blog_model)

	if db.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db.Error.Error())))
	}
}

//列表页
func (s *BlogService) GetList() (result []map[string]interface{}) {
	db := mysql.GetDefaultDBConnect()

	//连表表名
	blog_type_table_name := model.BlogTypeModel{}.TableName()
	blog_table_name := model.BlogModel{}.TableName()

	//博客需要的字段
	blog_felid := []string{"id", "blog_type_id", "cover_plan_id", "title", "create_time", "update_time"}

	for index, felid := range blog_felid {
		blog_felid[index] = fmt.Sprintf("%s.%s", blog_table_name, felid)
	}

	//博客类型需要的字段
	blog_type_felid := []string{"title as blog_type_title"}

	for index, felid := range blog_type_felid {
		blog_type_felid[index] = fmt.Sprintf("%s.%s", blog_type_table_name, felid)
	}

	select_felid := append(blog_felid, blog_type_felid...)

	//sql
	rows, _ := db.Table(blog_table_name).
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.blog_type_id = %s.id", blog_type_table_name, blog_table_name, blog_type_table_name)).
		Select(strings.Join(select_felid, ", ")).Rows()

	result = make([]map[string]interface{}, 0)

	cover_plan_ids := make([]int64, 0)

	for rows.Next() {
		var id int64
		var blog_type_id int64
		var cover_plan_id int64
		var title string
		var create_time int64
		var update_time int64
		var blog_type_title string
		rows.Scan(&id, &blog_type_id, &cover_plan_id, &title, &create_time, &update_time, &blog_type_title)

		item := make(map[string]interface{})
		item["ID"] = id
		item["BlogTypeId"] = blog_type_id
		item["CoverPlanId"] = cover_plan_id
		item["Title"] = title
		item["CreateTime"] = create_time
		item["UpdateTime"] = update_time
		blog_type_info := make(map[string]interface{})
		blog_type_info["Title"] = blog_type_title
		item["BlogTypeInfo"] = blog_type_info

		cover_plan_ids = append(cover_plan_ids, cover_plan_id)

		result = append(result, item)
	}

	//填充信息
	s.paddingListAttachemtInfo(cover_plan_ids, result)

	return
}

//填充附件信息
func (s *BlogService) paddingListAttachemtInfo(cover_plan_ids []int64, result []map[string]interface{}) {

	//获取图片的ids,填充图片信息
	attachment_list := common.GetAttachmentImages(cover_plan_ids)

	if attachment_list != nil {
		//转化成map
		log.Println(attachment_list)
		attachment_list_map := funk.ToMap(attachment_list, "ID").(map[uint]model.FullAttachmentExtend)
		log.Println(attachment_list_map)
		log.Println("哗啦啦")
		//填充图片信息
		for _, item := range result {
			log.Println(attachment_list_map[uint(item["CoverPlanId"].(int64))])

			if attachment_item, ok := attachment_list_map[uint(item["CoverPlanId"].(int64))]; ok {

				item["AttachmentInfo"] = attachment_item
			} else {
				item["AttachmentInfo"] = nil
			}

		}
	}
}
