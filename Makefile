.PHONY:
swag:
	swag init -g main.go

.PHONY:
fmt:
	gofmt -l -w .
