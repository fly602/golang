#   Network-Manager编译调试
##  一、编译环境：
*   sudo apt-get build-dep network-manager
*   sudo apt-get install libnss3-dev

##  二、构建：
### 1.  使用make
####    配置编译安装路径：
*   ./configure --prefix=/ --bindir=/usr/sbin
####    添加日志：
*   修改NetworkManager.service.in中ExecStart 添加参数--log-level参数，可设置为：off,err,warn,info,debug,trace
####    编译安装：
*   make;sudo make install

### 2.  打包
####    dpkg-buildpackage -b -us -uc -j8

##  三、Vscode gdb调试
### 1.  环境配置：
nm是在root环境下运行的，所以使用gdb调试需要sudo执行，而执行sudo需要密码。那就需要设置免密：
在/ect/sudoers 中添加:
```
<用户名>    ALL=(ALL:ALL) NOPASSWD:ALL
```
### 2.  在.vscode下编写gdb提权脚本替换gdb
```
#!/bin/bash
sudo /usr/bin/gdb "$@"
```
### 3.  编写vscode调试文件：launch.json
```
{
    // 使用 IntelliSense 了解相关属性。 
    // 悬停以查看现有属性的描述。
    // 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "(gdb) nm attach",
            "type": "cppdbg",
            "request": "attach",
            "program": "/usr/sbin/NetworkManager",
            "processId": "${command:pickProcess}",
            "MIMode": "gdb",
            // 使用自己编写提权的gdb进行调试
            "miDebuggerPath":"${workspaceFolder}/.vscode/gdb",
            "setupCommands": [
                {
                    "description": "为 gdb 启用整齐打印",
                    "text": "-enable-pretty-printing",
                    "ignoreFailures": true
                },
                {
                    "description":  "将反汇编风格设置为 Intel",
                    "text": "-gdb-set disassembly-flavor intel",
                    "ignoreFailures": true
                }
            ]
        },
        {
            "name": "(gdb) nm启动",
            "type": "cppdbg",
            "request": "launch",
            "program": "/usr/sbin/NetworkManager",
            // 添加运行参数
            "args": ["--log-level=TRACE", "--no-daemon"],
            "stopAtEntry": false,
            "cwd": "${fileDirname}",
            "environment": [],
            "externalConsole": false,
            "MIMode": "gdb",
            // 使用自己编写提权的gdb进行调试
            "miDebuggerPath":"${workspaceFolder}/.vscode/gdb",
            "setupCommands": [
                {
                    "description": "为 gdb 启用整齐打印",
                    "text": "-enable-pretty-printing",
                    "ignoreFailures": true
                },
                {
                    "description":  "将反汇编风格设置为 Intel",
                    "text": "-gdb-set disassembly-flavor intel",
                    "ignoreFailures": true
                }
            ]
        }

    ]
}
```