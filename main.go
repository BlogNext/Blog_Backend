package main

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db"
	"github.com/blog_backend/common-lib/db/mysql"
	_ "github.com/blog_backend/docs"
	my_router "github.com/blog_backend/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
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
	db_config, _ := config.GetConfig("db")
	db_map := db_config.GetStringMap("mysql")
	db_map_size := len(db_map)
	db_info_list := make(map[string]db.DBInfo, db_map_size)
	for key, item := range db_map {
		db_info_list[key] = db.DBInfo{
			Key: key,
			Dsn: item.(map[string]interface{})["dsn"].(string),
		}
	}
	mysql.InitDBConnect(db_info_list)
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
// @name token
// @termsOfService https://github.com/BlogNext

// @contact.name Ly
// @contact.url https://github.com/FlashFeiFei?tab=repositories
// @contact.email 51785816@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /
func main() {

	router := new(my_router.MyRouter)

	//goin的性能分析
	//ginpprof.Wrapper(router)

	//运行服务器
	server_config, err := config.GetConfig("server")
	if err != nil {
		log.Fatalln(err)
	}

	server_info := server_config.GetStringMap("servier")
	//gin的路由
	r := router.RunRouter()
	url := ginSwagger.URL("http://localhost:8083/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	http.ListenAndServe(fmt.Sprintf(":%d", server_info["port"].(int)), r)
}
