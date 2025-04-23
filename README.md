# 1. 构建并推送镜像
docker build -t your-registry/scheduler-extender:v1.0 .
docker push your-registry/scheduler-extender:v1.0

# 2. 部署扩展器服务
kubectl apply -f deploy/deployment.yaml

# 3. 更新kube-scheduler配置（根据集群部署方式选择合适的方式）
# 如果是kubeadm部署的集群：
cp deploy/scheduler-config.yaml /etc/kubernetes/
# 修改scheduler静态pod manifest，添加--config参数
vim /etc/kubernetes/manifests/kube-scheduler.yaml
# 在command部分添加：
# - --config=/etc/kubernetes/scheduler-config.yaml

# 4. 验证部署
kubectl get pod -n kube-system -l app=scheduler-extender
kubectl logs -n kube-system <extender-pod-name>
