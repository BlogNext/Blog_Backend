package mysql

import (
	"fmt"
	"github.com/blog_backend/common-lib/db"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"strings"
	"time"
)

//db连接信息
//不同的key，连接不同的数据库，主从
//不需要加锁,因为只是在main入口函数导入而已
var g_db *gorm.DB

func InitDBConnect(db_info ...db.DBInfo) {
	if len(db_info) == 0 {
		return
	}

	log.Println(fmt.Sprintf("v=%v ,t=%t, p=%p", db_info, db_info, db_info))
	//mysql主从配置先写死吧。。没得时间封装
	if !strings.EqualFold(db_info[0].Key, "sources") {
		panic("mysql第一个必须是主连接配置(sources),当前是" + db_info[0].Key)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)

	//主数据库
	myGdb, err := gorm.Open(gmysql.Open(db_info[0].Dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	//从数据库
	myGdb.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{gmysql.Open(db_info[0].Dsn)},
		Replicas: []gorm.Dialector{gmysql.Open(db_info[1].Dsn)},
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
