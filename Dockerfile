FROM node:16.17.1-alpine3.15 as web_image

WORKDIR /build

COPY ./web .

RUN npm install \
    && npm run build



FROM golang:1.19 as server_image

WORKDIR /build

COPY . .

COPY --from=web_image /build/dist /build/assets/frontend

# 执行指令 关闭链接确认
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && export PATH=$PATH:/go/bin \
    && go install -a -v github.com/go-bindata/go-bindata/...@latest \
    && go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest \
    && go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... \
    && go build -o li-calendar --ldflags="-X main.RunMode=release main.IsDocker=docker" main.go


FROM centos

WORKDIR /app

COPY --from=server_image /build/li-calendar /app/li-calendar

RUN ./li-calendar config

CMD ./li-calendar 