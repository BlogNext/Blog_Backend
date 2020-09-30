# Blog_backend


技术栈

- gin 框架（自己封装了下，多了一些工具类，和路由转发）
- mysql
- es
- canal或maxwell实现，各种数据源同步(未实现)


## 目录结构说明

- common-lib 一些工具包
- config 配置文件包
- controller 自己封装的一层控制器
- entity 实体类定义，穿插接口、与服务之间
- es_dls elastic的__setting 和__mapping的定义
- exception 框架的一些自定义异常
- help 也是一些工具包，和common-lib差不多，都已经乱了，放哪都可以
- model  orm定义
- router 自己封装的一层路由
- service 服务
- sql  sql文件
- validate gin的验证类，为了定义一些公用的验证，提示中文信息的工具
- upload 项目上传的静态文件



## 迭代日志

- 2020/07/20 添加es进入,ik分词实现搜索es搜索(埋点触发es同步,未来可加入canal或maxwell)
- 2020/09/30 blog搜索添加降级功能，es挂了之后降为mysql搜索

