package my_db_proxy

import (
	"gorm.io/gorm"
)

//我的链接代理
//db *gorm.DB  执行sql的
type Proxy func(db *gorm.DB)

type MyDBProxy struct {
	//执行sql
	db *gorm.DB
}

//创建一个代理
func NewMyDBProxy() *MyDBProxy {
	return nil
	//proxy := &MyDBProxy{
	//	db: mysql.GetNewDB(false),
	//}
	//
	//return proxy
}

//创建代理
func NewMyDBProxyByTable(table string) *MyDBProxy {
	proxy := NewMyDBProxy()
	proxy.ExecProxy(func(db *gorm.DB) {
		//需要改变一下db的内存值，gorm的clone值的问题
		*db = *db.Table(table)
	})
	return proxy
}

//创建代理
//func NewMyDBProxyByModel(model interface{}) *MyDBProxy {
//	proxy := NewMyDBProxy()
//	proxy.ExecProxy(func(db *gorm.DB) {
//		//需要改变一下db的内存值，gorm的clone值的问题
//		*db = *db.Model(model)
//	})
//	return proxy
//}

//执行一个代理
//返回一个interface
func (m *MyDBProxy) ExecProxy(proxy Proxy) {
	proxy(m.db)
	return
}
