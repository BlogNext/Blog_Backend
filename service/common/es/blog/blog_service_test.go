package blog

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"testing"
)

var client *elastic.Client

func TestNewClient(t *testing.T) {

	es_client, err := elastic.NewClient(
		elastic.SetURL("http://106.12.76.73:9200/"),
		elastic.SetBasicAuth("elastic", "你的密码"),
		elastic.SetTraceLog(log.New(os.Stdout, "ly_es", 0)), //跟踪请求和响应细节
	)
	if err != nil {
		panic(err)
	}

	client = es_client
	t.Log("连接成功")
	return
}

//ping一下es是否能ping通
func TestEsConnent(t *testing.T) {

	// Ping Elasticsearch服务器以获得节点信息
	info, code, err := client.Ping("http://106.12.76.73:9200/").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}

	t.Logf("Elasticsearch 返回信息， 状态码：code %d 版本号 %s\n", code, info.Version.Number)

}

//索引是否存在
func TestIndexExists(t *testing.T) {
	_,err := client.IndexExists("xiaochen_blog_next_blog").Do(context.Background())
	if err != nil {
		panic(err)
	}

	t.Log("xiaochen_blog_next_blog存在")
}
