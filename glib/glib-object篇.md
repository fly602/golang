# Glib-Object

Glib是一个通用的C语言工具库，提供了一组用于编写高效、可移植、可扩展的C语言程序的功能。它包含了许多常见的数据结构，例如动态数组、链表、哈希表等，还提供了字符串处理、文件操作、内存管理等功能。Glib还具有跨平台的特性，可以在多种操作系统上运行。许多开源项目都基于glib进行开发，例如lightdm、network-manager等。

GObject是Glib库的一个组成部分，它是一个面向对象的编程框架。它提供了一套用于创建和管理对象、信号和属性的机制，并且支持对象的继承和多态。GObject使用Glib的基本功能，并扩展了它以支持对象的特性，例如封装、继承和多态。

-   对象模型：GObject基于一种称为"基于类的对象模型"的设计范式。每个对象都是基于一个类创建的，类定义了对象的属性和方法。对象可以从类继承属性和方法，并且可以重写或扩展它们。类使用C结构来表示，并且在运行时会有一个唯一的类对象来代表它。

-   类和对象：在GObject中，类和对象是紧密相关且相互依赖的概念。类定义了一组属性、方法和信号，而对象是类的一个实例，它具有类定义的属性和方法。类通常用于创建对象的蓝图，对象则是根据类定义实例化的。

-   属性：属性是GObject中最重要的概念之一。属性是对象的状态或特征，例如颜色、大小、可见性等等。每个属性都有一个名称、类型和访问方法。GObject提供了一些宏和函数，用于定义和使用属性。属性支持变化通知，当属性的值改变时，可以发送信号通知其他部分。

-   信号：信号是GObject中的一个核心概念，它用于实现对象间的通信和事件处理。信号是一种机制，当对象的某个特定事件发生时，发送信号通知其他对象。信号通常与回调函数关联，当信号被触发时，相应的回调函数会被调用。GObject提供了一些宏和函数，用于定义和处理信号。

-   继承和多态：GObject支持继承和多态，允许通过派生类扩展已有的类。这种机制可以减少代码重复，并提供更高级的抽象和封装。通过继承，子类可以重写或添加新的方法，从而实现特定的行为。

-   对象系统：GObject提供了一系列函数和宏，用于创建、初始化、操作和销毁对象。这些函数和宏可以创建对象实例，设置和获取属性值，连接信号和槽，管理对象的引用计数等。GObject还提供了一些辅助函数和工具，用于对象的类型检查和转换。

本文将通过C语言编写具体的代码示例对GObject的功能和使用方法进行具体的说明。

## 一、调试环境搭建：
### 1. 调试包的安装
glib的so库一般放在libglib2.0-0这个包中，如果想通过gdb调试glib库，需要安装相应的调试包：libglib2.0-0-dbgsym。gdb在跟踪的时候可以清除的看到库函数的调用堆栈以及参数变量的内容。
### 2. 编译环境
引用glib库的时候，编译时需要带上glib的库和头文件，可以通过pkg-config --cflags --libs查询，例如引用了glib和gobject：
```sh
uos@uos-PC:~$ pkg-config --cflags --libs glib-2.0 gobject-2.0
-I/usr/include/glib-2.0 -I/usr/lib/x86_64-linux-gnu/glib-2.0/include -lgobject-2.0 -lglib-2.0
```

[Makefile](./glib-main/Makefile)中可以这样使用:
```Makefile
...
CFLAGS = -g -Wall -Wpedantic -Wno-padded -O $(shell pkg-config --cflags --libs glib-2.0 gobject-2.0)
...
```

### 3. vscode调试环境搭建
1.  添加工作区：为了方便调试和代码跟踪，可以将项目可glib源码添加到同一工作区，代码分析的时候，vscode可以跳转到glib源码中。
2.  添加头文件路径：一般情况下，直接打开引用glib库的项目的时候，会提示错误，找不到glib的头文件，可以在.vscode/c_cpp_properties.json中添加头文件路径。例如：
    ```json
    {
        "configurations": [
            {
                ...
                "includePath": [
                    "${workspaceFolder}/**",
                    "/usr/include/glib-2.0",
                    "/usr/lib/x86_64-linux-gnu/glib-2.0/include"
                ],
                ...
            }
        ],
        "version": 4
    }
    ```
3.  源码关联：在.vscode/launch.json中添加如下代码，在vscode中通过gdb调试的时候，可以逐步跟踪跳转到glib的源码中:
    ```json
    {
        ...
        "configurations": [
            {
                ...
                // glib debug的环境变量，使用g_log时的日志输出
                "environment": [
                    {
                        "name": "G_MESSAGES_DEBUG",
                        "value": "all"
                    }
                ],
                ...
                // 关联glib源码，调试的时候可以跟踪跳转到glib库的源码中
                "sourceFileMap": {
                    "./debian/build/deb/../../../": "/home/uos/dde-go/src/github.com/linuxdeepin/glib2.0/"
                }
            },
        ]
    }
    ```

