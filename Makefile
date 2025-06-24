# FILEPATH: /Users/yangyang/go/src/github.com/github-wyy/m-scheduler-extender/Makefile
# 定义变量
IMAGE_NAME ?= ccr.ccs.tencentyun.com/mervynwang/m-scheduler-extender
IMAGE_TAG ?= v1.0.2

.PHONY: build docker-build docker-push deploy undeploy test clean help

## 编译项目
build:
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -o m-scheduler-extender .

## 构建Docker镜像
docker-build: build
	@echo "Building Docker image..."
	docker build --platform linux/amd64 -t $(IMAGE_NAME):$(IMAGE_TAG) .

## 推送Docker镜像
docker-push: docker-build
	@echo "Pushing Docker image..."
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

## 部署到Kubernetes集群
deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f manifests/rbac.yaml
	kubectl apply -f manifests/deploy.yaml

## 卸载部署
undeploy:
	@echo "Removing deployment..."
	kubectl delete -f manifests/deploy.yaml
	kubectl delete -f manifests/rbac.yaml

## 运行单元测试
test:
	@echo "Running tests..."
	go test -v ./...

## 清理生成文件
clean:
	@echo "Cleaning up..."
	rm -f m-scheduler-extender

## 显示帮助信息
help:
	@echo "可用命令:"
	@echo "  build        - 编译Go项目"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-push  - 推送Docker镜像到仓库"
	@echo "  deploy       - 部署到Kubernetes集群"
	@echo "  undeploy     - 从集群移除部署"
	@echo "  test         - 运行单元测试"
	@echo "  clean        - 清理生成文件"
	@echo ""
	@echo "变量覆盖示例:"
	@echo "  make docker-build IMAGE_NAME=myregistry/extender IMAGE_TAG=latest"
