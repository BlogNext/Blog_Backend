package model

type UserModel struct {
	BaseModel
	Nickname string `gorm:"cloumn:nickname"`
}

func (UserModel) TableName() string {
	return "user"
}
