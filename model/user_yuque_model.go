package model

type UserYuQueModel struct {
	BaseModel
	UserId      int64  `gorm:"cloumn:user_id"`
	Login       string `gorm:"cloumn:login"`
	Name        string `gorm:"cloumn:name"`
	AvatarUrl   string `gorm:"cloumn:avatar_url"`
	Description string `gorm:"cloumn:description"`
}

func (UserYuQueModel) TableName() string {
	return "user_yuque"
}
