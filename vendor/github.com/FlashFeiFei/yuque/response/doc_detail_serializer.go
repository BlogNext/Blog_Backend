package response

/**
文档响应
*/
type ResponseDocDetailSerializer struct {
	Data *DocDetailSerializer `json:"data"`
}

/**
文档结构
*/
type DocDetailSerializer struct {
	ID            int64           `json:"id"`                   //文档编号
	Slug          string          `json:"slug"`                 //文档路径
	Title         string          `json:"title"`                //文章标题
	BookId        int64           `json:"book_id"`              //仓库编号，知识库编号
	Book          *BookSerializer `json:"book"`                 //仓库详情
	UserId        int64           `json:"user_id"`              //文章创建人
	User          *UserSerializer `json:"user"`                 //用户/团队信息
	Format        string          `json:"format"`               //文档格式
	Body          string          `json:"body"`                 //正文markdown源码
	BodyDraft     string          `json:"body_draft"`           //草稿markdown源码
	BodyHtml      string          `json:"body_html"`            //转换过后的正文html
	BodyLake      string          `json:"body_lake"`            //语雀lake格式的文档内容
	Public        int64           `json:"public"`               //公开级别 [0 - 私密, 1 - 公开]
	Status        int64           `json:"status"`               //状态 [0 - 草稿, 1 - 发布]
	LikesCount    int64           `json:"likes_count"`          //点赞数
	CommentsCount int64           `json:"comments_count"`       //评论总数
	CreatedAt     string          `json:"created_at"`           //文档创建时间
	UpdatedAt     string          `json:"updated_at"`           //文档更新时间
	DeletedAt     string          `json:"deleted_at,omitempty"` //删除时间，未删除为 null，在go里面是空字符串
}
