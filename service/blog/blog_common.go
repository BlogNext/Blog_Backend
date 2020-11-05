package blog

import (
	"fmt"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	attachment_service "github.com/blog_backend/service/attachment"
	"github.com/thoas/go-funk"
	"log"
)

//填充附件信息
func PaddingAttachemtInfo(cover_plan_ids []uint64, result []*blog.BlogEntity) {

	//获取图片的ids,填充图片信息
	attachment_list := attachment_service.GetAttachmentImages(cover_plan_ids)

	if attachment_list != nil {
		//转化成map
		log.Println(attachment_list)
		log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", attachment_list, attachment_list, attachment_list))

		attachment_list_map := funk.ToMap(attachment_list, "ID").(map[uint64]*attachment.AttachmentEntity)
		log.Println(attachment_list_map)
		log.Println("哗啦啦")
		//填充图片信息
		for _, item := range result {
			log.Println(attachment_list_map[item.CoverPlanId])

			if attachment_item, ok := attachment_list_map[item.CoverPlanId]; ok {
				item.AttachmentInfo = attachment_item
			} else {
				item.AttachmentInfo = nil
			}

		}
	}
}

//填充博客类型信息
func PaddingBlogTypeInfo(blog_type_ids []uint64, result []*blog.BlogEntity) {

	blog_type_service := new(BlogTypeRtService)
	blog_type_list := blog_type_service.getListByids(blog_type_ids)

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_type_list, blog_type_list, blog_type_list))

	if blog_type_list != nil {
		//填充博客类型

		for _, item := range result {
			if blog_type_entity, ok := blog_type_list[item.BlogTypeId]; ok {
				item.BlogTypeObject = blog_type_entity
			} else {
				item.BlogTypeObject = nil
			}
		}
	}
}

//模型转化成BlogEntity实体
func ChangeToBlogEntity(blog_model *model.BlogModel) *blog.BlogEntity {

	blog_entity := new(blog.BlogEntity)
	blog_entity.ID = uint64(blog_model.ID)
	blog_entity.CreatedAt = uint64(blog_model.CreatedAt)
	blog_entity.UpdatedAt = uint64(blog_model.UpdatedAt)
	blog_entity.BlogTypeId = uint64(blog_model.BlogTypeId)
	blog_entity.CoverPlanId = uint64(blog_model.CoverPlanId)
	blog_entity.Title = blog_model.Title
	blog_entity.Abstract = blog_model.Abstract
	blog_entity.Content = blog_model.Content
	blog_entity.DocID = blog_model.DocID

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity, blog_entity, blog_entity))

	blog_entity_list := []*blog.BlogEntity{blog_entity}
	//填充别的实体信息

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity_list, blog_entity_list, blog_entity_list))

	PaddingAttachemtInfo([]uint64{blog_entity.CoverPlanId}, blog_entity_list)

	log.Println("附件填充完毕", fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity_list, blog_entity_list, blog_entity_list))

	PaddingBlogTypeInfo([]uint64{blog_entity.BlogTypeId}, blog_entity_list)

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity_list, blog_entity_list, blog_entity_list))

	return blog_entity_list[0]
}

//模型转化为BlogEntity实体List操作
func ChangeToBlogEntityList(blog_model_list []*model.BlogModel) []*blog.BlogEntity {
	number := len(blog_model_list)

	if number <= 0 {
		return nil
	}

	blog_entity_list := make([]*blog.BlogEntity, number)

	cover_plan_ids := make([]uint64, number)

	blog_type_ids := make([]uint64, number)

	for index, item := range blog_model_list {
		blog_entity := new(blog.BlogEntity)
		blog_entity.ID = uint64(item.ID)
		blog_entity.CreatedAt = uint64(item.CreatedAt)
		blog_entity.UpdatedAt = uint64(item.UpdatedAt)
		blog_entity.BlogTypeId = uint64(item.BlogTypeId)
		blog_entity.CoverPlanId = uint64(item.CoverPlanId)
		blog_entity.Title = item.Title
		blog_entity.Abstract = item.Abstract
		blog_entity.Content = item.Content
		blog_entity.DocID = item.DocID

		cover_plan_ids[index] = blog_entity.CoverPlanId
		blog_type_ids[index] = blog_entity.BlogTypeId

		blog_entity_list[index] = blog_entity
	}

	PaddingAttachemtInfo(cover_plan_ids, blog_entity_list)

	PaddingBlogTypeInfo(blog_type_ids, blog_entity_list)

	return blog_entity_list

}
