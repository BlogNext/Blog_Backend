package blog

import (
	"github.com/blog_backend/common-lib/arithmetic/lru"
	"time"
)

//前端博客lru缓存的一些搜索

//用于搜索框的lru缓存
var search_blog_lru *lru.LruCache

func init() {
	if search_blog_lru == nil {
		search_blog_lru = lru.New(30)

		//定时清除lru到期的key
		go func() {

			ticker := time.NewTicker(2 * time.Minute)

			for {
				select {
				case <-ticker.C:
					search_blog_lru.RemoveExpire()
				default:

				}
			}
		}()

	}
}

//添加缓存到lru
func AddBlogToLru(key lru.Key, value interface{}) {
	search_blog_lru.Add(key, value)
}

//从缓冲中获取lru
func GetBlogByLru(key lru.Key) (value interface{}, ok bool) {
	return search_blog_lru.Get(key)
}
