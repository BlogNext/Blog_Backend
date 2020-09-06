package main

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db"
	"github.com/blog_backend/common-lib/db/mysql"
	my_router "github.com/blog_backend/router"
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

	//es配置可有可无
	_ = config.LoadConfig("es","config","yaml")
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
	http.ListenAndServe(fmt.Sprintf(":%d", server_info["port"].(int)), router.RunRouter())
}
