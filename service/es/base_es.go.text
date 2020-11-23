package es

import (
	"context"
	"encoding/json"
	"errors"
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
	Finish(*elastic.Client)
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
	esExecService.Finish(client)

	return
}

type BaseEsService struct {
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

func (bs *BaseEsService) Finish(client *elastic.Client) {
	es_default_connet_pool.Put(client)
}

//删除一个文档
func (bs *BaseEsService) DeleteDoc(index, doc_id string) (*elastic.DeleteResponse, error) {
	//构建一个命令
	commend := BuildDeleteDocCommend(index, doc_id)
	//设置命令
	bs.SetExecCommend(commend)
	//运行命令
	resule, err := RunCommend(bs)

	if err != nil {
		return nil, err
	}

	return resule.(*elastic.DeleteResponse), nil
}

//添加一个文档
func (bs *BaseEsService) AddDoc(index string, doc interface{}) (*elastic.IndexResponse, error) {
	if doc == nil {
		return nil, errors.New("blog_doc为空")
	}

	//构建一个命令
	commend := BuildAddDocCommend(index, doc)
	//设置命令
	bs.SetExecCommend(commend)
	//运行命令
	resule, err := RunCommend(bs)

	if err != nil {
		return nil, err
	}

	return resule.(*elastic.IndexResponse), nil
}

//通过mysql主键id获取一个es文档
//index 索引
//id mysql主键
//赋值的dock文档
func (bs *BaseEsService) GetDocByMysqlId(index string, id uint, doc_result interface{}) error {

	//构建一个命令
	commend := BuildGetDocCommendByMysql(index, id, doc_result)
	//设置命令
	bs.SetExecCommend(commend)
	//运行命令
	_, err := RunCommend(bs)

	if err != nil {
		return err
	}

	return nil
}

//更新一个文档
func (bs *BaseEsService) UpdateDoc(index, doc_id string, doc interface{}) (*elastic.UpdateResponse, error) {
	if doc == nil {
		return nil, errors.New("blog_doc为空")
	}

	//构建一个命令
	commend := BuildUpdateDocCommend(index, doc_id, doc)
	//设置命令
	bs.SetExecCommend(commend)
	//运行命令
	resule, err := RunCommend(bs)

	if err != nil {
		return nil, err
	}

	return resule.(*elastic.UpdateResponse), nil
}

//搜索文档
func (bs *BaseEsService) SearchDoc(callback func() Commend) (*elastic.SearchResult, error) {

	commend := BuildCallbackCommend(callback)

	//设置命令
	bs.SetExecCommend(commend)

	//运行命令
	commend_result, err := RunCommend(bs)
	if err != nil {
		return nil, err
	}

	return commend_result.(*elastic.SearchResult), nil
}

//一些命令的回调

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
func BuildGetDocCommend(index, doc_id string) Commend {
	return func(client *elastic.Client) (i interface{}, err error) {
		res, err := client.Get().Index(index).Id(doc_id).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func BuildGetDocCommendByMysql(index string, id uint, doc_result interface{}) Commend {
	return func(client *elastic.Client) (i interface{}, err error) {
		boolQ := elastic.NewBoolQuery()
		boolQ.Filter(elastic.NewTermQuery("id", id))
		res, err := client.Search().Index(index).Query(boolQ).Pretty(true).Do(context.Background())
		if err != nil {
			return nil, err
		}

		if res.TotalHits() <= 0 {
			return nil, errors.New("es查询不到单个文档数据")
		}

		json.Unmarshal(res.Hits.Hits[0].Source, doc_result)
		return nil, nil
	}
}

//构建一个回调命令
func BuildCallbackCommend(callback func() Commend) Commend {
	return callback()
}
