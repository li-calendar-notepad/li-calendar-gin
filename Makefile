IMAGE_NAME=li-calendar-build
VERSION=v1.0
APP_NAME=li-calendar

.PHONY:
swag:
	swag init -g main.go

.PHONY:
fmt:
	gofmt -l -w .

.PHONY:
build: build-web build-server
	@echo "编译完成，正在打包..."
	rm -rf ${PWD}/assets/frontend
	mv -f ${PWD}/web/dist ${PWD}/assets/frontend
	tar -czvf ${APP_NAME}_${VERSION}.tar.gz  ${APP_NAME}
	@echo "打包完成 - ${PWD}/${APP_NAME}_${VERSION}.tar.gz"
.PHONY:
build-web: 
	@echo "正在构建web..."
	docker run --rm -it \
		-v ${PWD}/web/:/app \
		node \
		/bin/sh -c "cd /app && npm install && npm run build"

.PHONY:
build-server: build-server-image
	@echo "正在构建server..."
	docker run --rm -it -e app_name=${APP_NAME} -v ${PWD}:/root/build ${IMAGE_NAME}
	tar -czvf ${APP_NAME}_${VERSION}.tar.gz  ${APP_NAME}

.PHONY:
build-server-image:
	docker build -t ${IMAGE_NAME} -f Dockerfile.base .

