package blog

import (
	"github.com/blog_backend/common-lib/arithmetic/lru"
	"time"
)

//前端博客lru缓存的一些搜索

//用于搜索框的lru缓存
var isInit bool
var SearchBlgLru *lru.LruCache
var ListBlogLru *lru.LruCache

func init() {
	if isInit == false {
		isInit = true
		SearchBlgLru = lru.New(30)
		ListBlogLru = lru.New(30)

		//定时清除lru到期的key
		go func() {

			ticker := time.NewTicker(2 * time.Minute)

			for {
				select {
				case <-ticker.C:
					SearchBlgLru.RemoveExpire()
					ListBlogLru.RemoveExpire()
				default:

				}
			}
		}()

	}
}
