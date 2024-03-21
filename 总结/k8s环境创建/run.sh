#!/bin/bash

#set -x
# 默认参数值
IP_ADDR_MASTER="192.168.3.180"
IP_ADDR_NODE1="192.168.3.181"
IP_ADDR_NODE2="192.168.3.182"
GATEWAY="192.168.3.1"
K8S_LOCAL_NODE="master"

print_usage() {
    echo "Usage: $0 [master | node <id> reset] [--help]"
}

# 判断参数个数
if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

create_env(){
  cat > /etc/k8s-env.sh << EOF
export K8S_GATEWAY=$GATEWAY
export USER_NAME="uos"
export K8S_LOCAL_NODE=$K8S_LOCAL_NODE
export IP_ADDR_MASTER=$IP_ADDR_MASTER
export IP_ADDR_NODE1=$IP_ADDR_NODE1
export IP_ADDR_NODE2=$IP_ADDR_NODE2
EOF
echo "生成k8s系统环境变量..."
}

print_env(){
  echo "当前用户名: $USER_NAME"
  USER_HOME_DIR=/home/"$USER_NAME"/
  echo "用户目录:        $USER_HOME_DIR"
  echo "当前k8s节点:     $K8S_LOCAL_NODE"
  if [ "$K8S_LOCAL_NODE" == "master" ]; then
    echo "当前节点IP地址:   $IP_ADDR_MASTER"
  elif [ "$K8S_LOCAL_NODE" == "node1" ]; then
    echo "当前节点IP地址:   $IP_ADDR_IP_ADDR_NODE1"
  elif [ "$K8S_LOCAL_NODE" == "node2" ]; then
    echo "当前节点IP地址:   $IP_ADDR_IP_ADDR_NODE2"
  else
    echo "未知节点，请修改/etc/k8s-env.sh."
    exit 0
  fi
}

init_env(){
  if [ ! -e /etc/k8s-env.sh ]; then
    touch /etc/k8s-env.sh
    create_env
  fi
  source /etc/k8s-env.sh
  USER_HOME_DIR=/home/"$USER_NAME"/
  print_env
  echo "请确认是否使用该配置环境，如需修改，请按'ctrl c'终止脚本，并修改/etc/k8s-env.sh."
  for ((i=5; i>=0; i--)); do
      echo -ne "脚本将在'$i'秒后自动配置...\r"
      sleep 1  # 等待1秒
  done
}

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
  apt install -y ssh sudo docker docker-compose curl gnupg zssh wget
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
  echo "k8s安装部署... 完成"
}

install_fannel(){
  if [ -e "kube-flannel.yml" ];then
    wget -O kube-flannel.yml https://github.com/flannel-io/flannel/releases/download/v0.24.3/kube-flannel.yml
  fi
  kubectl apply -f ./kube-flannel.yml
}

init_k8s(){
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
  if [ -d k8s-init ]; then
    rm -rf k8s-init/
  fi
  if [ -d .kube/config ]; then
    rm .kube/config
  fi
  echo "k8s环境重置... 完成"
}

config_k8s(){
  mkdir -p "$USER_HOME_DIR"/.kube
  cp -i /etc/kubernetes/admin.conf "$USER_HOME_DIR"/.kube/config
  chown "$(id -u "$USER_NAME")":"$(id -g "$USER_NAME")" "$USER_HOME_DIR"/.kube/config
  echo "export KUBECONFIG=\$USER_HOME_DIR/admin.conf" >> /etc/profile
  systemctl daemon-reload
  systemctl restart kubelet
}

check_args(){
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
      K8S_LOCAL_NODE="$1"
      IP_ADDR_CURRENT=$IP_ADDR_MASTER
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
          K8S_LOCAL_NODE="$1""$2"
          echo "===>>>node: $K8S_LOCAL_NODE"
      else
        echo "参数错误"
        print_usage
        exit 0
      fi
      ;;
  esac
}

main(){
  # 关闭虚拟内存
  check_root
  check_args "$@"
  init_env
  swapoff -a
  set_hosts
  preinstall_debs
  set_static_network
  set_kernel_ipv4
  install_k8s
  init_k8s
  config_k8s
  install_fannel
  echo "配置安裝完成，请重启系统..."
}

main "$@"
