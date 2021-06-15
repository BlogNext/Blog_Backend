package main

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db"
	"github.com/blog_backend/common-lib/db/mysql"
	_ "github.com/blog_backend/docs"
	my_router "github.com/blog_backend/router"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"golang.org/x/net/netutil"
	"log"
	"net"
	"net/http"
)

//加载各种配置文件
func loadConfig() {
	err := config.LoadConfig("db", "config", "yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = config.LoadConfig("server", "config", "yaml")
	if err != nil {
		log.Fatal(err)
	}

	//语雀配置
	err = config.LoadConfig("yuque", "config", "yaml")
	if err != nil {
		log.Fatal(err)
	}

	//es配置可有可无
	_ = config.LoadConfig("es", "config", "yaml")
}

func loadDB() {
	dbConfig, _ := config.GetConfig("db")
	dbMap := dbConfig.GetStringMap("mysql")
	dbMapSize := len(dbMap)
	dbInfoList := make(map[string]db.DBInfo, dbMapSize)
	for key, item := range dbMap {
		dbInfoList[key] = db.DBInfo{
			Key: key,
			Dsn: item.(map[string]interface{})["dsn"].(string),
		}
	}
	mysql.InitDBConnect(dbInfoList)
}

func init() {
	//加载配置文件
	loadConfig()
	////连接数据库
	loadDB()
}

// @title 晓琛博客
// @version 1.0
// @description 改swagger用于和前端联调用的，正在努力的测试中，目前的测试securityDefinitions.apikey ApiKeyAuth只能带一个请求头
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-access-token
// @termsOfService https://github.com/BlogNext

// @contact.name Ly
// @contact.url https://github.com/FlashFeiFei?tab=repositories
// @contact.email 51785816@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host blog.laughingzhu.cn
// @BasePath /
func main() {


	router := new(my_router.MyRouter)
	

	//goin的性能分析
	//ginpprof.Wrapper(router)

	//运行服务器
	serverConfig, err := config.GetConfig("server")
	if err != nil {
		log.Fatalln(err)
	}

	serverInfo := serverConfig.GetStringMap("server")
	//gin的路由
	r := router.RunRouter()
	//url := ginSwagger.URL("http://localhost:8083/swagger/doc.json") // The url pointing to API definition
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	//用了golang github官网的包限制tcp同一时刻的连接数，这个包没有纳入基础库中，目前不知道有什么用，随便加来玩
	l,err := net.Listen("tcp",fmt.Sprintf(":%d", serverInfo["port"].(int)))
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	netutil.LimitListener(l,1000)   //最多同时只能有1000个链接，防止压垮服务器
	//http.Serve(l,r)

	h2 := serverInfo["h2"].(map[string]interface{})
	http.ServeTLS(l,r, h2["certificate"].(string),h2["certificate_key"].(string))
	//http.ListenAndServe(fmt.Sprintf(":%d", server_info["port"].(int)), r)
}