## 二、GObject的使用

### 1. GObject几个比较重要的宏
在了解GObject的使用之前，先来认识几个关键的宏：

1.  G_DEFINE_TYPE: G_DEFINE_TYPE 是 GLib 中用于简化 GObject 类型定义的宏。它是一个宏模板，定义了一个用于创建 GObject 类型的标准化流程，包括类型注册、类结构体定义和类初始化函数。

    函数原形：
    ```
    G_DEFINE_TYPE (Type, type_name, PARENT_TYPE)
    ```
    参数说明：

    -   Type：要定义的 GObject 类型的名称。
    -   type_name：类型名称的小写形式，用于内部结构体的命名。
    -   PARENT_TYPE：父类型的 GObject 类型的名称，通常是你自己定义的类的父类。
    -   G_DEFINE_TYPE 宏会自动生成以下内容：
        -   类型注册函数（Type_register_type()）：该函数在运行时注册类型以便使用。
        -   内部结构体定义：struct _Type，用于存储实例字段和函数指针。
        -   类初始化函数（Type_class_init()）：用于初始化类的虚函数表，并可以添加其他初始化逻辑。
        -   实例初始化函数（Type_init()）：用于初始化实例字段和其他实例特定的逻辑。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】：
    ```c
    ...
    // BASE_OBJ_TYPE是父类类型
    // 这个宏包含了DbusObj和DbusObjClass类型注册、类结构体定义和类初始化函数
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)

    // dbus_obj_init和dbus_obj_class_init需要自己实现，G_DEFINE_TYPE只是做了函数的申明
    static void dbus_obj_init (DbusObj *self)
    {
    }

    static void dbus_obj_class_init (DbusObjClass *klass)
    {
    }
    ...

    ```
    通过使用 G_DEFINE_TYPE 宏，可以避免手动编写一些繁琐的类型定义代码和函数，在一定程度上简化了 GObject 类型的定义过程。

2.  *type_name*_get_type: 用于获取一个 GObject 类型的 *type_name* 类型标识符（Type Identifier），返回的是GType，GType 是一个代表类型的数据类型。它用于表示在程序中定义的各种数据类型，如对象类、接口、枚举、结构体等。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】：
    ```c
    ...
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    // 获取DBUS_OBJ的类型，在NM项目中，一般是宏封装
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    GType dbus_obj_get_type (void);
    ...


    int main(){
        // g_type_name用于获取GType的名称
        g_log(domain, G_LOG_LEVEL_INFO, "g_type_name(DBUS_OBJ_TYPE) =%s!",g_type_name(DBUS_OBJ_TYPE));
    }
    ...
    ```
    运行输出结果：
    ```sh
    DBUS-OBJ-INFO: 16:59:07.002: dbus_obj.c:50: g_type_name(DBUS_OBJ_TYPE) =DbusObj!
    ```

3.  G_TYPE_FROM_INSTANCE: 通过对象的实例获取对象的类型。

    函数原形：
    ```c
    G_TYPE_FROM_INSTANCE(instance)
    ```
    参数说明：

    -   instance：要进行类型转换的对象实例。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】：
    ```c
    ...
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    // 获取DBUS_OBJ的类型
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    GType dbus_obj_get_type (void);
    ...


    int main(){
        // g_type_name用于获取GType的名称
        g_log(domain, G_LOG_LEVEL_INFO, "g_type_name(G_TYPE_FROM_INSTANCE(obj))=%s!",g_type_name(G_TYPE_FROM_INSTANCE(obj)));
    }
    ...
    ```

    运行输出结果：
    ```sh
    DBUS-OBJ-INFO: 17:09:44.275: dbus_obj.c:51: g_type_name(G_TYPE_FROM_INSTANCE(obj))=DbusObj!
    ```

