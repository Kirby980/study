#/bin/bash

# Install golang
apt update && apt upgrade -y
apt install -y golang-1.23 
sudo ln -s /usr/lib/go-1.23/bin/go /usr/bin/go
echo export GOPROXY=https://goproxy.cn >>~/.bashrc
sources ~/.bashrc
# Install git 
sudo apt install git 
git config --global user.name "Your Name"
git config --global user.email "Your Email"

# Install web VPN 
#https://v2raya.org/docs/prologue/installation

wget -qO - https://apt.v2raya.org/key/public-key.asc | sudo tee /etc/apt/keyrings/v2raya.asc
echo "deb [signed-by=/etc/apt/keyrings/v2raya.asc] https://apt.v2raya.org/ v2raya main" | sudo tee /etc/apt/sources.list.d/v2raya.list
sudo apt update
sudo apt install v2raya v2ray ## 也可以使用 xray 包
sudo systemctl start v2raya.service
sudo systemctl enable v2raya.service

# Install docker
sudo apt update && sudo apt install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl enable --now docker
sudo docker run hello-world 

# Docker Proxy
sudo tee /etc/docker/daemon.json <<EOF
{
    "registry-mirrors": [
        "https://docker.1ms.run",
        "https://docker.xuanyuan.me"
    ]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker

# Install nodejs && typescript
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs
sudo npm install -g typescript

# Install kubernetes
# https://kubernetes.io/zh-cn/docs/tasks/tools/install-kubectl-linux/#install-kubectl-binary-with-curl-on-linux

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
kubectl version --client --output=yaml
# if did you specify the right host or port?
# Install minikube
# https://minikube.sigs.k8s.io/docs/start/?arch=%2Flinux%2Fx86-64%2Fstable%2Fbinary+download
curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
minikube start --driver=none

# 由于虚拟机非云主机不支持LoadBalancer,可以安装MetalLB部署
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.12/config/manifests/metallb-native.yaml
kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
#添加metallb-config.yaml 配置文件
kubectl apply -f metallb-config.yaml
# 数据如下: 可以通过访问192.168.3.200
# get services -o wide 
# NAME         TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)        AGE   SELECTOR
# kubernetes   ClusterIP      10.96.0.1     <none>          443/TCP        62m   <none>
# webook       LoadBalancer   10.108.93.3   192.168.3.200   81:30697/TCP   15m   app=webook

# Install helm 
# 用于部署k8s-ingress 
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx   --namespace ingress-nginx   --create-namespace
# 后续跟普通的k8s服务一样 通过kubectl apply -f 启动

# package main.go ard框架
cd ~/study/week_2   
GOOS=linux GOARCH=arm go build -o webook
#x86_64框架
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o webook .
docker build -t Kirby980/webook:v0.0.1 .