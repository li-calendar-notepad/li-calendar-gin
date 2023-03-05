

<p align="center">
<img alt="li calendar" src="./readme_src/logo.png" width="100px" />
<br>
Li Calendar - 锂日历记事本
<br>
<a title="Github" target="_blank" href="https://github.com/li-calendar-notepad">Github</a> |
<a title="Gitee" target="_blank" href="https://gitee.com/li-calendar-notepad">Gitee</a> 
</p>





## 🛸 介绍

前身[日历记事本PHP版本](https://gitee.com/hslr/calendar_notepad)，因为工作中常常要记录每天的工作日志，所以2020年上半年，抽了几天的下班时间开发了PHP版本，稳定运行了一年，但是它有些不足，2021年决定重新启动此项目，对他进行优化和增加功能并进行了技术升级。后期接触了GO+Gin+Vue3+ElementUI-Plus并重新开发了本项目 锂日历记事本。

## 🌱 相对PHP版本增加以及准备做的
- [x] 全新UI
- [x] 内容选用高级编辑器支持传文件，粘贴图片
- [x] 深色模式支持
- [x] 强化事件模板功能，并支持拖拽
- [x] 风格支持自定义，支持导入导出
- [x] 节假日改为特殊日期，可自定义上传
- [x] docker运行
- [ ] 按时间范围，分享日历视图
- [ ] 设置待办邮件提醒
- [ ] 速记功能
- [ ] 事件时间线视图
- [ ] 单事件收藏、分享


## ⌨️ 前端源码
项目进行了前后端分离，所以本源码不包含前端，前端是由`Vue3`+`Element-UI Plus`+`Fullcalendar`，前端项目源码请访问：[github](https://github.com/li-calendar-notepad/li-calendar-vue) | [gitee](https://gitee.com/li-calendar-notepad/li-calendar-vue)

## 🚥 说明
目前项目仍处于开发阶段，部分功能未完善，欢迎体验，有问题可以提Issues，暂时不建议作为正式项目使用。

## 🖼️ 截图

日历首页
<img alt="li calendar" style="border:1px solid #dce1e4;" src="./readme_src/screenshot/item_home.png" />

深色模式
<img alt="li calendar" style="border:1px solid #dce1e4;" src="./readme_src/screenshot/dark.png" />

事件内容
<img alt="li calendar" style="border:1px solid #dce1e4;" src="./readme_src/screenshot/event_content.png" />

事件模板
<img alt="li calendar" style="border:1px solid #dce1e4;" src="./readme_src/screenshot/item_home_model.png" />

## 💾 编译

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

## 🚄 运行

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


## 💎 Docker 运行

**请将前端项目拉取到当前目录，并将前端项目命名为`web`，否则无法编译成功**

```shell
# 编译镜像
docker build -t licalendar:latest . 

# 运行
docker run --name li-calendar -p 9090:9090 \
-v ~/licalendar/conf:/app/conf \
-v ~/licalendar/uploads:/app/uploads \
-v ~/licalendar/runtime:/app/runtime \
-v ~/licalendar/lang:/app/lang \
licalendar:latest
```

更多数据卷说明
```
-v ~/licalendar/conf:/app/conf # 项目配置目录
-v ~/licalendar/uploads:/app/uploads # 上传的文件目录
-v ~/licalendar/runtime:/app/runtime # 运行缓存、日志等
-v ~/licalendar/lang:/app/lang # 语言文件目录
-v ~/licalendar/database:/app/database # sqlite数据库目录
```

## 🎁 打赏

开源不易，如果你喜欢本项目或者觉得项目对你有帮助，欢迎进行[🧧打赏作者🧧](https://blog.enianteam.com/u/sun/content/11#%E6%89%93%E8%B5%8F)。记得加作者留名。在此感谢