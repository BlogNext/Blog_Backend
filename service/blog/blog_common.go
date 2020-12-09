package blog

import (
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	attachment_service "github.com/blog_backend/service/attachment"
	"github.com/blog_backend/service/user"
	"github.com/thoas/go-funk"
)

//填充附件信息
func PaddingAttachemtInfo(cover_plan_ids []uint64, result []*blog.BlogEntity) {

	//获取图片的ids,填充图片信息
	attachment_list := attachment_service.GetAttachmentImages(cover_plan_ids)

	if attachment_list != nil {
		//转化成map

		attachment_list_map := funk.ToMap(attachment_list, "ID").(map[uint64]*attachment.AttachmentEntity)

		//填充图片信息
		for _, item := range result {

			if attachment_item, ok := attachment_list_map[item.CoverPlanId]; ok {
				item.CoverPlanInfo = attachment_item
			} else {
				item.CoverPlanInfo = nil
			}

		}
	}
}

//填充附件信息，blogEntity
func PaddingAttachemtInfoByBlogSortEntityList(cover_plan_ids []uint64, result []*blog.BlogSortEntity) {
	//获取图片的ids,填充图片信息
	attachment_map_list := attachment_service.GetAttachmentImagesMap(cover_plan_ids)

	if attachment_map_list != nil {
		//填充图片信息
		for _, item := range result {

			if attachment_item, ok := attachment_map_list[uint(item.CoverPlanId)]; ok {
				item.CoverPlanInfo = attachment_item
			} else {
				item.CoverPlanInfo = nil
			}

		}
	}
}

//填充附件信息，blogList实体
func PaddingAttachemtInfoByBlogListEntity(cover_plan_ids []uint64, result []*blog.BlogListEntity) {
	//获取图片的ids,填充图片信息
	attachment_list := attachment_service.GetAttachmentImages(cover_plan_ids)

	if attachment_list != nil {
		//转化成map

		attachment_list_map := funk.ToMap(attachment_list, "ID").(map[uint64]*attachment.AttachmentEntity)

		//填充图片信息
		for _, item := range result {

			if attachment_item, ok := attachment_list_map[item.CoverPlanId]; ok {
				item.CoverPlanInfo = attachment_item
			} else {
				item.CoverPlanInfo = nil
			}

		}
	}
}

//填充用户信息
func PaddingUserInfo(user_ids []uint, result []*blog.BlogEntity) {

	user_entity_map := user.GetUserEntityByUserIds(user_ids)
	for _, item := range result {
		if user_info, ok := user_entity_map[uint(item.UserId)]; ok {
			item.UserInfo = user_info
		} else {
			item.UserInfo = nil
		}

	}
}

//填充用户信息
func PaddingUserInfoByBlogListEntity(user_ids []uint, result []*blog.BlogListEntity) {
	user_entity_map := user.GetUserEntityByUserIds(user_ids)

	for _, item := range result {

		if user_info, ok := user_entity_map[uint(item.UserId)]; ok {
			item.UserInfo = user_info
		} else {
			item.CoverPlanInfo = nil
		}

	}
}