4.  G_TYPE_CHECK_INSTANCE_CAST：GLib 开发中用于类型转换的宏定义之一。它用于将一个 GObject 对象转换为其派生类型的对象。

    函数原形：
    ```c
    G_TYPE_CHECK_INSTANCE_CAST(instance, g_type, c_type)
    ```
    参数说明：

    -   instance：要进行类型转换的对象实例。
    -   g_type：目标类型的 GObject 类型。
    -   c_type：目标类型的 C 类型。

    代码示例【[完整代码示例](./glib-main/dbus_obj.h)】：
    ```c
    // base_obj.h
    ...
    #define BASE_OBJ_TYPE (base_obj_get_type()) 
    // 在使用时一般用宏进行再次封装
    #define BASE_OBJ(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), BASE_OBJ_TYPE, BaseObj))
    ...

    // dbus_obj.c
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    int main(){
        DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);

        BaseObj *base = BASE_OBJ(obj);
        g_log(domain, G_LOG_LEVEL_INFO, "obj=0x%p base=0x%p!",obj,&obj->parent,base);
    }
    ...

    ```
    运行输出结果：
    ```sh
    // 地址都一样，其实就是类型转换
    DBUS-OBJ-INFO: 17:28:51.255: dbus_obj.c:54: obj=0x0x411ac0 base=0x0x411ac0!
    ```
5.  G_TYPE_INSTANCE_GET_CLASS: 用于获取一个 GObject 实例的类结构体的指针。

    函数原形：
    ```c
    G_TYPE_INSTANCE_GET_CLASS(instance, g_type, c_type)
    ```
    参数说明：
    -   instance：需要检查的 GObjectClass 结构体。
    -   g_type：目标类的 GType。
    -   c_type：目标类的结构体类型。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】：
    ```c
        ...
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    // 获取DBUS_OBJ的类型，在NM项目中，一般是宏封装
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    #define DBUS_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  DBUS_OBJ_TYPE, DbusObjClass))
    GType dbus_obj_get_type (void);
    ...


    int main(){
        DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);
        DbusObjClass *class = DBUS_OBJ_GET_CLASS(obj);
    }
    ...
    ```


6.  G_TYPE_CHECK_INSTANCE_TYPE: GLib 开发中用于检查类型的宏定义之一。用于检查一个 GObject 实例的类型是否匹配给定的类型。返回true或者false。

    函数原形：
    ```
    G_TYPE_CHECK_INSTANCE_TYPE(instance, g_type) 
    ```
    参数说明：
    -   instance：需要检查类型的 GObject 实例。
    -   g_type：目标类的 GType。

    代码示例【[完整代码示例](./glib-main/main.c)】：
    ```c
    // base_obj.h
    ...
    #define BASE_OBJ_TYPE (base_obj_get_type()) 
    // 在使用时一般用宏进行再次封装，检查并转换klass为BaseObjClass
    #define BASE_OBJ_CLASS(klass)     (G_TYPE_CHECK_CLASS_CAST ((klass),  BASE_OBJ_TYPE, BaseObjClass))
    ...

    // dbus_obj.c
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    #define DBUS_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  DBUS_OBJ_TYPE, DbusObjClass))

    // main.c
    int main(){
        DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);

        BaseObjClass *base_class = BASE_OBJ_CLASS(DBUS_OBJ_GET_CLASS(obj));
        ...
    }
    ...

    ```
7.  G_TYPE_CHECK_CLASS_TYPE: GLib 开发中用于检查类型的宏定义之一。用于检查一个 GObject 类的类型是否匹配给定的类型。返回true或者false。

    函数原形：
    ```
    G_TYPE_CHECK_CLASS_TYPE(g_class, g_type) 
    ```
    参数说明：
    -   g_class：需要检查类型的 GObject 类的类结构体指针。
    -   g_type：目标类的 GType。

    代码示例【[完整代码示例](./glib-main/main.c)】：

8.  G_TYPE_CHECK_CLASS_CAST:  GLib 开发中用于类型转换的宏定义之一。这个宏用于检查一个给定的 GObjectClass 结构体是否与指定的 GType 类型相匹配，并进行相应的类型转换。这个宏通常用于在 GObject 的类型系统中进行类型安全的操作。

    函数原形：
    ```
    G_TYPE_CHECK_CLASS_CAST(g_class, g_type, c_type)
    ```
    参数说明：
    -   g_class：需要检查的 GObjectClass 结构体。
    -   g_type：目标类的 GType。
    -   c_type：目标类的结构体类型。

    代码示例【[完整代码示例](./glib-main/main.c)】：
    ```c
    // base_obj.h
    ...
    #define BASE_OBJ_TYPE (base_obj_get_type()) 
    // 在使用时一般用宏进行再次封装，检查并转换klass为BaseObjClass
    #define BASE_OBJ_CLASS(klass)     (G_TYPE_CHECK_CLASS_CAST ((klass),  BASE_OBJ_TYPE, BaseObjClass))
    ...

    // dbus_obj.c
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    #define DBUS_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  DBUS_OBJ_TYPE, DbusObjClass))

    // main.c
    int main(){
        DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);

        BaseObjClass *base_class = BASE_OBJ_CLASS(DBUS_OBJ_GET_CLASS(obj));
        ...
    }
    ...

    ```
