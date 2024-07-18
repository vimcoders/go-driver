#直接下载并安装 Minikube
#如果你不想通过包安装，你也可以下载并使用一个单节点二进制文件。
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 \
  && chmod +x minikube
#将 Minikube 可执行文件添加至 path：
sudo mkdir -p /usr/local/bin/
sudo install minikube /usr/local/bin/直接下载并安装 Minikube
# 如果你不想通过包安装，你也可以下载并使用一个单节点二进制文件。
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 \
  && chmod +x minikube
# 将 Minikube 可执行文件添加至 path：
sudo mkdir -p /usr/local/bin/
sudo install minikube /usr/local/bin/

# https://docker.aityp.com/image/docker.io/kicbase/stable:v0.0.44
# Docker拉取命令
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kicbase/stable:v0.0.44
docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kicbase/stable:v0.0.44  docker.io/kicbase/stable:v0.0.44
# Containerd拉取命令
ctr images pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kicbase/stable:v0.0.44
ctr images tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kicbase/stable:v0.0.44  docker.io/kicbase/stable:v0.0.44
# 使用docker启动
minikube start --driver=docker

