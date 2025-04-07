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

# Install docker
# https://kubernetes.io/zh-cn/docs/tasks/tools/install-kubectl-linux/#install-kubectl-binary-with-curl-on-linux

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
kubectl version --client --output=yaml

# package main.go
cd ~/study/week_2   
GOOS=linux GOARCH=arm go build -o webook
docker build -t Kirby980/webook:v0.0.1 .