package user

import (
	"github.com/blog_backend/entity"
)

//User实体文档
type UserEntity struct {
	entity.BaseEntity
	Nickname string `json:"nickname"`
}
