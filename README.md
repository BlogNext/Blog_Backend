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
- 2020/12/31 搜索博客接口添加lru算法，实现搜索缓存


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



# docker 部署blog_front遇到的问题


##问题描述

- 进入容器里面，手动执行npm install的时候，有些包总是下载不下载了（不是源的问题），下载一直卡着，卡了很久
- 然后npm install下载的进程就被killer掉了，没有任何的报错信息，npm install就直接被停止掉了
- 但是docker容器没有killer掉，单单只是里面的npm失败了
- 重复了很多遍，都是下载不下载，npm install，就是没有报错信息,直接被killer掉了


## 如何解决的问题，针对被killer掉的程序，

###  查看被killer程序的原因
- dmesg | grep -i -B100 'killed process'

```cassandraql
dmesg | grep -i -B100 'killed process'


#输出

[9561396.672577] [99372]     0 99372    38226      319   307200        0             0 sshd
[9561396.674020] [99375]  1000 99375    23438      368   212992        0             0 systemd
[9561396.675449] [99376]  1000 99376    77579      681   311296      151             0 (sd-pam)
[9561396.676894] [99387]  1000 99387    38226      321   299008        0             0 sshd
[9561396.678341] [99388]  1000 99388     5840      169    77824        0             0 bash
[9561396.679756] [99883]  1000 99883   287794    93719  3919872        0             0 npm
[9561396.681146] Out of memory: Kill process 99883 (npm) score 294 or sacrifice child
[9561396.682727] Killed process 99883 (npm) total-vm:1151176kB, anon-rss:374876kB, file-rss:0kB, shmem-rss:0kB

```

- 可以看到99883这个程序内存溢出了，npm 总共占用了 total-vm:1151176多的内存


### 然后查看机器的内存 free -h

```cassandraql
free -h

#输出

              total        used        free      shared  buff/cache   available
Mem:          997Mi       270Mi       233Mi       3.0Mi       492Mi       576Mi
```

- 发现自己的内存才997M，内存不足，估计是内存不足导致的npm install失败


### 开启swap增加虚拟内存


[Linux设置Swap交换分区](https://renqiang.xyz/2018/08/21/Linux%E8%AE%BE%E7%BD%AESwap%E4%BA%A4%E6%8D%A2%E5%88%86%E5%8C%BA/)

```cassandraql
dd if=/dev/zero of=/tmp/swapfile bs=1024 count=1024k

#详细操作参看上面链接
```


- 再次查看机器的内存信息


```cassandraql

free -h
              total        used        free      shared  buff/cache   available
Mem:          997Mi       271Mi       233Mi       3.0Mi       492Mi       576Mi
Swap:         2.3Gi       528Mi       1.7Gi

```


### 再次运行docker容器，这次就成功了



## 总结一下

- 对于killer掉程序，在没有任何的报错信息的前提下，可以通过 dmesg | grep -i -B100 'killed process'看原因
- free可以看到机器的内存信息
- 开启swap增大虚拟内存
- swap有个缺点，就是会有io消耗，会拖慢运行在这个机器上所有程序的性能.(没办法，穷，没钱买机器)



