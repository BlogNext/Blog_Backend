# user-system-backend
用户系统，cas单点登录，oauth2.0

- gin 1.6版本

# oauth2.0只完成的

- authorization code grant 预授权码模式
- implicati Grant token颁发之简化模式
- client credential grant 以后可以考虑要

## 实现方案
- 采用jwt实现
- token和refreshToken的数据其实是一样的,区别在于token有过期时间，refreshToken没有过期时间


# gin1.6的一些改造
- v10的验证器，做了中文提示
- 简单的反射封装了一下路由到控制器
- validate提取出来
- 封装了一些业务的异常处理

# 一张好看的流程图

![image](./一张好看的流程图.jpg)