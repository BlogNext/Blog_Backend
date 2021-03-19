package controller

import (
	"errors"
	"fmt"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/validate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
)

//创建一个控制器
//参数是一个执行的controller
func NewController(execController Controller) func(*gin.Context) {

	return func(context *gin.Context) {
		//通过反射创建一个新的controller，为什么要这样做？
		//因为如果不这么做，所有用户都公用一个controller，又因为gin.context是指针，所以用户第二次请求会覆盖第一次请求的context
		execControllerType := reflect.TypeOf(execController) //获取controller的指针的reflect.Type
		trueType := execControllerType.Elem()                 //获取controller的真实类型
		ptrValue := reflect.New(trueType)                       //获取controller的真实值
		controller := ptrValue.Interface().(Controller)         //底层的“值” =>  转interface{} => 再转具体类型 Controller
		//捕获异常 try/catch
		defer func() {

			if err := recover(); err != nil {

				switch err.(type) {
				case exception.MyException:
					//自定义异常，返回200,自定义的错误码
					help.Gin200ErrorResponse(context, err.(exception.MyException).GetErrorCode(), err.(exception.MyException).Error(), nil)
				default:
					//非自定义异常，返回500错误,往上抛，让gin当前运行的协程自己捕获，输出堆栈信息
					panic(err)
				}

			}
			//final
			controller.Finish() //做一些释放资源的操作
		}()

		controller.Init(context) //控制器初始化,类似构造函数

		var param []reflect.Value         // 反射调用方法所需要的参数
		action := context.Param("action") //获取执行控制器的方法
		idString := context.Param("id")

		id, err := strconv.ParseUint(idString, 10, 64) //id
		if err == nil {
			//转成功才添加
			param = make([]reflect.Value, 1)
			param[0] = reflect.ValueOf(id)
		}


		//通过反射去执行具体的方法
		value := reflect.ValueOf(controller)
		//通过路径中的action参数，去解析，调用具体的控制器方法
		action = strings.Trim(action, " ") //修剪一下空格
		var actionSplit []string
		actionSplit = strings.Split(action, "-") // '-' 符号分割方法, 例如 hello-world，谷歌是 '-'
		if len(actionSplit) == 1 {
			actionSplit = strings.Split(action, "_") // '_'符号分割方法，例如hello_world,百度是 '_'
		}



		for index, item := range actionSplit {
			actionSplit[index] = help.StrFirstToUpper(item) //首字母大写
		}

		method := strings.Join(actionSplit, "") //拼接方法

		callMethod := value.MethodByName(method)


		if !callMethod.IsValid() {
			//没有找到action参数，通过请求类型去执行具体对应的方法
			method := context.Request.Method
			panic(errors.New(fmt.Sprintf("还不支持的方法: %s", method)))
		}

		//一些钩子吧,在真正执行到控制器请求前在做一下操作，例如权限认证等
		baseException := controller.Prepare()

		prefixPath := strings.Split(context.FullPath(),":")[0]

		//转化成uri路径
		controller.SetMyFullPath(prefixPath + method)

		if baseException != nil {
			help.Gin200ErrorResponse(context, baseException.GetErrorCode(), baseException.Error(), nil)
			return //结束请求
		}

		//调用方法
		callMethod.Call(param)
		return
	}
}

//参考beego
type Controller interface {
	//控制器的生命周期
	Init(ctx *gin.Context)          //ctx是gin的Context controller是当前执行的控制器,初始化
	Prepare() exception.MyException //解析
	Finish()                        //这个函数是在执行完相应的 HTTP Method 方法之后执行的，默认是空，用户可以在子 struct 中重写这个函数，执行例如数据库关闭，清理数据之类的工作。

	//控制器的一些属性
	SetMyFullPath(path string) //设置访问的路由
}

type BaseController struct {
	Ctx        *gin.Context //gin框架的Context
	UniqFullPath string       //路由全路径,真正的路由地址，接口地址uri;转化成自己唯一的路由路径,因为自己做的url 支持 '-','_'两个解析
}

//初始化函数
func (c *BaseController) Init(ctx *gin.Context) {
	c.Ctx = ctx
}

//做一些鉴权操作等
func (c *BaseController) Prepare() exception.MyException {
	return nil
}

func (c *BaseController) Finish() {
	//可以做一些释放资源的操作
}

func (c *BaseController) SetMyFullPath(uniqFullPath string) {
	c.UniqFullPath = uniqFullPath
}

//一些重写方法,重写验证的方法，当验证出错的时候，如果是自定义的异常，则按标准的json格式返回
func (c *BaseController) ShouldBindWith(obj interface{}, b binding.Binding) {

	if err := c.Ctx.ShouldBindWith(obj, b); err != nil {
		if value, ok := obj.(validate.ValidateRequest); ok {
			panic(value.GetError(err.(validator.ValidationErrors)))
		} else {
			panic(err)
		}
	}
}
