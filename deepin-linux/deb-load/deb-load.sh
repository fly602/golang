#!/bin/bash

########################## 使用说明 ##############################
# 1、下载指定包：需要配置下面的仓库“ORIGIN”参数和包名"PKGNAMES" 参数即可
# 例如：
# ORIGIN=https://aptly.uniontech.com/pkg/kelvinu-sp2/release-candidate/a2x1LXNqcy1rd2luMjAyMi0wOS0xNSAxNzoyMDowNQ
# PKGNAMES=(kwin kwayland)
#
# 2、下载仓库全量包：需配置下面的仓库“ORIGIN”参数即可
# 例如：
# ORIGIN=https://aptly.uniontech.com/pkg/kelvinu-sp2/release-candidate/a2x1LXNqcy1rd2luMjAyMi0wOS0xNSAxNzoyMDowNQ
#
# 3、生成文件说明：
# debs-all：下载全量包
# debs-install：交付补丁包，去掉了非deb文件以及调试和dev编译依赖包
#
#################################################################

######################### 手动修改仓库和包 #########################
## 仓库地址,根据对应仓库修改
ORIGIN=https://aptly.uniontech.com/pkg/kelvinu-sp2/release-candidate/a2x1LXNqcy1rd2luMjAyMi0wOS0xNSAxNzoyMDowNQ
## 包名，根据需要自行补充
PKGNAMES=()
#################################################################

######################### 以下配置保持不变 #########################
## 仓库base路径
BASE_PATH=pool/main
## wget
WGET=wget
## options
OPTIONS='--mirror -r -np -R "*index.html*" -nd -e robots=off'
DIR_LOCAL="debs-all"
DIR_INSTALL="debs-install"
#################################################################

echo "仓库源地址: "$ORIGIN
echo "仓库路径："${PKGNAMES[*]}

if [ ! $ORIGIN];then
echo "脚本仓库源地址为空，请配置脚本仓库源地址。。。"
fi

if [ ! -d $DIR_LOCAL ];then
mkdir $DIR_LOCAL -v
fi

if [ ! -d $DIR_INSTALL ];then
mkdir $DIR_INSTALL -v
fi

cd $DIR_LOCAL
if [ ! $PKGNAMES];then
    ##  如果数组为空， 默认下载所有deb包
    $WGET $OPTIONS   $ORIGIN/
else
    for pkg in ${PKGNAMES[*]}
    do
        ## 检测是否是lib开头
        if [ ${pkg:0:3} = "lib" ];then
            pkg_path=$BASE_PATH/${pkg:0:4}/$pkg
        else
            pkg_path=$BASE_PATH/${pkg:0:1}/$pkg
        fi
        $WGET $OPTIONS   $ORIGIN/$pkg_path    
    done
fi

##  将需要的包压缩
cp *.deb ../$DIR_INSTALL
cd ../$DIR_INSTALL
ls *-dev* *-dbgsym* *-doc* |xargs rm -rf
cd ..
tar -cvzf deb-install.tar.gz $DIR_INSTALL
#