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
func PaddingAttachemtInfo(coverPlanIds []uint64, result []*blog.BlogEntity) {

	//获取图片的ids,填充图片信息
	attachmentList := attachment_service.GetAttachmentImages(coverPlanIds)

	if attachmentList != nil {
		//转化成map

		attachmentListMap := funk.ToMap(attachmentList, "ID").(map[uint64]*attachment.AttachmentEntity)

		//填充图片信息
		for _, item := range result {

			if attachmentItem, ok := attachmentListMap[item.CoverPlanId]; ok {
				item.CoverPlanInfo = attachmentItem
			} else {
				item.CoverPlanInfo = nil
			}

		}
	}
}

//填充附件信息，blogEntity
func PaddingAttachemtInfoByBlogSortEntityList(coverPlanIds []uint64, result []*blog.BlogSortEntity) {
	//获取图片的ids,填充图片信息
	attachmentMapList := attachment_service.GetAttachmentImagesMap(coverPlanIds)

	if attachmentMapList != nil {
		//填充图片信息
		for _, item := range result {

			if attachmentItem, ok := attachmentMapList[uint(item.CoverPlanId)]; ok {
				item.CoverPlanInfo = attachmentItem
			} else {
				item.CoverPlanInfo = nil
			}

		}
	}
}

//填充附件信息，blogList实体
func PaddingAttachemtInfoByBlogListEntity(coverPlanIds []uint64, result []*blog.BlogListEntity) {
	//获取图片的ids,填充图片信息
	attachmentList := attachment_service.GetAttachmentImages(coverPlanIds)

	if attachmentList != nil {
		//转化成map

		attachmentListMap := funk.ToMap(attachmentList, "ID").(map[uint64]*attachment.AttachmentEntity)

		//填充图片信息
		for _, item := range result {

			if attachmentItem, ok := attachmentListMap[item.CoverPlanId]; ok {
				item.CoverPlanInfo = attachmentItem
			} else {
				item.CoverPlanInfo = nil
			}

		}
	}
}

//填充用户信息
func PaddingUserInfo(userIds []uint64, result []*blog.BlogEntity) {

	userEntityMap := user.GetUserEntityByUserIds(userIds)
	for _, item := range result {
		if userInfo, ok := userEntityMap[item.UserId]; ok {
			item.UserInfo = userInfo
		} else {
			item.UserInfo = nil
		}

	}
}

//填充用户信息
func PaddingUserInfoByBlogListEntity(userIds []uint64, result []*blog.BlogListEntity) {
	userEntityMap := user.GetUserEntityByUserIds(userIds)

	for _, item := range result {

		if userInfo, ok := userEntityMap[item.UserId]; ok {
			item.UserInfo = userInfo
		} else {
			item.CoverPlanInfo = nil
		}

	}
}

//填充博客类型信息
func PaddingBlogTypeInfo(blogTypeIds []uint64, result []*blog.BlogEntity) {

	blogTypeService := new(BlogTypeRtService)
	blogTypeList := blogTypeService.getListByids(blogTypeIds)

	if blogTypeList != nil {
		//填充博客类型

		for _, item := range result {
			if blogTypeEntity, ok := blogTypeList[item.BlogTypeId]; ok {
				item.BlogTypeObject = blogTypeEntity
			} else {
				item.BlogTypeObject = nil
			}
		}
	}
}

//填充博客类型信息
func PaddingBlogTypeInfoByBlogListEntity(blogTypeIds []uint64, result []*blog.BlogListEntity) {
	blogTypeService := new(BlogTypeRtService)
	blogTypeList := blogTypeService.getListByids(blogTypeIds)

	if blogTypeList != nil {
		//填充博客类型

		for _, item := range result {
			if blogTypeEntity, ok := blogTypeList[item.BlogTypeId]; ok {
				item.BlogTypeObject = blogTypeEntity
			} else {
				item.BlogTypeObject = nil
			}
		}
	}
}

//模型转化成BlogEntity实体
func ChangeToBlogEntity(blogModel *model.BlogModel) *blog.BlogEntity {

	blogEntity := new(blog.BlogEntity)
	blogEntity.ID = blogModel.ID
	blogEntity.UserId = blogModel.UserID
	blogEntity.CreatedAt = uint64(blogModel.CreatedAt)
	blogEntity.UpdatedAt = uint64(blogModel.UpdatedAt)
	blogEntity.BlogTypeId = blogModel.BlogTypeId
	blogEntity.CoverPlanId = blogModel.CoverPlanId
	blogEntity.BrowseTotal = blogModel.BrowseTotal
	blogEntity.YuqueFormat = blogModel.YuqueFormat
	blogEntity.Title = blogModel.Title
	blogEntity.Abstract = blogModel.Abstract
	blogEntity.Content = blogModel.Content
	blogEntity.DocID = blogModel.DocID

	blogEntityList := []*blog.BlogEntity{blogEntity}
	//填充别的实体信息

	PaddingAttachemtInfo([]uint64{blogEntity.CoverPlanId}, blogEntityList)

	PaddingBlogTypeInfo([]uint64{blogEntity.BlogTypeId}, blogEntityList)

	PaddingUserInfo([]uint64{blogEntity.UserId}, blogEntityList)

	return blogEntityList[0]
}

