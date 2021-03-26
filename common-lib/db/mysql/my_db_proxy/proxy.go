package my_db_proxy

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"gorm.io/gorm"
)

//我的链接代理
//db *gorm.DB  执行sql的
//dbDryRun *gorm.DB 获取将要执行的sql的
type Proxy func(db *gorm.DB, dbDryRun *gorm.DB)

type MyDBProxy struct {
	//执行sql
	db *gorm.DB
	//只是得到sql，不真正执行
	dbDryRun *gorm.DB
}

//创建一个代理
func NewMyDBProxy() *MyDBProxy {
	return &MyDBProxy{
		db:       mysql.GetNewDB(false),
		dbDryRun: mysql.GetNewDB(true),
	}
}

//执行一个代理
func (m *MyDBProxy) ExecProxy(proxy Proxy) {
	proxy(m.db, m.dbDryRun)
}

//获取执行的sql
func (m *MyDBProxy) GetExecSql() string {
	statement := m.dbDryRun.Statement
	return m.dbDryRun.Dialector.Explain(statement.SQL.String(), statement.Vars...)
}

//生成缓存的key，缓存key生成规则:  prefix + 当前执行的sql语句
func (m *MyDBProxy) BuildCacheKey(prefix string) string {
	return prefix + m.GetExecSql()
}
