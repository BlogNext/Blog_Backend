package model

type UsereModel struct {
	BaseModel
	Nickname string `gorm:"cloumn:nickname"`
}

func (UsereModel) TableName() string {
	return "user"
}