//模型转化成BlogListEntity实体
func ChangeToBlogListEntity(blogModel *model.BlogModel) *blog.BlogListEntity {

	blogEntity := new(blog.BlogListEntity)
	blogEntity.ID = blogModel.ID
	blogEntity.CreatedAt = uint64(blogModel.CreatedAt)
	blogEntity.UpdatedAt = uint64(blogModel.UpdatedAt)
	blogEntity.BlogTypeId = blogModel.BlogTypeId
	blogEntity.CoverPlanId = blogModel.CoverPlanId
	blogEntity.Title = blogModel.Title
	blogEntity.Abstract = blogModel.Abstract
	blogEntity.BrowseTotal = blogModel.BrowseTotal
	blogEntity.DocID = blogModel.DocID
	blogEntity.UserId = blogModel.UserID

	blogEntityList := []*blog.BlogListEntity{blogEntity}
	//填充别的实体信息

	PaddingAttachemtInfoByBlogListEntity([]uint64{blogEntity.CoverPlanId}, blogEntityList)

	PaddingBlogTypeInfoByBlogListEntity([]uint64{blogEntity.BlogTypeId}, blogEntityList)

	PaddingUserInfoByBlogListEntity([]uint64{blogEntity.UserId}, blogEntityList) //填充用户信息

	return blogEntityList[0]
}

//模型转化为BlogEntity实体List操作
func ChangeToBlogEntityList(blogModelList []*model.BlogModel) []*blog.BlogEntity {
	number := len(blogModelList)

	if number <= 0 {
		return nil
	}

	blogEntityList := make([]*blog.BlogEntity, number)

	coverPlanIds := make([]uint64, number)

	blogTypeIds := make([]uint64, number)

	userIdIds := make([]uint64, number)

	for index, item := range blogModelList {
		blogEntity := new(blog.BlogEntity)
		blogEntity.ID = item.ID
		blogEntity.UserId = item.UserID
		blogEntity.CreatedAt = uint64(item.CreatedAt)
		blogEntity.UpdatedAt = uint64(item.UpdatedAt)
		blogEntity.BlogTypeId = item.BlogTypeId
		blogEntity.CoverPlanId = item.CoverPlanId
		blogEntity.Title = item.Title
		blogEntity.YuqueFormat = item.YuqueFormat
		blogEntity.Abstract = item.Abstract
		blogEntity.BrowseTotal = item.BrowseTotal
		blogEntity.Content = item.Content
		blogEntity.DocID = item.DocID

		//填充数据
		coverPlanIds[index] = blogEntity.CoverPlanId
		blogTypeIds[index] = blogEntity.BlogTypeId
		userIdIds[index] = item.UserID

		//返回的集合
		blogEntityList[index] = blogEntity
	}

	PaddingAttachemtInfo(coverPlanIds, blogEntityList)

	PaddingBlogTypeInfo(blogTypeIds, blogEntityList)

	PaddingUserInfo(userIdIds, blogEntityList) //填充用户信息

	return blogEntityList

}

//模型转化为blogSortEntity实体
func ChangeBlogSortEntityByList(blogModelList []*model.BlogModel) []*blog.BlogSortEntity {
	number := len(blogModelList)

	if number <= 0 {
		return nil
	}

	blogSortEntityList := make([]*blog.BlogSortEntity, number)

	coverPlanIds := make([]uint64, number)

	for index, item := range blogModelList {
		blogSortEntity := new(blog.BlogSortEntity)
		blogSortEntity.ID = item.ID
		blogSortEntity.CoverPlanId = uint64(item.CoverPlanId)
		blogSortEntity.Title = item.Title
		blogSortEntity.BrowseTotal = item.BrowseTotal
		blogSortEntity.CreatedAt = uint64(item.CreatedAt)
		blogSortEntity.UpdatedAt = uint64(item.UpdatedAt)

		//填充数据
		coverPlanIds[index] = blogSortEntity.CoverPlanId

		//返回的集合
		blogSortEntityList[index] = blogSortEntity
	}

	PaddingAttachemtInfoByBlogSortEntityList(coverPlanIds, blogSortEntityList)

	return blogSortEntityList
}

//模型转化为BlogListEntity实体List操作
func ChangeToBlogListEntityList(blogModelList []*model.BlogModel) []*blog.BlogListEntity {
	number := len(blogModelList)

	if number <= 0 {
		return nil
	}

	blogEntityList := make([]*blog.BlogListEntity, number)

	coverPlanIds := make([]uint64, number)

	blogTypeIds := make([]uint64, number)

	userIdIds := make([]uint64, number)

	for index, item := range blogModelList {
		blogEntity := new(blog.BlogListEntity)
		blogEntity.ID = item.ID
		blogEntity.UserId = item.UserID
		blogEntity.CreatedAt = uint64(item.CreatedAt)
		blogEntity.UpdatedAt = uint64(item.UpdatedAt)
		blogEntity.BlogTypeId = item.BlogTypeId
		blogEntity.CoverPlanId = item.CoverPlanId
		blogEntity.Title = item.Title
		blogEntity.Abstract = item.Abstract
		blogEntity.BrowseTotal = item.BrowseTotal
		blogEntity.DocID = item.DocID

		//填充数据
		coverPlanIds[index] = blogEntity.CoverPlanId
		blogTypeIds[index] = blogEntity.BlogTypeId
		userIdIds[index] = item.UserID

		//返回的集合
		blogEntityList[index] = blogEntity
	}

	PaddingAttachemtInfoByBlogListEntity(coverPlanIds, blogEntityList)

	PaddingBlogTypeInfoByBlogListEntity(blogTypeIds, blogEntityList)

	PaddingUserInfoByBlogListEntity(userIdIds, blogEntityList) //填充用户信息

	return blogEntityList

}
