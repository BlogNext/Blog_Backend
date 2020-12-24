package user

import "github.com/blog_backend/entity"

//User前端用户实体
type UserFrontEntity struct {
	entity.BaseEntity
	UserId uint64 `json:"user_id"` //用户id
	Login  string `json:"login"`   //语雀的login
	Name   string `json:"name"`    //语雀名称
}
