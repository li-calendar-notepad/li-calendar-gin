#!/bin/bash


# 输出脚本文件到共享目录
echo "
#!/bin/bash
cd /app 
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy 

# 环境变量
export PATH=\$PATH:/go/bin

# 静态资源编译
go install -a -v github.com/go-bindata/go-bindata/...@latest
go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...

# go run main.go

# 编译
go build -o calendar main.go
" > auto_build.sh

chmod 777 auto_build.sh

docker run --rm -it \
    -v $PWD:/app \
    --name go_build \
    golang \
    /bin/bash -c  "cd /app && ./auto_build.sh"

# 删除脚本
rm auto_build.sh