//填充博客类型信息
func PaddingBlogTypeInfo(blog_type_ids []uint64, result []*blog.BlogEntity) {

	blog_type_service := new(BlogTypeRtService)
	blog_type_list := blog_type_service.getListByids(blog_type_ids)

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

//填充博客类型信息
func PaddingBlogTypeInfoByBlogListEntity(blog_type_ids []uint64, result []*blog.BlogListEntity) {
	blog_type_service := new(BlogTypeRtService)
	blog_type_list := blog_type_service.getListByids(blog_type_ids)

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
	blog_entity.UserId = uint64(blog_model.UserID)
	blog_entity.CreatedAt = uint64(blog_model.CreatedAt)
	blog_entity.UpdatedAt = uint64(blog_model.UpdatedAt)
	blog_entity.BlogTypeId = uint64(blog_model.BlogTypeId)
	blog_entity.CoverPlanId = uint64(blog_model.CoverPlanId)
	blog_entity.BrowseTotal = blog_model.BrowseTotal
	blog_entity.YuqueFormat = blog_model.YuqueFormat
	blog_entity.Title = blog_model.Title
	blog_entity.Abstract = blog_model.Abstract
	blog_entity.Content = blog_model.Content
	blog_entity.DocID = blog_model.DocID

	blog_entity_list := []*blog.BlogEntity{blog_entity}
	//填充别的实体信息

	PaddingAttachemtInfo([]uint64{blog_entity.CoverPlanId}, blog_entity_list)

	PaddingBlogTypeInfo([]uint64{blog_entity.BlogTypeId}, blog_entity_list)

	PaddingUserInfo([]uint{uint(blog_entity.UserId)}, blog_entity_list)

	return blog_entity_list[0]
}

//模型转化成BlogListEntity实体
func ChangeToBlogListEntity(blog_model *model.BlogModel) *blog.BlogListEntity {

	blog_entity := new(blog.BlogListEntity)
	blog_entity.ID = uint64(blog_model.ID)
	blog_entity.CreatedAt = uint64(blog_model.CreatedAt)
	blog_entity.UpdatedAt = uint64(blog_model.UpdatedAt)
	blog_entity.BlogTypeId = uint64(blog_model.BlogTypeId)
	blog_entity.CoverPlanId = uint64(blog_model.CoverPlanId)
	blog_entity.Title = blog_model.Title
	blog_entity.Abstract = blog_model.Abstract
	blog_entity.BrowseTotal = blog_model.BrowseTotal
	blog_entity.DocID = blog_model.DocID
	blog_entity.UserId = uint64(blog_model.UserID)

	blog_entity_list := []*blog.BlogListEntity{blog_entity}
	//填充别的实体信息

	PaddingAttachemtInfoByBlogListEntity([]uint64{blog_entity.CoverPlanId}, blog_entity_list)

	PaddingBlogTypeInfoByBlogListEntity([]uint64{blog_entity.BlogTypeId}, blog_entity_list)

	PaddingUserInfoByBlogListEntity([]uint{uint(blog_entity.UserId)}, blog_entity_list) //填充用户信息

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

	user_id_ids := make([]uint, number)

	for index, item := range blog_model_list {
		blog_entity := new(blog.BlogEntity)
		blog_entity.ID = uint64(item.ID)
		blog_entity.UserId = uint64(item.UserID)
		blog_entity.CreatedAt = uint64(item.CreatedAt)
		blog_entity.UpdatedAt = uint64(item.UpdatedAt)
		blog_entity.BlogTypeId = uint64(item.BlogTypeId)
		blog_entity.CoverPlanId = uint64(item.CoverPlanId)
		blog_entity.Title = item.Title
		blog_entity.YuqueFormat = item.YuqueFormat
		blog_entity.Abstract = item.Abstract
		blog_entity.BrowseTotal = item.BrowseTotal
		blog_entity.Content = item.Content
		blog_entity.DocID = item.DocID

		//填充数据
		cover_plan_ids[index] = blog_entity.CoverPlanId
		blog_type_ids[index] = blog_entity.BlogTypeId
		user_id_ids[index] = item.UserID

		//返回的集合
		blog_entity_list[index] = blog_entity
	}

	PaddingAttachemtInfo(cover_plan_ids, blog_entity_list)

	PaddingBlogTypeInfo(blog_type_ids, blog_entity_list)

	PaddingUserInfo(user_id_ids, blog_entity_list) //填充用户信息

	return blog_entity_list

}

//模型转化为blogSortEntity实体
func ChangeBlogSortEntityByList(blog_model_list []*model.BlogModel) []*blog.BlogSortEntity {
	number := len(blog_model_list)

	if number <= 0 {
		return nil
	}

	blog_sort_entity_list := make([]*blog.BlogSortEntity, number)

	cover_plan_ids := make([]uint64, number)

	for index, item := range blog_model_list {
		blog_sort_entity := new(blog.BlogSortEntity)
		blog_sort_entity.ID = uint64(item.ID)
		blog_sort_entity.CoverPlanId = uint64(item.CoverPlanId)
		blog_sort_entity.Title = item.Title
		blog_sort_entity.BrowseTotal = item.BrowseTotal

		//填充数据
		cover_plan_ids[index] = blog_sort_entity.CoverPlanId

		//返回的集合
		blog_sort_entity_list[index] = blog_sort_entity
	}

	PaddingAttachemtInfoByBlogSortEntityList(cover_plan_ids, blog_sort_entity_list)

	return blog_sort_entity_list
}

//模型转化为BlogListEntity实体List操作
func ChangeToBlogListEntityList(blog_model_list []*model.BlogModel) []*blog.BlogListEntity {
	number := len(blog_model_list)

	if number <= 0 {
		return nil
	}

	blog_entity_list := make([]*blog.BlogListEntity, number)

	cover_plan_ids := make([]uint64, number)

	blog_type_ids := make([]uint64, number)

	user_id_ids := make([]uint, number)

	for index, item := range blog_model_list {
		blog_entity := new(blog.BlogListEntity)
		blog_entity.ID = uint64(item.ID)
		blog_entity.UserId = uint64(item.UserID)
		blog_entity.CreatedAt = uint64(item.CreatedAt)
		blog_entity.UpdatedAt = uint64(item.UpdatedAt)
		blog_entity.BlogTypeId = uint64(item.BlogTypeId)
		blog_entity.CoverPlanId = uint64(item.CoverPlanId)
		blog_entity.Title = item.Title
		blog_entity.Abstract = item.Abstract
		blog_entity.BrowseTotal = item.BrowseTotal
		blog_entity.DocID = item.DocID

		//填充数据
		cover_plan_ids[index] = blog_entity.CoverPlanId
		blog_type_ids[index] = blog_entity.BlogTypeId
		user_id_ids[index] = item.UserID

		//返回的集合
		blog_entity_list[index] = blog_entity
	}

	PaddingAttachemtInfoByBlogListEntity(cover_plan_ids, blog_entity_list)

	PaddingBlogTypeInfoByBlogListEntity(blog_type_ids, blog_entity_list)

	PaddingUserInfoByBlogListEntity(user_id_ids, blog_entity_list) //填充用户信息

	return blog_entity_list

}
