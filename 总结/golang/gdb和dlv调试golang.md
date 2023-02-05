# GDB和dlv调试

## 1. 摘要

本章节讲解gdb和dlv调试GOLANG程序的入门配置，以及gdb和dlv命令详解备忘。

## 2. gdb调试go程序入门

gdb是linux系统自带的调试器，功能十分强大，它不仅支持C/C++调试，也支持GO程序调试。

### 2.1 配置gdb

(1) 打开gdb初始化配置文件：

vim ~/.gdbinit
(2) 增加一行，:wq!保存后退出：

add-auto-load-safe-path /usr/local/go/src/runtime/runtime-gdb.py

### 2.2 编译golang

虽然gdb也支持golang了，但是在编译golang仍然需要加一些特殊的参数，不然打印变量会提示找不到：No symbol in current context

编译添加参数： -gcflags "-N -l"

## 3.gdb相关命令

- gdb -tui test打开调试程序，界面分页，上面是代码，下面是命令；
- gdbtui的开关快捷键：ctrl+x ctrl+a或者ctrl+x A
- file test在运行gdb下打开某个文件
- run/r 运行
- continue/c 继续运行
- step/s 如果有函数则进入函数执行
- finish 跳出当前的函数
- stop 停止运行
- until xxx 可用于跳出循环
- guit/ctrl+d 退出GDB
- print/p var 打印变量的值
- print/p &var 打印变量地址
- printf/p *addr 打印地址的值
- printf/p /x var 用16进制显示数据 x十六进制/d十进制/u十六进制无符号/t二进制/c字符/f浮点
- break/b xxx 在某行打断点
- break/b fun 在某个函数处加断点
- break/b 30 if n==100 //当变量n等于100的时候在30行处加断点
- break fileName:N 在某个文件的N行加断点
- info break/b 查看断点
- clear N 删除N行断点
- delete N 删除N号断点
- delete 删除所有断点
- disable xxx 失能断点
- enable xxx 使能断点
- info b 查看断点
- info source 查看当前程序
- info stack 查看堆栈信息
- info args 查看当前参数值
- display args 查看当前参数值
- bt 查看函数堆栈
- pwd查看程序路径
- ctrl+p 前一条命令
- ctrl+n 下一条命令
- watch xxx 设置监控点，在变量改变的时候停下来。(不可直接设置，先加断点在监测)
- ctrl+l可能layout会造成控制台花屏,使用ctrl+L清屏
- list linenum：以linenum指定的行号为中心，显示10行
- list function：以指定的函数为中心，显示10行
- list：重复上一次的list指令，也可以直接按回车键，重复上次指令。
- set listsize count：设置每次显示的行数。
- show listsize：显示已设置的显示行数。
- list first,last：显示指定起始行到结束结束行的源文件。
- list ,last：显示以指定的last为结束行，显示10行。
- list first,：以first为第一行，显示10行。
- list +：以上次显示的结束行为起始行显示后10行
- list –：以上次显示的起始行为结束行，显示前10行

## 4.dlv安装

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

go版本小于1.16的用下面方式安装

```bash
 git clone https://github.com/go-delve/delve
 cd delve
 go install github.com/go-delve/delve/cmd/dlv
```

## 5.dlv指令

仅列出常用或者会用到的

|指令 |用处|
|----|:---|
|attach |这个命令将使Delve控制一个已经运行的进程，并开始一个新的调试会话。 当退出调试会话时，你可以选择让该进程继续运行或杀死它。|
|exec |这个命令将使Delve执行二进制文件，并立即附加到它，开始一个新的调试会话。请注意，如果二进制文件在编译时没有关闭优化功能，可能很难正确地调试它。请考虑在Go 1.10或更高版本上用-gcflags="all=-N -l "编译调试二进制文件，在Go的早期版本上用-gcflags="-N -l"。|
|help |使用手册|
|debug |默认情况下，没有参数，Delve将编译当前目录下的 "main "包，并开始调试。或者，你可以指定一个包的名字，Delve将编译该包，并开始一个新的调试会话。|
|test |test命令允许你在单元测试的背景下开始一个新的调试会话。默认情况下，Delve将调试当前目录下的测试。另外，你可以指定一个包的名称，Delve将在该包中调试测试。双破折号`--`可以用来传递参数给测试程序。|
|version |查看dlv版本|

## 6.dlv调试指令

仅记录个人觉得会用到的指令

### 6.1断点管理

|指令 |缩写| 用法|
|----|:---|:---|
|break |b| 设置断点|
|breakpoints |bp| 查看当前所有断点|
|clear |/| 删除断点|
|clearall |/| 删除多个断点|
|toggle |/| 启用或关闭断点|

### 6.2程序执行中的调试指令

|指令| 缩写| 用法|
|---|:----|:---|
|continue |c| 继续执行到一个断点或者程序结束吗|
|next |n| 执行下一行代码|
|restart |r| 重新执行程序|
|step |s| 执行代码的下一步|
|step-instruction |si| 执行下一行机器码|
|stepout |so| 跳出当前执行函数|

### 6.3参数管理

|指令 |缩写| 用法|
|----|:---|----|
|args |/| 打印函数input|
|display |/| 打印加入到display的变量的值，每次执行下一行代码或下一个断点时|
|locals |/| 打印局部变量|
|print |p| 打印表达式的结果|
|set |/| 设置某个变量的值|
|vars |/| 查看全局变量|
|whatis |/| 查看变量类型|

## 6.4其他

|指令 |缩写| 用法|
|-------|:--|-----|
|disassemble |disass|查看反编译后的代码，机器码|
|exit| quit / q| 退出|
|funcs |/| 打印程序用到的所有函数|
|help |h| 帮助信息|
list |ls / l |打印代码|
