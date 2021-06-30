package blog

import (
	"github.com/blog_backend/common-lib/arithmetic/lru"
)

//前端博客lru缓存的一些搜索
//这个缓存主要用来存储一些不需要线程安全展示的
var BlgLruUnsafety *lru.LruCache

func init() {
	if BlgLruUnsafety == nil {
		BlgLruUnsafety = lru.New(30)
		//定时清除lru到期的key
		//go func() {
		//	ticker := time.NewTicker(2 * time.Minute)
		//	for {
		//		select {
		//		case <-ticker.C:
		//			BlgLruUnsafety.RemoveExpire()
		//		default:
		//		}
		//	}
		//}()

	}
}
