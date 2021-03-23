![](https://upload-images.jianshu.io/upload_images/7236980-ec85513460bb5cd6.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/535/format/webp)


# 主从架构图解

![](https://upload-images.jianshu.io/upload_images/7236980-ec85513460bb5cd6.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/535/format/webp)

- binlog，主要服务器要的，记录写操作
- relay-log 从服务器要的，中继日志

### 统一的镜像
```
version: '3.1'

services:
  db:
    image: mysql:8.0.23
    command: --default-authentication-plugin=mysql_native_password
    restart: always
    environment:
       - MYSQL_ROOT_PASSWORD=admin123
       - MYSQL_USER=ly
       - MYSQL_PASSWORD=admin123
    ports:
       - 3306:3306
    volumes:
    - ./dataDir:/var/lib/mysql
    - ./conf.d:/etc/mysql/conf.d:ro
```


### 主服务器.cnf配置


```
[mysqld]
#默认的字符集编码
character-set-server=utf8mb4
#默认创建数据库时候使用的数据库迎请
default-storage-engine=INNODB
#数据在网络上传输的时候,数据包的大小值
max_allowed_packet=10M
#log-bin自定义名字
log-bin=ly-mysql-bin
#log-bin的格式是row模式
binlog_format=row
#日志三十天会自动删除
expire_logs_days=30
#服务器的唯一标识
server-id=1
#最大的binlog文件，超过了会开一个新的日志文件
max-binlog-size = 500M
```


### 从服务器.cnf配置

```
[mysqld]
character-set-server=utf8mb4
default-storage-engine=INNODB
max_allowed_packet=10M
#设置relay-log，docker重启的时候reldy-log名字会变
relay-log=mysql-relay-bin
expire_logs_days=30
server-id=2
max-binlog-size = 500M
```



## 从0开始做主从

### 1  主库、从库，导入上面的配置,然后重启mysql；主库开启binlog，从库开启relay-log

主要是开启主库的binlog和从库的relay-log


注意：

- docker环境搭建的，从库一定要配置relay-log,不然docker重启，relay-log的命名会不一样



### 2 主库

#### 2.1 创建slave用户

```
//创建用户
CREATE USER 'slave_account'@'%' IDENTIFIED BY 'admin123';

//授权用户为slave用户
GRANT REPLICATION SLAVE ON *.* TO 'slave_account'@'%';

//刷新权限
flush privileges;
```

#### 2.2 主库进行 FLUSH TABLES WITH READ LOCK; 锁表

```

//锁表,全局的，所有数据库，所有表
FLUSH TABLES WITH READ LOCK;

```

#### 2.3 导出要备份的数据库sql结构和数据；（mysqldump或者navicat右键）

导出数据方案

- mysqldump进行导出数据
- navicat右键导出数据和结构




#### 2.4 SHOW MASTER STATUS; 获取当前binlog文件和位置

```
mysql> SHOW MASTER STATUS;
+---------------------+----------+--------------+------------------+-------------------+
| File                | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------------+----------+--------------+------------------+-------------------+
| ly-mysql-bin.000023 |    695 |              |                  |                   |
+---------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)


```


#### 2.5 unlock tables;  释放锁表

```
//释放锁表
unlock tables;
```



### 3 从库

#### 3.1 创建对应的数据库(navicat或者自己执行命令)


```
CREATE DATABASE 数据库名;

```


#### 3.2 从库导入主库导出的数据库sql数据

导入数据方案

- mysqldump进行导入数据
- navicat右键导入数据和结构

#### 3.3 配置主库的信息

```

 CHANGE MASTER TO
        MASTER_HOST='主库ip',
        MASTER_USER='主库创建的salve用户',
        MASTER_PASSWORD='主库创建的salve密码',
        #主库 SHOW MASTER STATUS得到的log文件和位置
        MASTER_LOG_FILE='ly-mysql-bin.000023',
        MASTER_LOG_POS=695;
        
```

#### 3.4 开启slave

```
START SLAVE;
```

#### 3.5 查看是开启slave的状态

```
show slave status;


*************************** 1. row ***************************
               #这个是指slave 连接到master的状态。
               Slave_IO_State: Waiting for master to send event
                  Master_Host: 154.8.162.38
                  # 这个是master上面的一个用户。用来负责主从复制的用户 ，创建主从复制的时候建立的（具有reolication slave权限）
                  Master_User: slave_account
                  #master服务器的端口
                  Master_Port: 3306
                Connect_Retry: 60
                #I/O线程当前正在读取的主服务器二进制日志文件的名称。
              Master_Log_File: ly-mysql-bin.000023
              #在当前的主服务器二进制日志中，I/O线程已经读取的位置。
          Read_Master_Log_Pos: 2305
          #slave的SQL线程当前正在读取和执行的中继日志文件的名称。
               Relay_Log_File: mysql-relay-bin.000004
               #在当前的中继日志中，slave的SQL线程已读取和执行的位置。
                Relay_Log_Pos: 1132
        Relay_Master_Log_File: ly-mysql-bin.000023
             #Slave_IO_Running及Slave_SQL_Running进程必须正常运行，即Yes状态，否则说明同步失败
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB: 
          Replicate_Ignore_DB: 
           Replicate_Do_Table: 
       Replicate_Ignore_Table: 
      Replicate_Wild_Do_Table: 
  Replicate_Wild_Ignore_Table: 
                   Last_Errno: 0
                   Last_Error: 
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 2305
              Relay_Log_Space: 1341
              Until_Condition: None
               Until_Log_File: 
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File: 
           Master_SSL_CA_Path: 
              Master_SSL_Cert: 
            Master_SSL_Cipher: 
               Master_SSL_Key: 
               #表示主从之间的时间差 是数字的时候表示相差多少秒  null表示未知数，一般主从复制出问题了会出现null的情况。 
        Seconds_Behind_Master: 0
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error: 
               Last_SQL_Errno: 0
               Last_SQL_Error: 
  Replicate_Ignore_Server_Ids: 
             Master_Server_Id: 1
                  Master_UUID: ba542a7e-77cc-11eb-82ef-0242ac130002
             Master_Info_File: mysql.slave_master_info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
           Master_Retry_Count: 86400
                  Master_Bind: 
      Last_IO_Error_Timestamp: 
     Last_SQL_Error_Timestamp: 
               Master_SSL_Crl: 
           Master_SSL_Crlpath: 
           Retrieved_Gtid_Set: 
            Executed_Gtid_Set: 
                Auto_Position: 0
         Replicate_Rewrite_DB: 
                 Channel_Name: 
           Master_TLS_Version: 
       Master_public_key_path: 
        Get_master_public_key: 0
            Network_Namespace: 
1 row in set, 1 warning (0.01 sec)

```


## 在已有主从的情况下同步另外一个db实例


### 1. 主库选择数据库blog_ly

```
CREATE DATABASE 数据库名;
```

#### 1.1 进行锁表 FLUSH TABLES WITH READ LOCK;


```
//锁表,全局的，所有数据库，所有表
FLUSH TABLES WITH READ LOCK;
```

#### 1.2 导出sql数据

导出数据方案

- mysqldump进行导出数据
- navicat右键导出数据和结构


#### 1.3 SHOW MASTER STATUS; 查看binlog当前的文件名已经当前的位置

```

mysql> SHOW MASTER STATUS;
+---------------------+----------+--------------+------------------+-------------------+
| File                | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------------+----------+--------------+------------------+-------------------+
| ly-mysql-bin.000023 |    33697 |              |                  |                   |
+---------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

```

#### 1.4 释放表unlock tables;

```
//释放锁表，全局，所有数据库，所有表
unlock tables;
```

### 2. 从库创建数据库，导入sql文件

导入数据方案

- mysqldump进行导入数据
- navicat右键导入数据和结构


#### 2.1 停止之前的slave

```

stop slave

```

#### 2.2  重新指定binlog文件

```
 CHANGE MASTER TO
        MASTER_HOST='主库ip',
        MASTER_USER='主库创建的salve用户',
        MASTER_PASSWORD='主库创建的salve密码',
        MASTER_LOG_FILE='ly-mysql-bin.000023',
        MASTER_LOG_POS=33697;

```

#### 2.3 开启slave  

```
START SLAVE;
```

#### 2.4 查看状态

```

show slave status;

```








