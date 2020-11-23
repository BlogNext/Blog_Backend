# Blog_backend


技术栈

- gin 框架（自己封装了下，多了一些工具类，和路由转发）
- mysql
- es
- docker
- docker-compose
- maxwell实现，各种数据源同步(未实现)
- redis


### docker部署和本地运行遇到的问题

- 2020/09/30 docker go版本 golang:1.15.2，本地windows 1.4, 请求带有中文参数，中文需要urlencode
## 迭代日志

- 2020/05/24 项目docker自动部署
- 2020/06/02 妈蛋老子找错orm包了，好多山寨包呀，修改为gorm.io/gorm
- 2020/07/20 添加es进入,ik分词实现搜索es搜索(埋点触发es同步,未来可加入maxwell)
- 2020/09/30 blog搜索添加降级功能，es挂了之后降为mysql搜索
- 2020/11/05 重构数据库，加入语雀，es没有接入
- 2020/11/09 接入es
- 2020/11/23 接入swagger和前端联调,访问地址/swaager/index.html
- 2020/11/23 es太消耗内存了，总是编译的时候报错，内存不足,关掉它，项目没有了es


## MAXWELL和redis部署

[maxwell官网](http://maxwells-daemon.io/quickstart/)

- docker部署maxwells参考 maxwell-docker-compose.yml

```cassandraql
解释下对应的sql
//创建用户mawell ‘%’标识任何ip都能连,密码
mysql> CREATE USER 'maxwell'@'%' IDENTIFIED BY 'XXXXXX';
//授权maxwell数据库给maxwell用户,maxwell数据库用于记录一些maxwell作为模拟mysql slave的信息
mysql> GRANT ALL ON maxwell.* TO 'maxwell'@'%';
//未用户赋予slave从库权限
mysql> GRANT SELECT, REPLICATION CLIENT, REPLICATION SLAVE ON *.* TO 'maxwell'@'%';
```

##测试
```cassandraql
//监听redis maxwell的订阅频道
SUBSCRIBE maxwell

```

## es遇到的问题

```cassandraql
elastic: Error 429 (Too Many Requests): [parent] Data too large, data for [<
```

原因: 
- fileddata占用内存过大,解决办法如下

```cassandraql
PUT _cluster/settings
{
  "persistent" : {
    "indices.breaker.fielddata.limit" : "20%" 
  }
}

```


## nginx

```cassandraql
server
 {
     listen       443 ssl;
     
    # ssl_certificate    /home/xiaochen/BlogNext/ssl_certs/laughingzhu.com/blog/server.crt;
    # ssl_certificate_key  /home/xiaochen/BlogNext/ssl_certs/laughingzhu.com/blog/server.key;
    ssl_certificate    /etc/letsencrypt/live/laughingzhu.com/fullchain.pem;
    ssl_certificate_key  /etc/letsencrypt/live/laughingzhu.com/privkey.pem;
    
    server_name blog.laughingzhu.com;


    add_header Access-Control-Allow-Origin $http_origin;
    add_header 'Access-Control-Allow-Credentials' 'true';
    #add_header 'Access-Control-Allow-Headers' 'Authorization,Content-Type,Accept,Origin,User-Agent,DNT,Cache-Control,X-Mx-ReqToken,X-Requested-With';
    add_header 'Access-Control-Allow-Headers' '*';
    add_header 'Access-Control-Allow-Methods' 'GET,POST,OPTIONS';
    access_log  /var/log/nginx/blog_next_access.log  main;

    location /upload/ {
        expires 30s;
        root /home/xiaochen/BlogNext/code/Blog_Backend/;
   }


   location /swagger {

    auth_basic "swagger登录";
    auth_basic_user_file /home/xiaochen/BlogNext/nginx/htpasswd;
    proxy_pass  http://127.0.0.1:8083;

   }

   location / {

     proxy_pass  http://127.0.0.1:8083;
   }
   
}
```


- /home/xiaochen/BlogNext/nginx/htpasswd文件的生成

```cassandraql
printf "admin:$(openssl passwd -crypt 123456)\n" >> htpasswd
```

