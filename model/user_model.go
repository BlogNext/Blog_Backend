package model

type UsereModel struct {
	BaseModel
	NickName string `gorm:"cloumn:nickname"`
}

func (UsereModel) TableName() string {
	return "user"
}
