package response

/**
仓库结构
*/
type BookSerializer struct {
	ID           int64           `json:"id"`            //文档编号
	Type         string          `json:"type"`          // 类型 [Book - 文档]
	Slug         string          `json:"slug"`          //知识库路径
	Name         string          `json:"name"`          //知识库名
	UserId       int64           `json:"user_id"`       //所属的团队/用户编号，仓库的创建者
	User         *UserSerializer `json:"user"`          //用户详情
	Description  string          `json:"description"`   //仓库简介
	Public       int64           `json:"public"`        //公开状态 1-公开，0-私密
	ItemsCount   int64           `json:"items_count"`   //知识库下有多少篇文章
	LikesCount   int64           `json:"likes_count"`   //喜欢数量
	WatchesCount int64           `json:"watches_count"` //订阅数量
	CreatedAt    string          `json:"created_at"`    //创建时间
	UpdatedAt    string          `json:"updated_at"`    //更新时间
}
