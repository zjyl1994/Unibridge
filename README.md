# Unibridge
Unibridge for XJCraft

适用于小鸡服的统一登录桥，方便第三方接入用户名密码系统

## 安装

下载编译后，启动UniBridge即可。建议使用第三方工具监视进程，崩溃后重启。
注意：.env文件包含所需的配置项，如果找不到.env文件则会从环境变量中读取对应的值。

- UNIBRIDGE_LISTEN 为Web服务监听的地址和端口
- UNIBRIDGE_DSN 为连接mysql数据库所用的连接字符串

连接字符串完整格式为 username:password@protocol(address)/dbname?charset=utf8&parseTime=True&loc=Local
(&parseTime=True&loc=Local 不可以缺少)

## 使用

HTTP GET 访问 `http://[UNIBRIDGE_LISTEN]/checkpass?name=xxx&pass=xxx`

- name 为玩家的游戏ID，游戏ID区分大小写
- pass 为玩家密码的哈希结果,hex(sha256(玩家ID+玩家密码)),哈希结果不分大小写

当密码正确时：
```json
{
    "success": true,
    "name": "zjyl1994",
    "lastAction": 1582292632
}
```
当发生错误时：
```json
{
    "success": false,
    "reason": "错误原因"
}
```
可根据success查看是否密码是否通过。

当密码错误尝试超过3次时，系统会锁定，此时需要玩家在MC中登录一次才能重置密码重试次数。