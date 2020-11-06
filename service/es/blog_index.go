package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
)

//创建博客的索引
func CreateBlogIndex() Commend {
	return func(client *elastic.Client) (i interface{}, err error) {

		exists, err := client.IndexExists(BLOG_INDEX).Do(context.Background())
		if err != nil {
			// 链接出错
			panic(err)
		}

		if exists {
			//索引已经存在,返回
			return nil, errors.New(fmt.Sprintf("博客索引:%s已存在", BLOG_INDEX))
		}

		mapping := ``
		//创建索引
		createIndex, err := client.CreateIndex(BLOG_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			// 创建索引失败
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			panic("??? Not acknowledged 啥意思？？？")
		}

		return nil, nil
	}
}
