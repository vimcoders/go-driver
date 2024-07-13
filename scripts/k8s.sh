# 安装依赖包以使apt可以通过HTTPS使用repository：
sudo apt-get install -y apt-transport-https ca-certificates curl
# 首先，导入阿里云的apt-key
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
# 接着，添加阿里云的Kubernetes APT源到系统中：
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
# 更新软件包列表并安装Kubernetes组件：
sudo apt update
sudo apt install -y kubelet kubeadm kubectl
# 禁止这些组件自动更新：
sudo apt-mark hold kubelet kubeadm kubectl
# 初始化Kubernetes master节点：
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
# 设置kubectl的配置文件：
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
# 安装Pod网络插件（例如Flannel）：
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
# 检查Kubernetes集群状态：
kubectl get nodes
kubectl get pods --all-namespaces