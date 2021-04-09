package oauth_sso

import (
	"github.com/blog_backend/common-lib/oauth_sso/oauth"
	"testing"
)

func TestOauthManage(t *testing.T) {
	SetOauthSSoSchemeConfig("http")
	SetOauthSSoHostConfig("127.0.0.1")

	manage := oauth.NewManage("blog_1616644960", "blog_b09bfdf65bb51bb50307f93ab930dd7708a5b6dc")
	cpacr := new(oauth.CreatePreAuthCodeResponse)
	err := manage.CreatePreAuthCode("LaughingZhu", "LaughingZhu", "", cpacr)
	if err != nil {
		t.Error(err)
	}

	t.Log("预授权码", cpacr)
}
