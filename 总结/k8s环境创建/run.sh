#!/bin/bash

#set -x
# 默认参数值
IP_ADDR_MASTER="192.168.3.180"
IP_ADDR_NODE1="192.168.3.181"
IP_ADDR_NODE2="192.168.3.182"
K8S_GATEWAY="192.168.3.1"
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
export K8S_GATEWAY=$K8S_GATEWAY
export USER_NAME="uos"
export IP_ADDR_MASTER=$IP_ADDR_MASTER
export IP_ADDR_NODE1=$IP_ADDR_NODE1
export IP_ADDR_NODE2=$IP_ADDR_NODE2
EOF
echo "生成k8s系统环境变量..."
}

print_env(){
  echo "当前用户名：        $USER_NAME"
  USER_HOME_DIR=/home/"$USER_NAME"
  echo "用户目录：         $USER_HOME_DIR"
  echo "主节点的IP:        $IP_ADDR_MASTER"
  echo "节点1的IP:         $IP_ADDR_NODE1"
  echo "节点2的IP:         $IP_ADDR_NODE2"
  echo "当前IP掩码:        $K8S_GATEWAY"
}

init_env(){
  if [ ! -e /etc/k8s-env.sh ]; then
    touch /etc/k8s-env.sh
    create_env
  fi
  source /etc/k8s-env.sh
  USER_HOME_DIR=/home/"$USER_NAME"
  print_env
  echo "请确认是否使用该配置环境，如需修改，请按'ctrl c'终止脚本，并修改/etc/k8s-env.sh."
  for ((i=5; i>=0; i--)); do
      echo -ne "脚本将在'$i'秒后自动配置...\r"
      sleep 1  # 等待1秒
  done
}

set_hosts(){
host_file=/etc/hosts
sed -i '/master$/d' "$host_file"
sed -i '/node1$/d' "$host_file"
sed -i '/node2$/d' "$host_file"
hosts=("127.0.1.1 $K8S_LOCAL_NODE" "$IP_ADDR_MASTER master" "$IP_ADDR_NODE1 node1" "$IP_ADDR_NODE2 node2")
for item in "${hosts[@]}"; do
    # 检查是否已经存在相同的记录
    if ! grep -q "$item" "$host_file"; then
        # 追加内容到 /etc/hosts 文件
        echo "$item" | sudo tee -a "$host_file"
    fi
done
echo "设置hosts... 完成"
}

