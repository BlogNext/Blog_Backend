package model

type BaseModel struct {
	ID         uint  `gorm:"primary_key;AUTO_INCREMENT;not null"`
	CreateTime int64 `gorm:"cloumn:create_time"`
	UpdateTime int64 `gorm:"cloumn:update_time"`
}
