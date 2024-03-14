#!/bin/bash

#set -x
# 默认参数值
IP_ADDR_MASTER="192.168.3.180"
IP_ADDR_NODE1="192.168.3.181"
IP_ADDR_NODE2="192.168.3.182"
GATEWAY="192.168.3.1"

print_usage() {
    echo "Usage: $0 [master | node <id> reset] [--help]"
}

# 判断参数个数
if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

set_hosts(){
hosts=("$IP_ADDR_MASTER master" "$IP_ADDR_NODE1 node1" "$IP_ADDR_NODE2 node2")
for item in "${hosts[@]}"; do
    # 检查是否已经存在相同的记录
    if ! grep -q "$item" /etc/hosts; then
        # 追加内容到 /etc/hosts 文件
        echo "$item" | sudo tee -a /etc/hosts
    fi
done
echo "设置hosts... 完成"
}

config_k8s_source(){
  # 安装k8s源
  if [ ! -f /etc/apt/sources.list.d/kubernetes.list ] ; then
  touch /etc/apt/sources.list.d/kubernetes.list
  chmod 666 /etc/apt/sources.list.d/kubernetes.list
  tee /etc/apt/sources.list.d/kubernetes.list <<-'EOF'
deb http://mirrors.ustc.edu.cn/kubernetes/apt kubernetes-xenial main
EOF
  curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add -
  cp /etc/apt/trusted.gpg /etc/apt/trusted.gpg.d/
  # 重新更新仓库
  fi
  echo "k8s仓库配置... 完成"
  apt update
}

preinstall_debs(){
  apt update
  # 安装前置软件包
  apt install -y ssh sudo docker docker-compose curl gnupg zssh
  # 启动ssh
  systemctl enable  ssh && systemctl start ssh

  if [ ! -f cri-dockerd_0.3.10.3-0.debian-bookworm_amd64.deb ]; then
    curl -OL https://github.com/Mirantis/cri-dockerd/releases/download/v0.3.10/cri-dockerd_0.3.10.3-0.debian-bookworm_amd64.deb
  fi
  apt install -y ./cri-dockerd_0.3.10.3-0.debian-bookworm_amd64.deb

  sed -i 's|^ExecStart=/usr/bin/cri-dockerd --container-runtime-endpoint fd://|ExecStart=/usr/bin/cri-dockerd --pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.9 --container-runtime-endpoint fd://|' /lib/systemd/system/cri-docker.service


  cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "registry-mirrors": [ "https://1nj0zren.mirror.aliyuncs.com", "https://docker.mirrors.ustc.edu.cn", "http://f1361db2.m.daocloud.io", "https://registry.docker-cn.com" ]
}
EOF
  systemctl daemon-reload
  systemctl restart docker
  systemctl enable docker
  systemctl restart cri-docker.service
  echo "安装前置软件包，docker安装部署... 完成"
}

set_static_network(){
  # 配置网络
  sed -i "/^iface \([^ ]*\) inet dhcp$/s/dhcp/static\\
      address $IP_ADDR_CURRENT\\
      netmask 255.255.255.0\\
      gateway $GATEWAY/" /etc/network/interfaces
  echo "静态网络配置... 完成"
}

set_kernel_ipv4(){
  # 内核开启IPv4转发需要开启下面的模块
  modprobe br_netfilter

  cat > /etc/sysctl.d/k8s.conf <<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
  echo "内核开启IPv4转发... 完成"
}

check_root(){
  # 检查脚本是否以 root 身份运行
  if [[ $EUID -ne 0 ]]; then
     echo "错误：请以 root 身份运行此脚本"
     exit 1
  fi
}

install_k8s(){
  config_k8s_source
  # 安装k8s
  apt-get install -y kubelet kubernetes-cni kubeadm kubectl
  kubeadm version

  systemctl enable kubelet
  systemctl daemon-reload
  systemctl restart kubelet
  echo "k8s安装部署... 完成"
}

config_k8s(){
  mkdir k8s-init
  cd k8s-init || exit 0
  kubeadm config print init-defaults > kubeadm.conf
  # 替换 imageRepository 仓库地址
  sed -i 's/imageRepository:.*$/imageRepository: registry.aliyuncs.com\/google_containers/' kubeadm.conf
  # 替换 criSocket 的为cri-docker
  sed -i 's/criSocket:.*$/criSocket: unix:\/\/\/var\/run\/cri-dockerd.sock/' kubeadm.conf
  # 替换 advertiseAddress 的值
  sed -i "s/advertiseAddress:.*$/advertiseAddress: $IP_ADDR_MASTER/" kubeadm.conf
  # 替换 bindPort 的值
  sed -i 's/bindPort:.*$/bindPort: 6443/' kubeadm.conf
  sed -i 's/name: node/name: master/' kubeadm.conf
  kubeadm config images list --config kubeadm.conf
  echo "生成初始化配置... 完成"
  kubeadm config images pull --config kubeadm.conf
  kubeadm init --config kubeadm.conf
  echo "k8s环境初始化... 完成"
}

reset(){
  echo "y" | kubeadm reset --cri-socket unix:///var/run/cri-dockerd.sock
  rm -rf /etc/cni/net.d
  echo "k8s环境重置... 完成"
}

main(){
  # 关闭虚拟内存
  check_root
  swapoff -a
  set_hosts
  preinstall_debs
  set_static_network
  set_kernel_ipv4
  install_k8s
  config_k8s
  echo "配置安裝完成，请重启系统..."
}

case "$1" in
  --help)
    print_usage
    exit 0
    ;;
  reset)
    reset
    exit 0
    ;;
  master)
    IP_ADDR_CURRENT=$IP_ADDR_MASTER
    main
    ;;
  node)
    # 判断参数个数
    if [ $# -eq 2 ]; then
        if [ "$2" -eq 1 ]; then
          IP_ADDR_CURRENT=$IP_ADDR_NODE1
        elif [ "$2" -eq 2 ]; then
          IP_ADDR_CURRENT=$IP_ADDR_NODE1
        else
          echo "id not found: $1"
          exit 0
        fi
        main
    else
      echo "参数错误"
      print_usage
      exit 0
    fi
    ;;
esac