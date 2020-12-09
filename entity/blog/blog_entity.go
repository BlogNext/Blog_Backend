package blog

import (
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/entity/user"
)


//blogSort实体
type BlogSortEntity struct {
	entity.BaseEntity
	//文章标题
	Title string `json:"title"`
	BrowseTotal uint `json:"browse_total"` //浏览量
	//封面图id
	CoverPlanId uint64 `json:"cover_plan_id"`
	//封面图信息
	CoverPlanInfo *attachment.AttachmentEntity `json:"cover_plan_info"`
}


//blog列表实体
type BlogListEntity struct {
	entity.BaseEntity
	UserId     uint64 `json:"user_id"`
	BlogTypeId uint64 `json:"blog_type_id"`
	//文章标题
	Title string `json:"title"`
	//文章摘要
	Abstract string `json:"abstract"`

	BrowseTotal uint `json:"browse_total"` //浏览量

	//封面图id
	CoverPlanId uint64 `json:"cover_plan_id"`

	//封面图信息
	CoverPlanInfo *attachment.AttachmentEntity `json:"cover_plan_info"`

	BlogTypeObject *BlogTypeEntity `json:"blog_type_object"`

	//用户信息
	UserInfo *user.UserEntity `json:"user_info"`
}

//blog详情文档
type BlogEntity struct {
	entity.BaseEntity

	UserId     uint64 `json:"user_id"`
	BlogTypeId uint64 `json:"blog_type_id"`

	YuqueFormat string `json:"yuque_format"`
	//文章标题
	Title string `json:"title"`
	//文章摘要
	Abstract string `json:"abstract"`
	//文章内容
	Content     string `json:"content"`
	BrowseTotal uint   `json:"browse_total"` //浏览量

	//封面图id
	CoverPlanId uint64 `json:"cover_plan_id"`

	//封面图信息
	CoverPlanInfo *attachment.AttachmentEntity `json:"cover_plan_info"`

	BlogTypeObject *BlogTypeEntity `json:"blog_type_object"`

	//用户信息
	UserInfo *user.UserEntity `json:"user_info"`
}

//blog_type文档
type BlogTypeEntity struct {
	entity.BaseEntity
	Title string `json:"title"`
}
