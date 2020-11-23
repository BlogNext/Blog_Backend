package response


/**
文档响应
*/
type ResponseUserSerializer struct {
	Data *UserSerializer `json:"data"`
}


/**
用户结构
*/
type UserSerializer struct {
	ID          int64  `json:"id"`          //用户id
	Type        string `json:"type"`        //类型[`User`  - 用户, Group - 团队]
	Login       string `json:"login"`       //用户个人路径
	Name        string `json:"name"`        //用户昵称
	AvatarUrl   string `json:"avatar_url"`  //用户头像(绝对地址)
	Description string `json:"description"` //个性签名，描述
}
