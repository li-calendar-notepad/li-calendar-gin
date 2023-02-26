## 静态文件编译文件夹


### 1. 安装（...必须带上）
```ssh
go get github.com/go-bindata/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...

# go版本>=1.17 使用intsall方式
go install -a -v github.com/go-bindata/go-bindata/...@latest
go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest
```
### 2. 安装成功后将 GOPATH/bin 加入环境变量

参考各自系统环境变量配置即可


### 3. 压缩静态文件 到 asset目录
```ssh
# 测试
# go-bindata-assetfs -debug -o=assets/bindata.go -pkg=assets static/... view/...
go-bindata-assetfs -debug -o=assets/bindata.go -pkg=assets assets/... 

# 正式
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... 
```
> 正式环境需要 去掉` -debug `

#### 参考文章
Go | Go 语言打包静态文件以及如何与Gin一起使用Go-bindata
https://www.jianshu.com/p/a7f5885679ef

[golang]Go内嵌静态资源go-bindata的安装及使用
https://www.cnblogs.com/landv/p/11577213.html