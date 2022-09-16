#!/bin/bash

## 仓库地址,根据对应仓库修改
ORIGIN=https://aptly.uniontech.com/pkg/kelvinu-sp2/release-candidate/a2x1LXNqcy1rd2luMjAyMi0wOS0xNSAxNzoyMDowNQ
## 包名，根据需要自行补充
PKGNAMES=(kwin kwayland dde-kwin)


## 仓库base路径
BASE_PATH=pool/main
## wget
WGET=wget
## options
OPTIONS='--mirror -r -np -R "*index.html*" -nd -e robots=off'

DIR_LOCAL="debs-all"
DIR_INSTALL="debs-install"

echo "仓库源地址: "$ORIGIN
echo "仓库路径："${PKGNAMES[*]}
echo "option=" $OPTIONS

if [ ! -d $DIR_LOCAL ];then
mkdir $DIR_LOCAL -v
fi

if [ ! -d $DIR_INSTALL ];then
mkdir $DIR_INSTALL -v
fi

cd $DIR_LOCAL

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

##  将需要的包压缩
cp *.deb ../$DIR_INSTALL
cd ../$DIR_INSTALL
ls *-dev* *-dbgsym* *-doc* |xargs rm -rf
cd ..
tar -cvzf deb-install.tar.gz $DIR_INSTALL
#