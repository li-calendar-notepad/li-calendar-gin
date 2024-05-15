FROM node AS web_image

# 华为源
# RUN npm config set registry https://repo.huaweicloud.com/repository/npm/


WORKDIR /build

COPY ./web .

RUN npm install \
    && npm run build



FROM golang:1.21-alpine3.18 as server_image

WORKDIR /build

COPY . .

COPY --from=web_image /build/dist /build/assets/frontend

# 中国国内源
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories \
#     && go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache bash curl gcc git musl-dev

RUN go env -w GO111MODULE=on \
    && export PATH=$PATH:/go/bin \
    && go install -a -v github.com/go-bindata/go-bindata/...@latest \
    && go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest \
    && go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... \
    && go mod tidy \
    && go build -o li-calendar --ldflags="-X main.RunMode=release -X main.IsDocker=docker" main.go


FROM alpine

# 中国国内源
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories

WORKDIR /app

COPY --from=server_image /build/li-calendar /app/li-calendar

EXPOSE 9090

RUN apk add --no-cache bash ca-certificates su-exec tzdata \
    && chmod +x ./li-calendar 

CMD ./li-calendar 