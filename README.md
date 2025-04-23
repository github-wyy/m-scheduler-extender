```shell
# 1. 构建并推送镜像
make build
make docker-build
make docker-push

# 2. 部署扩展器服务
make deploy

# 3. 更新kube-scheduler配置（根据集群部署方式选择合适的方式）
参考 manifests/scheduler-config.yaml

# 4. 验证部署
调度pod，查看调度器日志

```