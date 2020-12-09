package mysql

import (
	"fmt"
	"github.com/blog_backend/common-lib/db"
	"github.com/gin-gonic/gin"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
)

//db连接信息
//不同的key，连接不同的数据库，主从
//不需要加锁,因为只是在main入口函数导入而已
var g_db *gorm.DB

func InitDBConnect(db_info map[string]db.DBInfo) {
	if len(db_info) == 0 {
		return
	}

	log.Println(fmt.Sprintf("v=%v ,t=%t, p=%p", db_info, db_info, db_info))
	log.Println(gin.Mode())

	//mysql一些配置
	var newLogger logger.Interface

	if gin.Mode() == "release" {
		//正式环境
		//mysql慢查询日志打印
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: 4 * time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Error,    // Log level, 生产环境下，超过阈值的sql开启慢查询
				Colorful:      true,            // 禁用彩色打印
			},
		)
	} else {
		//非正式环境
		//mysql慢查询日志打印
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Nanosecond,   // 慢 SQL 阈值
				LogLevel:      logger.Info, // Log level  warn和error级别下，slowthreshold才生效，非正式环境下开启sql打印
				Colorful:      true,          // 禁用彩色打印
			},
		)
	}

	//主数据库
	myGdb, err := gorm.Open(gmysql.Open(db_info["sources"].Dsn), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: true,
	})

	if err != nil {
		panic(err)
	}

	//从数据库
	myGdb.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{gmysql.Open(db_info["sources"].Dsn)},
		Replicas: []gorm.Dialector{gmysql.Open(db_info["replicas"].Dsn)},
		// sources/replicas 负载均衡策略
		Policy: dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200))

	if err != nil {
		panic(err)
	}

	g_db = myGdb
	log.Println(fmt.Sprintf("g_db初始化完成 v=%v ,t=%T, p=%p", g_db, g_db, g_db))
}

func GetDefaultDBConnect() *gorm.DB {
	return g_db
}