set_hostname(){
# 获取当前主机名
current_hostname=$(hostname)

# 检查当前主机名是否为 "master"
if [ "$current_hostname" != "$K8S_LOCAL_NODE" ]; then
    echo "当前主机名为 $current_hostname，不是 $K8S_LOCAL_NODE，将修改主机名和 hosts 文件..."

    # 修改主机名
    sudo hostnamectl set-hostname "$K8S_LOCAL_NODE"

    # 修改 hosts 文件
    sudo sed -i "s/$current_hostname/$K8S_LOCAL_NODE/g" /etc/hosts

    echo "主机名和 hosts 文件已修改为 $K8S_LOCAL_NODE"
else
    echo "当前主机名已经是 $K8S_LOCAL_NODE，无需修改。"
fi
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

do_swapoff() {
  swapoff -a
  # 永久关闭swapoff
  sed -i '/swap/{/^#/!s/^/#/}' /etc/fstab
}

set_static_network(){
# 网络配置文件路径
interfaces_file="/etc/network/interfaces"

# 检查网络配置文件是否存在
if [ -f "$interfaces_file" ]; then
    # 检查是否为静态 IP
    if grep -Eq "iface [a-zA-Z0-9]+ inet dhcp" "$interfaces_file"; then
        echo "当前配置为 DHCP，将修改为静态 IP..."

        # 修改配置为静态 IP
        sudo sed -Ei 's/iface ([a-zA-Z0-9]+) inet dhcp/iface \1 inet static/g' "$interfaces_file"
    elif grep -Eq "iface [a-zA-Z0-9]+ inet static" "$interfaces_file"; then
        echo "当前配置为静态 IP，将修改 IP 地址"
    else
        echo "当前网络配置不是静态 IP 也不是 DHCP，尝试修复。"
        echo "allow-hotplug ens33" | tee -a "$interfaces_file" >/dev/null
        echo "iface ens33 inet static" | tee -a "$interfaces_file" >/dev/null
    fi
    sed -i '/^address/d' "$interfaces_file"
    sed -i '/^netmask/d' "$interfaces_file"
    sed -i '/^gateway/d' "$interfaces_file"
    echo "address $IP_ADDR_CURRENT" | tee -a "$interfaces_file" >/dev/null
    echo "netmask 255.255.255.0" | tee -a "$interfaces_file" >/dev/null
    echo "gateway $K8S_GATEWAY" | tee -a "$interfaces_file" >/dev/null
else
    echo "网络配置文件 $interfaces_file 不存在。"
    exit 0
fi

echo "网络配置如下："
cat /etc/network/interfaces
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
  if [ ! -e "kube-flannel.yml" ];then
    wget -O kube-flannel.yml https://github.com/flannel-io/flannel/releases/download/v0.24.3/kube-flannel.yml
  fi
  sudo -u $USER_NAME kubectl apply -f ./kube-flannel.yml
  echo "k8s网络插件[fannel]安装... 完成"
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
  sed -i "s/name: node$/name: $K8S_LOCAL_NODE/" kubeadm.conf
  kubeadm config images list --config kubeadm.conf
  echo "生成初始化配置... 完成"
  kubeadm config images pull --config kubeadm.conf
  kubeadm certs renew all --config=kubeadm.conf
  kubeadm init --config kubeadm.conf
  echo "k8s环境初始化... 完成"
}

reset(){
  if [ ! -e /etc/k8s-env.sh ]; then
      echo "/etc/k8s-env.sh不存在，环境可能已经初始化"
      exit 0
  else
    source /etc/k8s-env.sh
    USER_HOME_DIR=/home/"$USER_NAME"
  fi

  init_dir_path="$USER_HOME_DIR"/k8s-init
  config_dir_path="$USER_HOME_DIR"/.kube/
  pki_dir_path="/var/lib/kubelet/pki/"

  echo "y" | kubeadm reset --cri-socket unix:///var/run/cri-dockerd.sock
  rm -rf /etc/cni/net.d
  if [ -e "$init_dir_path" ]; then
    echo "删除$init_dir_path"
    rm -rf "$init_dir_path"
  fi
  if [ -e "$config_dir_path" ]; then
    echo "删除$config_dir_path"
    rm -rf "$config_dir_path"
  fi

  if [ -e "$pki_dir_path" ]; then
    rm -rf "$pki_dir_path"
  fi
  echo "k8s环境重置... 完成"
}

config_k8s(){
  mkdir -p "$USER_HOME_DIR"/.kube
  cp -i /etc/kubernetes/admin.conf "$USER_HOME_DIR"/.kube/config
  chown "$(id -u "$USER_NAME")":"$(id -g "$USER_NAME")" "$USER_HOME_DIR"/.kube/config
  echo "export KUBECONFIG=\$HOME/.kube/config" >> /etc/profile

  if ! grep -q "export KUBECONFIG=\$HOME/.kube/config" /etc/profile; then
          # 追加内容到 /etc/hosts 文件
          echo "export KUBECONFIG=\$HOME/.kube/config" >> /etc/profile
  fi
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
      init_env
      K8S_LOCAL_NODE="$1"
      IP_ADDR_CURRENT=$IP_ADDR_MASTER
      ;;
    node)
      init_env
      # 判断参数个数
      if [ $# -eq 2 ]; then
          if [ "$2" -eq 1 ]; then
            IP_ADDR_CURRENT=$IP_ADDR_NODE1
          elif [ "$2" -eq 2 ]; then
            IP_ADDR_CURRENT=$IP_ADDR_NODE2
          else
            echo "id not found: $1"
            exit 0
          fi
          K8S_LOCAL_NODE="$1""$2"
      else
        echo "参数错误"
        print_usage
        exit 0
      fi
      ;;
  esac
  echo "当前节点名称：     $K8S_LOCAL_NODE"
  echo "当前节点地址：     $IP_ADDR_CURRENT"
}

main(){
  # 关闭虚拟内存
  check_root
  check_args "$@"
  do_swapoff
  set_hostname
  set_hosts
  preinstall_debs
  set_static_network
  set_kernel_ipv4
  install_k8s
  if [ "$K8S_LOCAL_NODE" == "master" ];then
  init_k8s
  config_k8s
  install_fannel
  fi
  echo "配置安裝完成，请重启系统..."
}

main "$@"
