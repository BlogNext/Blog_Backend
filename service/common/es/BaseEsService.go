package es

import (
	"context"
	"github.com/blog_backend/exception"
	"github.com/olivere/elastic/v7"
	"log"
)

const (
	BLOG_INDEX = "xiaochen_blog_next_blog"
)

//执行命令的服务
type EsExecService interface {
	//链接es
	Connet() (*elastic.Client, error)
	//命令
	SetExecCommend(commend Commend)
	GetExecCommend() Commend
	//析构函数
	Finish()
}

//回调执行命令
type Commend func(client *elastic.Client) (interface{}, error)

//运行命令的抽象流程
//result是具体命令运行返回的结果"elastic.result"
func RunCommend(esExecService EsExecService) (result interface{}, err error) {
	client, err := esExecService.Connet()
	if err != nil {
		log.Println("es出问题了", err.Error())
		return nil, exception.NewException(exception.ES_ERROR_CONNET, err.Error())
	}

	//得到回调命令
	commend := esExecService.GetExecCommend()
	result, err = commend(client)
	if err != nil {
		log.Println("回调命令出错")
		panic(err)
	}

	//执行一些释放字段的操作
	esExecService.Finish()

	return
}

type BaseEsService struct {
	//es连接
	Client elastic.Client
	//命令的执行
	Commend Commend
}

func (bs *BaseEsService) Connet() (client *elastic.Client, err error) {
	//异常捕获一下
	defer func() {
		if unknown_err := recover(); unknown_err != nil {
			log.Println("链接es失败：", unknown_err)
			client = nil
			err = unknown_err.(error)
		}
	}()

	client = es_default_connet_pool.Get().(*elastic.Client)

	return client, nil
}

func (bs *BaseEsService) SetExecCommend(commend Commend) {
	bs.Commend = commend
}

func (bs *BaseEsService) GetExecCommend() Commend {
	return bs.Commend
}

func (bs *BaseEsService) Finish() {
	es_default_connet_pool.Put(bs.Client)
}

//更新一个文档
func BuildUpdateDocCommend(index, doc_id string, doc interface{}) Commend {
	return func(client *elastic.Client) (i interface{}, err error) {
		res, err := client.Update().Index(index).Id(doc_id).Doc(doc).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

//添加一个文档
func BuildAddDocCommend(index string, entity interface{}) Commend {
	return func(client *elastic.Client) (i interface{}, err error) {
		res, err := client.Index().Index(index).BodyJson(entity).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

//删除一个文档
func BuildDeleteDocCommend(index, doc_id string) Commend {

	return func(client *elastic.Client) (i interface{}, err error) {
		res, err := client.Delete().Index(index).Id(doc_id).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

//获取一个文档
func BuildGetDocCoCommend(index, doc_id string) Commend {
	return func(client *elastic.Client) (i interface{}, err error) {
		res, err := client.Get().Index(index).Id(doc_id).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}


//构建一个回调命令
func BuildCallbackCommend(callback func() Commend) Commend {
	return callback()
}
