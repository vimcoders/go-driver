# 安装依赖包以使apt可以通过HTTPS使用repository：
sudo apt-get install -y apt-transport-https ca-certificates curl
# 首先，导入阿里云的apt-key
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
# 接着，添加阿里云的Kubernetes APT源到系统中：
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
#关闭防火墙，禁用防火墙开机自启动
systemctl stop firewalld
systemctl disable firewalld

# 临时禁用SeLinux，重启失效
setenforce 0
# 修改SeLinux配置，永久禁用
sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/selinux/config

# 临时关闭Swap
swapoff -a
# 修改 /etc/fstab 删除或者注释掉swap的挂载，可永久关闭swap
sed -i '/swap/s/^/#/' /etc/fstab

#修改k8s.conf
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system
# 安装docker
sudo apt-get install docker-ce docker-ce-cli containerd.io
# 更新软件包列表并安装Kubernetes组件：
sudo apt update
sudo apt install -y kubelet kubeadm kubectl
# 禁止这些组件自动更新：
sudo apt-mark hold kubelet kubeadm kubectl
# 初始化Kubernetes master节点：
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
# 镜像拉取失败的处理方式
sudo kubeadm init --apiserver-advertise-address=192.168.16.132 --image-repository registry.aliyuncs.com/google_containers --service-cidr=10.96.0.0/12 --pod-network-cidr=10.244.0.0/16
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
# 通过如下指令创建默认的kubeadm-config.yaml文件
kubeadm config print init-defaults  > kubeadm-config.yaml
journalctl -xeu kubelet -l #查看详细错误信息
# 重新初始化
kubeadm reset
# 设置kubectl的配置文件：
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
# 安装Pod网络插件（例如Flannel）：
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
# 检查Kubernetes集群状态：
kubectl get nodes
kubectl get pods --all-namespaces
# 查看所需要的镜像与版本
kubeadm config images list
# 重启docker
systemctl restart docker
# 重新加载配置文件
sudo systemctl daemon-reload
# 重新加载容器
systemctl restart containerd
# 查看镜像文件
sudo containerd config default | grep "sandbox"
# 设置主机名字
sudo hostnamectl set-hostname node1
# 指定配置文件
kubectl --kubeconfig ~/.kube/config  get nodes
# 查看节点状态
kubectl get pod -n kube-system

