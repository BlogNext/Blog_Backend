package model

type BaseModel struct {
	ID         uint  `gorm:"primary_key;AUTO_INCREMENT;not null"`
	CreateTime int64 `gorm:"cloumn:created_at,autoCreateTime"`
	UpdateTime int64 `gorm:"cloumn:updated_at,autoUpdateTime"`
}
