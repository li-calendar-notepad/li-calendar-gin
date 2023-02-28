

<p align="center">
<img alt="li calendar" src="./readme_src/logo.png" width="100px" />
<br>
Li Calendar - 锂日历记事本
<br>
<a title="Github" target="_blank" href="https://github.com/li-calendar-notepad">Github</a> |
<a title="Gitee" target="_blank" href="https://gitee.com/li-calendar-notepad">Gitee</a> 
</p>





## 介绍

前身[日历记事本PHP版本](https://gitee.com/hslr/calendar_notepad)，因为工作中常常要记录每天的工作日志，所以2020年上半年，抽了几天的下班时间开发了PHP版本，稳定运行了一年，但是它有些不足，2021年决定重新启动此项目，对他进行优化和增加功能并进行了技术升级。后期接触了GO+Gin+Vue3+ElementUI-Plus并重新开发了本项目 锂日历记事本。

## 前端代码地址
本项目不包含前端代码，前端代码是独立项目请访问：[github](https://github.com/li-calendar-notepad/li-calendar-vue) | [gitee](https://gitee.com/li-calendar-notepad/li-calendar-vue)

## 编译

#### 方式一 （通用）

1. 前端文件编译后，将dist下文件全部移植`./assets/frontend`文件夹下。编译教程请参考[前端项目](#前端代码地址)的`readme.md`文件

2. 按照[此教程](./assets/readme.md)安装工具。然后将
    `assets`文件夹编译成go文件（目的是把静态资源打包在可执行文件内）
3. 依次执行
    ```shell
    # 编译静态资源（上一步执行完成了，可以不用重复执行）
    go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... 
    
    # 开始编译，编译成功后在项目根目录生成可执行文件：li-calendar win平台: li-calendar.exe
    go build -o li-calendar main.go
    ```
#### 方式二 （Docker）推荐此方式

前提：docker环境，并且可以执行make命令，暂时不适用于windows平台

1. 将前端代码克隆在当前项目的根目录并将文件夹命名为`web`
    示例：
    ```shell
    # github
    git clone https://github.com/li-calendar-notepad/li-calendar-vue web
    
    # gitee
    git clone https://gitee.com/li-calendar-notepad/li-calendar-vue web
    ```
2. 执行make命令
    ```shell
    # 编译程序，成功后项目根目录会生成压缩包
    make build
    ```

## 运行

#### 生成配置文件：
```
# 生成配置文件（必须）
./li-calendar config

# 执行完成之后同级目录会生成两个配置文件，根据自己的需求修改`config.ini`文件内容
```

#### 运行：
```
# 运行
./li-calendar 
```

#### 访问：
浏览器打开：http://[你的域名或ip]:9090


## Docker 运行

**请将前端项目拉取到当前目录，并项目文件夹命名为`web`，否则无法编译成功**

```shell
# 编译镜像
docker build -t licalendar:latest . 

# 运行
docker run --name li-calendar -p 9090:9090 \
-v ~/licalendar/conf:/app/conf \
-v ~/licalendar/runtime:/app/runtime \
-v ~/licalendar/lang:/app/lang \
-v ~/licalendar/database:/app/database \
licalendar:latest
```

数据卷说明
```
-v ~/licalendar/conf:/app/conf # 项目配置目录
-v ~/licalendar/runtime:/app/runtime # 运行缓存、日志等
-v ~/licalendar/lang:/app/lang # 语言文件
-v ~/licalendar/database:/app/database # sqlite数据库文件夹
```