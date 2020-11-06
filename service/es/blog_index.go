package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

//创建博客的索引
func CreateBlogIndex() Commend {
	return func(client *elastic.Client) (i interface{}, err error) {

		exists, err := client.IndexExists(BLOG_INDEX).Do(context.Background())
		if err != nil {
			// 链接出错
			log.Println(err)
			return nil, errors.New(fmt.Sprintf("创建博客索引失败，链接出错"))
		}

		if exists {
			//索引已经存在,返回
			return nil, errors.New(fmt.Sprintf("博客索引:%s已存在", BLOG_INDEX))
		}

		//blog索引
		mapping := `
{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "my_ik_synonym": {
          "type": "custom",
          "tokenizer": "ik_max_word"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "abstract": {
        "type": "text",
        "analyzer": "my_ik_synonym"
      },
      "title": {
        "type": "text",
        "analyzer": "my_ik_synonym"
      },
      "content": {
        "type": "text",
        "analyzer": "my_ik_synonym"
      },
      "user_info": {
        "type": "object",
        "properties": {
          "nickname": {
            "type": "text",
            "analyzer": "my_ik_synonym"
          }
        }
      }
    }
  }
}

`
		//创建索引
		createIndex, err := client.CreateIndex(BLOG_INDEX).Body(mapping).Do(context.Background())

		if err != nil {
			// 创建索引失败
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			//可能是由于网络原因，中断了，重试就好了
			panic("??? Not acknowledged 啥意思？？？")
		}

		return nil, nil
	}
}
