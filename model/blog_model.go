package model

const (
	//语雀公开级别
	BLOG_MODEL_YUQUE_PUBLIC_0 = iota //私密的
	BLOG_MODEL_YUQUE_PUBLIC_1        //公开的
)

type BlogModel struct {
	BaseModel
	UserID      uint   `gorm:"cloumn:user_id"` //userid
	DocID       string `gorm:"cloumn:doc_id"`  //es文档id
	CoverPlanId int64  `gorm:"cloumn:cover_plan_id"`
	BlogTypeId  int64  `gorm:"cloumn:blog_type_id"`
	YuqueId     int64  `gorm:"cloumn:yuque_id"` //语雀文档id
	YuqueSlug   string `gorm:"cloumn:yuque_slug"`
	YuqueFormat string `gorm:"cloumn:yuque_format"`
	YuqueLake   string `gorm:"cloumn:yuque_lake"`
	YuquePublic int    `gorm:"cloumn:yuque_public"`
	Title       string `gorm:"cloumn:title"`
	BrowseTotal uint   `gorm:"cloumn:browse_total"` //浏览量
	Abstract    string `gorm:"cloumn:abstract"`
	Content     string `gorm:"cloumn:content"`
}

func (BlogModel) TableName() string {
	return "blog"
}
