

<p align="center">
<img alt="li calendar" src="./readme_src/logo.png" width="100px" />
<br>
Li Calendar - 锂日历记事本
<br>
<a title="Github" target="_blank" href="https://github.com/li-calendar">Github</a> |
<a title="Gitee" target="_blank" href="https://gitee.com/li-calendar-notepad">Gitee</a> 
</p>





## 介绍

前身[日历记事本PHP版本](https://gitee.com/hslr/calendar_notepad)，因为工作中常常要记录每天的工作日志，所以2020年上半年，抽了几天的下班时间开发了PHP版本，稳定运行了一年，但是它有些不足，2021年决定重新启动此项目，对他进行优化和增加功能并进行了技术升级。后期接触了GO+Gin+Vue3+ElementUI-Plus并重新开发了本项目 锂日历记事本。





## 目录结构

生成方式`tree --dirsfirst --charset=ascii . -d`
```
|-- api                         api文档
|   `-- v1                      v1版本
|       |-- admin               admin路由
|       |-- common              公共路由
|       |   |-- apiReturn       api返回
|       |   `-- base            基础路由
|       |-- install             安装相关
|       |-- middleware          中间件
|       `-- system              系统相关
|-- assets                      前端文件
|-- DOC                         文档相关
|-- initialize                  初始化
|-- lang                        web语言包
|-- lib                         公共库
|   |-- cache                   缓存
|   |-- captcha                 图形验证码
|   |-- cmn                     常用封装
|   |-- computerInfo            系统信息(cpu, memory)
|   |-- global                  全局配置相关
|   |-- jsonConfig              json 解析配置
|   |-- language                语言翻译相关
|   |-- mail                    邮件相关
|   |-- systemSetting           系统配置
|   `-- user                    用户相关
|-- models                      数据库模型
|   `-- datatype
`-- routers                     路由集合
    `-- admin                   admin路由
```

## 文档构建(构建中)
执行 `make swag` 或者 `swag init -g main.go` </br>
打开 [`http://127.0.0.1:9090/swagger/index.html`](http://127.0.0.1:9090/swagger/index.html) 查看文档

## 代码格式化
执行`make fmt`或者`gofmt -w -l .`(仅限Linux)

## 编译流程

1. 前端文件(编译完成的)

将前端文件存放至`./assets/frontend`

2. 编译静态文件到包内

```
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...
```

3. 构建

```
go build
```
