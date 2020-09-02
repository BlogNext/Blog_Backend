# Blog_backend


技术栈

- gin 框架（自己封装了下，多了一些工具类，和路由转发）
- mysql
- es


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


