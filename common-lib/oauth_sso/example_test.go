package oauth_sso

import (
	"github.com/blog_backend/common-lib/oauth_sso/core"
	"github.com/blog_backend/common-lib/oauth_sso/oauth"
	"github.com/blog_backend/common-lib/oauth_sso/token"
	"github.com/blog_backend/common-lib/oauth_sso/user"
	"testing"
)

//测试oauth模块
func TestOauthManage(t *testing.T) {
	core.SetOauthSSoSchemeConfig("http")
	core.SetOauthSSoHostConfig("127.0.0.1:8084")

	//获取预授权码
	manage := oauth.NewManage("blog_1616644960", "blog_b09bfdf65bb51bb50307f93ab930dd7708a5b6dc")
	cpacr := new(oauth.CreatePreAuthCodeResponse)
	err := manage.CreatePreAuthCode("LaughingZhu", "LaughingZhu", "http://www.baidu.com", cpacr)
	if err != nil {
		t.Fatal("创建预授权码",err)
	}

	t.Log("预授权码", cpacr)


	//预授权码换取accessToken
	pacatr := new(oauth.PreAuthCodeAccessTokenResponse)
	err = manage.PreAuthCodeAccessToken(cpacr.PreAuthCode,pacatr)
	if err != nil {
		t.Fatal("预授权码获取accessToken失败",err)
	}
	t.Log("预授权码换取accessToken",pacatr)


	//刷新预授权码
	rtr := new(oauth.RefreshTokenResponse)
	err = manage.RefreshToken(pacatr.RefreshToken,rtr)

	if err != nil {
		t.Fatal("通过refreshToken刷新AccessToken失败",err)
	}

	t.Log("通过refreshToken刷新AccessToken",rtr)
}


//测试user模块
func TestUserManagecc(t *testing.T){
	core.SetOauthSSoSchemeConfig("http")
	core.SetOauthSSoHostConfig("127.0.0.1:8084")

	//获取预授权码
	manage := oauth.NewManage("blog_1616644960", "blog_b09bfdf65bb51bb50307f93ab930dd7708a5b6dc")
	cpacr := new(oauth.CreatePreAuthCodeResponse)
	err := manage.CreatePreAuthCode("LaughingZhu", "LaughingZhu", "http://www.baidu.com", cpacr)
	if err != nil {
		t.Fatal("创建预授权码",err)
	}

	t.Log("预授权码", cpacr)


	//预授权码换取accessToken
	pacatr := new(oauth.PreAuthCodeAccessTokenResponse)
	err = manage.PreAuthCodeAccessToken(cpacr.PreAuthCode,pacatr)
	if err != nil {
		t.Fatal("预授权码获取accessToken失败",err)
	}
	t.Log("预授权码换取accessToken",pacatr)

	//token管理
	tokenManage := token.NewTokenManage(pacatr.RefreshToken,"blog_1616644960","blog_b09bfdf65bb51bb50307f93ab930dd7708a5b6dc")
	userManage := user.NewManage(tokenManage)
	uir := new(user.UserInfoResponse)
	err = userManage.UserInfo(uir)

	if err != nil {
		t.Fatal("用户信息获取失败",err)
	}
	t.Log("用户信息获取",uir)
}
