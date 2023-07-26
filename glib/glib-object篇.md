# Glib-Object

Glib是一个通用的C语言工具库，提供了一组用于编写高效、可移植、可扩展的C语言程序的功能。它包含了许多常见的数据结构，例如动态数组、链表、哈希表等，还提供了字符串处理、文件操作、内存管理等功能。Glib还具有跨平台的特性，可以在多种操作系统上运行。许多开源项目都基于glib进行开发，例如lightdm、network-manager等。

GObject是Glib库的一个组成部分，它是一个面向对象的编程框架。它提供了一套用于创建和管理对象、信号和属性的机制，并且支持对象的继承和多态。GObject使用Glib的基本功能，并扩展了它以支持对象的特性，例如封装、继承和多态。

-   **对象模型**: GObject基于一种称为"基于类的对象模型"的设计范式。每个对象都是基于一个类创建的，类定义了对象的属性和方法。对象可以从类继承属性和方法，并且可以重写或扩展它们。类使用C结构来表示，并且在运行时会有一个唯一的类对象来代表它。

-   **类和对象**: 在GObject中，类和对象是紧密相关且相互依赖的概念。类定义了一组属性、方法和信号，而对象是类的一个实例，它具有类定义的属性和方法。类通常用于创建对象的蓝图，对象则是根据类定义实例化的。

-   **属性**: 属性是GObject中最重要的概念之一。属性是对象的状态或特征，例如颜色、大小、可见性等等。每个属性都有一个名称、类型和访问方法。GObject提供了一些宏和函数，用于定义和使用属性。属性支持变化通知，当属性的值改变时，可以发送信号通知其他部分。

-   **信号**: 信号是GObject中的一个核心概念，它用于实现对象间的通信和事件处理。信号是一种机制，当对象的某个特定事件发生时，发送信号通知其他对象。信号通常与回调函数关联，当信号被触发时，相应的回调函数会被调用。GObject提供了一些宏和函数，用于定义和处理信号。

-   **继承和多态**: GObject支持继承和多态，允许通过派生类扩展已有的类。这种机制可以减少代码重复，并提供更高级的抽象和封装。通过继承，子类可以重写或添加新的方法，从而实现特定的行为。

-   **对象系统**: GObject提供了一系列函数和宏，用于创建、初始化、操作和销毁对象。这些函数和宏可以创建对象实例，设置和获取属性值，连接信号和槽，管理对象的引用计数等。GObject还提供了一些辅助函数和工具，用于对象的类型检查和转换。

本文将通过C语言编写具体的代码示例对GObject的功能和使用方法进行具体的说明。

## 一、调试环境搭建: 
### 1. 调试包的安装
glib的so库一般放在libglib2.0-0这个包中，如果想通过gdb调试glib库，需要安装相应的调试包: libglib2.0-0-dbgsym。gdb在跟踪的时候可以清除的看到库函数的调用堆栈以及参数变量的内容。
### 2. 编译环境
引用glib库的时候，编译时需要带上glib的库和头文件，可以通过pkg-config --cflags --libs查询，例如引用了glib和gobject: 
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
1.  **添加工作区**: 为了方便调试和代码跟踪，可以将项目可glib源码添加到同一工作区，代码分析的时候，vscode可以跳转到glib源码中。
2.  **添加头文件路径**: 一般情况下，直接打开引用glib库的项目的时候，会提示错误，找不到glib的头文件，可以在.vscode/c_cpp_properties.json中添加头文件路径。例如: 
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
3.  **源码关联**: 在.vscode/launch.json中添加如下代码，在vscode中通过gdb调试的时候，可以逐步跟踪跳转到glib的源码中:
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

### 1. 对象模型和类和对象

GObject基于一种称为"基于类的对象模型"的设计范式。每个对象都是基于一个类创建的，类定义了对象的属性和方法。对象可以从类继承属性和方法，并且可以重写或扩展它们。类和对象是紧密相关且相互依赖的概念。类定义了一组属性、方法和信号，而对象是类的一个实例，它具有类定义的属性和方法。

GObject基本类和基础对象是GObjectClass和GObject，新的类型和对象大多数都是基于这个类和对象派生出来的。

1.  GObjectClass结构体的介绍

    在了解GObject是怎么运作之前，还要认识一下GObjectClass，它的结构体原形是struct _GObjectClass，它是GObject的类结构体，也是GObject继承的基础。_GObjectClass 结构体中包含了一组函数指针，这些函数定义了在实例化 GObject 类型对象时的行为。这些函数指针包括构造函数 (constructor)、析构函数 (destructor)、对象属性的设置和获取函数、信号处理函数等。

    ```c
    struct  _GObjectClass
    {
    GTypeClass   g_type_class;

    /*< private >*/
    GSList      *construct_properties;

    /*< public >*/
    /* seldom overidden */
    GObject*   (*constructor)     (GType                  type,
                                    guint                  n_construct_properties,
                                    GObjectConstructParam *construct_properties);
    /* overridable methods */
    void       (*set_property)		(GObject        *object,
                                            guint           property_id,
                                            const GValue   *value,
                                            GParamSpec     *pspec);
    void       (*get_property)		(GObject        *object,
                                            guint           property_id,
                                            GValue         *value,
                                            GParamSpec     *pspec);
    void       (*dispose)			(GObject        *object);
    void       (*finalize)		(GObject        *object);
    /* seldom overidden */
    void       (*dispatch_properties_changed) (GObject      *object,
                            guint	   n_pspecs,
                            GParamSpec  **pspecs);
    /* signals */
    void	     (*notify)			(GObject	*object,
                        GParamSpec	*pspec);

    /* called when done constructing */
    void	     (*constructed)		(GObject	*object);

    /*< private >*/
    gsize		flags;

    /* padding */
    gpointer	pdummy[6];
    };
    ```

    _GObjectClass结构体成员: 

    1.  **g_type_class**:
        -   类型: GTypeClass
        -   描述: 用于存储 GObject 类的类型信息。GObject 类是 GObject 类型系统中的基本类型，g_type_class 可以用于访问或操作该类的类型信息。

    2.  **construct_properties**:
        -   类型: GSList *
        -   描述: : 这是一个指向 GSList 的指针，GSList 是一个单链表结构，用于存储构造属性（construct properties）。构造属性是在创建对象时设置的属性，通常用于初始化对象的状态。通过 construct_properties 列表，可以将一组属性关联到对象的构造过程中。

    3.  **constructor**:
        -   类型: GObject* ()(GType, guint, GObjectConstructParam)
        -   描述: 构造函数指针，用于实例化 GObject 类型的对象。
        -   参数:
            -   type: GObject 类型的 GType。
            -   n_construct_properties: 构造函数需要的属性数量。
            -   construct_properties: 构造函数需要的属性列表。
        -   返回值: 新创建的 GObject 类型的对象的指针。

    4.  **set_property**:
        -   类型: void ()(GObject, guint, const GValue*, GParamSpec*)
        -   描述: 设置对象属性的函数指针。
        -   参数:
            -   object: 要设置属性的 GObject 对象指针。
            -   property_id: 属性的标识符。
            -   value: 属性的值。
            -   pspec: 属性的参数规范。

    5.  **get_property**:
        -   类型: void ()(GObject, guint, GValue*, GParamSpec*)
        -   描述: 获取对象属性的函数指针。
        -   参数:
            -   object: 要获取属性值的 GObject 对象指针。
            -   property_id: 属性的标识符。
            -   value: 存储属性值的 GValue 结构体指针。
            -   pspec: 属性的参数规范。

    6.  **dispose**:
        -   类型: void ()(GObject)
        -   描述: dispose 方法用于释放对象所拥有的资源，但是并不销毁对象本身。通常在此方法中进行资源的清理、断开连接、解除订阅等操作。调用 g_object_unref() 函数时会自动调用 dispose 方法，但也可以手动调用该方法。

    7.  **finalize**:
        -   类型: void ()(GObject)
        -   描述: finalize 方法是在对象的引用计数为0时自动调用的。它用于执行对象的最终化操作，例如释放内存、销毁相关的资源等。finalize 方法是 GObject 类的虚拟方法，可以在子类中对其进行重写以添加自定义的最终化逻辑。通过g_object_ref()可以增加引用计数。

        >   dispose和finalize的区别: 
        >   -   dispose 方法主要用于释放对象所拥有的资源，而不涉及对象自身的销毁。
        >   -   finalize 方法用于对象的最终化操作，在对象被销毁之前执行，包括释放对象自身所占用的内存和资源。

    8.  **dispatch_properties_changed**:
        -   类型: void ()(GObject, guint, GParamSpec**)
        -   描述: 在对象的属性发生改变时调用的函数指针。用于发送属性改变的信号。
        -   参数:
            -   object: 属性发生改变的 GObject 对象指针。
            -   n_pspecs: 发生改变的属性数量。
            -   pspecs: 发生改变的属性的参数规范数组。

    9.  **notify**:
        -   类型: void ()(GObject, GParamSpec*)
        -   描述: 在对象的属性被修改时调用的函数指针。用于发送属性修改的信号。
        -   参数:
            -   object: 属性被修改的 GObject 对象指针。
            -   pspec: 被修改的属性的参数规范。

    10. **constructed**:
        -   类型: void ()(GObject)
        -   描述: 对象构造完成后调用的函数指针。可以在此函数中执行初始化操作或触发其他相关处理。

    11. **flags**:
        -   类型: gsize
        -   描述: 该成员位域标志位，用于存储一些标志信息。可以使用该成员变量来存储与 GObject 类相关的标志或状态信息。它是一个私有成员，一般不做外部调用。
        
    12. **pdummy**:
        -   类型: gpointer[6]
        -   描述: 填充字段，仅用于对齐结构体成员。

2.  GObject结构体的介绍

    GObject内部结构是struct _GObject，它是GObject这个库的基础结构。_GObject的内部结构如下: 
    ```c
    struct  _GObject
    {
    GTypeInstance  g_type_instance;
    
    /*< private >*/
    volatile guint ref_count;
    GData         *qdata;
    };
    ```
    _GObject结构体成员: 

    1.  **g_type_instance**:
            -   类型: GTypeInstance
            -   描述: 用于存储对象的类型信息和实例相关的数据。

    2.  **ref_count**:
        -   类型: volatile guint
        -   描述: 私有成员。用于表示对象的引用计数（reference count）。引用计数是一种内存管理机制，用于跟踪对象被引用的次数。当引用计数变为0时，对象就可以被释放或销毁。
        
    3.  **qdata**:
        -   类型: GData*
        -   描述: 私有成员。用于存储与对象相关的私有数据。GData是GLib库提供的一种数据结构，用于关联任意类型的私有数据到对象上。


### 2. 类和对象
在了解GObject的使用之前，先来认识几个关键的宏: 

1.  **G_DEFINE_TYPE**: GLib 中用于简化 GObject 类型定义的宏。它是一个宏模板，定义了一个用于创建 GObject 类型的标准化流程，包括类型注册、类结构体定义和类初始化函数。

    函数原形: 
    ```c
    G_DEFINE_TYPE (Type, type_name, PARENT_TYPE)
    ```
    参数说明: 

    -   Type: 要定义的 GObject 类型的名称。
    -   type_name: 类型名称的小写形式，用于内部结构体的命名。
    -   PARENT_TYPE: 父类型的 GObject 类型的名称，通常是你自己定义的类的父类。
    -   G_DEFINE_TYPE 宏会自动生成以下内容: 
        -   类型注册函数（Type_register_type()）: 该函数在运行时注册类型以便使用。
        -   内部结构体定义: struct _Type，用于存储实例字段和函数指针。
        -   类初始化函数（Type_class_init()）: 用于初始化类的虚函数表，并可以添加其他初始化逻辑。
        -   实例初始化函数（Type_init()）: 用于初始化实例字段和其他实例特定的逻辑。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】: 
    ```c
    ...
    typedef struct _DbusObj DbusObj;
    typedef struct _DbusObjClass DbusObjClass;

    struct _DbusObj
    {
        /* data */
        BaseObj parent;
    };

    struct _DbusObjClass
    {
        /* data */
        BaseObjClass parent_class;
    };
    // BASE_OBJ_TYPE是父类类型，BASE_OBJ_TYPE是以G_TYPE_OBJECT为基础类
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

2.  ***type_name*_get_type**: 用于获取一个 GObject 类型的 *type_name* 类型标识符（Type Identifier），返回的是GType，GType 是一个代表类型的数据类型。它用于表示在程序中定义的各种数据类型，如对象类、接口、枚举、结构体等。注意，它是通过**G_DEFINE_TYPE**扩展而来的，也就是说它的实现在这个宏中，不需要手动实现。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】: 
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
    运行输出结果: 
    ```sh
    DBUS-OBJ-INFO: 16:59:07.002: dbus_obj.c:50: g_type_name(DBUS_OBJ_TYPE) =DbusObj!
    ```

3.  **G_TYPE_FROM_INSTANCE**: 从给定的 GObject 实例 instance 中获取其对应的 GType。

    函数原形: 
    ```c
    G_TYPE_FROM_INSTANCE(instance)
    ```
    参数说明: 

    -   instance: 要进行类型转换的对象实例。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】: 
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

    运行输出结果: 
    ```sh
    DBUS-OBJ-INFO: 17:09:44.275: dbus_obj.c:51: g_type_name(G_TYPE_FROM_INSTANCE(obj))=DbusObj!
    ```

4.  **G_TYPE_CHECK_INSTANCE_CAST**: GLib 开发中用于类型转换的宏定义之一。将给定的 GObject 实例 instance 转换为目标类型 target_type 的实例。这个宏会检查实例的类型是否与目标类型匹配，如果匹配，则返回目标类型的实例；如果不匹配，则返回 NULL。

    函数原形: 
    ```c
    G_TYPE_CHECK_INSTANCE_CAST(instance, g_type, c_type)
    ```
    参数说明: 

    -   instance: 要进行类型转换的对象实例。
    -   g_type: 目标类型的 GObject 类型。
    -   c_type: 目标类型的 C 类型。

    代码示例【[完整代码示例](./glib-main/dbus_obj.h)】: 
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
    运行输出结果: 
    ```sh
    // 地址都一样，其实就是类型转换
    DBUS-OBJ-INFO: 17:28:51.255: dbus_obj.c:54: obj=0x0x411ac0 base=0x0x411ac0!
    ```
5.  **G_TYPE_INSTANCE_GET_CLASS**: 获取给定 GObject 实例 instance 所属类的类结构体指针。

    函数原形: 
    ```c
    G_TYPE_INSTANCE_GET_CLASS(instance, g_type, c_type)
    ```
    参数说明: 
    -   instance: 需要检查的 GObjectClass 结构体。
    -   g_type: 目标类的 GType。
    -   c_type: 目标类的结构体类型。

    代码示例【[完整代码示例](./glib-main/dbus_obj.c)】: 
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


6.  **G_TYPE_CHECK_INSTANCE_TYPE**: GLib 开发中用于检查类型的宏定义之一。判断给定的 GObject 实例 instance 是否属于指定的类型 g_type。如果实例的类型等于或是指定类型的子类型，则宏会返回 TRUE，否则返回 FALSE。

    函数原形: 
    ```
    G_TYPE_CHECK_INSTANCE_TYPE(instance, g_type) 
    ```
    参数说明: 
    -   instance: 需要检查类型的 GObject 实例。
    -   g_type: 目标类的 GType。

    代码示例【[完整代码示例](./glib-main/main.c)】: 
    ```c
    // base_obj.h
    ...
    #define BASE_OBJ_TYPE (base_obj_get_type()) 
    #define OBJ_IS_BASE(obj)          (G_TYPE_CHECK_INSTANCE_TYPE ((obj), BASE_OBJ_TYPE))
    ...

    // dbus_obj.c
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    #define OBJ_IS_DBUS(obj)          (G_TYPE_CHECK_INSTANCE_TYPE ((obj), DBUS_OBJ_TYPE))

    // main.c
    int main(){
        DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);

        // 检查obj是否是base类型或者是base的子类型
        if(OBJ_IS_BASE(obj)){
            // do somethings
            ...
        }
        ...
    }
    ...

    ```
7.  **G_TYPE_CHECK_CLASS_TYPE**: GLib 开发中用于检查类型的宏定义之一。判断给定的 GObject 类的类型是否与指定的类型 g_type 匹配。如果类的类型等于或是指定类型的子类，则宏会返回 TRUE，否则返回 FALSE。

    函数原形: 
    ```
    G_TYPE_CHECK_CLASS_TYPE(g_class, g_type) 
    ```
    参数说明: 
    -   g_class: 需要检查类型的 GObject 类的类结构体指针。
    -   g_type: 目标类的 GType。

    代码示例【[完整代码示例](./glib-main/main.c)】: 
    ```c
    // base_obj.h
    ...
    #define BASE_OBJ_TYPE (base_obj_get_type()) 
    #define CLASS_IS_BASE_CLASS(klass)  (G_TYPE_CHECK_CLASS_TYPE ((klass),  BASE_OBJ_TYPE))
    ...

    // dbus_obj.c
    // BASE_OBJ_TYPE是父类类型
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)
    #define DBUS_OBJ_TYPE (dbus_obj_get_type())
    #define CLASS_IS_DBUS_CLASS(klass)  (G_TYPE_CHECK_CLASS_TYPE ((klass),  DBUS_OBJ_TYPE))
    #define DBUS_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  DBUS_OBJ_TYPE, DbusObjClass))

    // main.c
    int main(){
        DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);
        DbusObjClass *dbus_class = DBUS_OBJ_GET_CLASS(dbus_class);

        // 检查dbus_class是否是base_class类型或者是base_class的子类型
        if(dbus_class && CLASS_IS_BASE_CLASS(dbus_class)){
            // do somethings
            ...
        }
        ...
    }
    ...

    ```

8.  **G_TYPE_CHECK_CLASS_CAST**: GLib 开发中用于类型转换的宏定义之一。这个宏用于检查一个给定的 GObjectClass 结构体是否与指定的 GType 类型相匹配，并进行相应的类型转换。这个宏通常用于在 GObject 的类型系统中进行类型安全的操作。

    函数原形: 
    ```
    G_TYPE_CHECK_CLASS_CAST(g_class, g_type, c_type)
    ```
    参数说明: 
    -   g_class: 需要检查的 GObjectClass 结构体。
    -   g_type: 目标类的 GType。
    -   c_type: 目标类的结构体类型。

    代码示例【[完整代码示例](./glib-main/main.c)】: 
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

### 3.  继承和多态
GObject支持继承和多态，允许通过派生类扩展已有的类。这种机制可以减少代码重复，并提供更高级的抽象和封装。通过继承，子类可以重写或添加新的方法，从而实现特定的行为。

在前面的例子中其实已经实现了继承，再来看看多态是怎么实现的:

代码示例【[完整代码示例](./glib-main/dbus_obj.c)】:
```c
    // 基类baseclass中包含了一个虚函数base_hello，然后在子类继承的时候重写
    struct _BaseObjClass
    {
        /* data */
        GObjectClass parent_class;
        void	     (*base_hello)		(void);
    };
    ...

    // 构建子类dbus
    typedef struct _DbusObj DbusObj;
    typedef struct _DbusObjClass DbusObjClass;

    struct _DbusObj
    {
        /* data */
        BaseObj parent;
    };

    struct _DbusObjClass
    {
        /* data */
        BaseObjClass parent_class;
    };
    // BASE_OBJ_TYPE是父类类型，BASE_OBJ_TYPE是以G_TYPE_OBJECT为基础类
    // 这个宏包含了DbusObj和DbusObjClass类型注册、类结构体定义和类初始化函数
    G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)

    
    static void dbus_hello (void)
    {
        g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class say hello!");
    }

    static void dbus_obj_class_init (DbusObjClass *klass)
    {
        GObjectClass *object_class = G_OBJECT_CLASS (klass);
        object_class->finalize = finalize;
        object_class->dispose = dispose;
        object_class->constructor = constructor;

        // baseclass可以派生多个子类，都可以对基类的base_hello进行重写
        BaseObjClass *base_class = BASE_OBJ_CLASS(klass);
        base_class->base_hello = dbus_hello;
    }

    ...
    // 在main.c中
    DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);
    // obj可以是base的子类
    BASE_OBJ_GET_CLASS(obj)->base_hello();
    ...

```

### 4.  GParamSpec的介绍和使用
GParamSpec 结构体定义了属性的各种属性，例如名称、类型、默认值、范围等。GParamSpec 用于定义和管理 GObject 的属性系统。它的内部结构是GParamSpec。

结构体原形如下: 
```c
struct _GParamSpec
{
  GTypeInstance  g_type_instance;

  const gchar   *name;          /* interned string */
  GParamFlags    flags;
  GType		 value_type;
  GType		 owner_type;	/* class or interface using this property */

  /*< private >*/
  gchar         *_nick;
  gchar         *_blurb;
  GData		*qdata;
  guint          ref_count;
  guint		 param_id;	/* sort-criteria */
};
```

1.  _GParamSpec结构体成员: 
    1.  g_type_instance:
            -   类型: GTypeInstance
            -   描述: GTypeInstance 结构体用于存储 GObject 类型实例的基本信息。

    2.  name:
        -   类型: const gchar   *
        -   描述: 属性的名称，表示一个内部化的字符串，用作属性的唯一标识符。
        
    3.  flags:
        -   类型: GParamFlags
        -   描述: 属性的标志，用于指定属性的行为和特性，例如属性的可读写性、是否可序列化等。

    4.  value_type:
        -   类型: GType
        -   描述: 属性值的 GType 类型，表示属性所支持的数据类型。

    5.  owner_type:
        -   类型: GType
        -   描述: 拥有该属性的 GObject 类型，可以是类或接口。
        
    6.  _nick:
        -   类型: gchar         *
        -   描述: 私有成员。属性的简短描述，作为属性的昵称。

    7.  _blurb:
        -   类型: gchar         *
        -   描述: 私有成员。属性的详细描述，用于提供有关属性的详细信息。

    8.  qdata:
        -   类型: GData		*
        -   描述: 私有成员。指向属性的用户数据，用于存储属性的附加信息。

    9.  ref_count:
        -   类型: guint
        -   描述: 私有成员。属性的引用计数，用于管理结构体的内存管理。

    10. param_id:
        -   类型: guint
        -   描述: 私有成员。属性的排序标识符，用于属性的排序。

    它的作用以及对比普通属性参数的区别如下: 
    1.  **提供属性元数据**: GParamSpec结构体中包含了属性的各种元数据，如名称、类型、默认值、访问权限等。这些信息可以用于属性的查询、验证和映射等操作。普通的属性只包含属性值本身，而GParamSpec可以提供更多的属性信息。
    2.  **支持属性验证**: GParamSpec结构体中可以定义属性的取值范围、约束条件和验证规则。通过使用这些元数据，可以对属性进行验证，以确保属性值的合法性和一致性。
    3.  **支持信号传递**: GParamSpec结构体还可以与属性相关联的信号进行关联。这样，当属性的值发生变化时，可以触发相关的信号传递，从而允许其他对象对属性的变化做出响应。
    4.  **提供属性的获取和设置接口**: 通过GParamSpec结构体，可以定义属性的获取和设置接口，以方便属性值的读取和修改。这样，属性的访问可以通过get和set方法进行，使得属性的操作更加统一和易用。

2.  GParamSpec的安装
    1.  **g_param_spec_*type*** : 用于创建type类型的属性参数。返回的是GParamSpec的指针。一般用于GObject类构造的时候，常用的类型都有对应的函数添加属性。

        代码示例【[完整代码示例](./glib-main/base_obj.c)】: 
        ```c
        // 添加字符串类型的属性
        obj_properties[PROP_TYPE_STRING] = 
            g_param_spec_string(
                BASE_PROP_STRING,  // 属性名称
                "",  // nick
                "",  // blurb
                "",  // 默认值
                G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
        );
        // 添加INT类型的属性
        obj_properties[PROP_TYPE_INT] = 
            g_param_spec_int(
                BASE_PROP_INT,  // 属性名称
                "int type",  // nick
                "prop type is int",  // blurb
                0,  // 最小值
                100,  // 最大值
                50,  // 默认值
                G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
        );
        // 添加结构体指针类型的属性
        obj_properties[PROP_TYPE_POINTER] = 
            g_param_spec_pointer(
                BASE_PROP_POINTER,  // 属性名称
                "pointer type",  // nick
                "prop type is pointer",  // blurb
                G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
        );

        // 添加结构体类型的属性
        obj_properties[PROP_TYPE_BOXED] = 
            g_param_spec_boxed(
                BASE_PROP_BOXED,  // 属性名称
                "boxed type",  // nick
                "prop type is boxed",  // blurb
                MY_BOXED_POINT_TYPE,     // GType，将结构体转成指针存放
                G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
        );
        ```
    2.  **g_object_class_install_properties**: 
        -   类型: void ()(GObjectClass*, guint, GParamSpec**)
        -   描述: 在GObject类中安装属性列表。这个函数将一组属性参数与GObject类进行关联，使得类的实例可以使用这些属性。
        -   参数:
            -   oclass: 要安装属性参数的GObject类的指针。。
            -   n_pspecs: 属性参数的数量。
            -   pspecs: 一个指向属性参数（GParamSpec）指针的数组。

        代码示例【[完整代码示例](./glib-main/base_obj.c)】:
        ```c
        g_object_class_install_properties(klass,_PROPERTY_ENUMS_LAST,obj_properties);
        ```
    
    3.  **g_object_class_install_property**: 
        -   类型: void ()(GObjectClass*, guint, GParamSpec*)
        -   描述: 在GObject类中安装属性列表。这个函数将一个属性参数与GObject类进行关联，和g_object_class_install_properties功能一样，区别是它安装单个属性。
        -   参数:
            -   class: 要安装属性参数的GObject类的指针。
            -   property_id: 属性的唯一标识符，类型为guint（无符号整数）。
            -   pspecs: 指向GParamSpec结构的指针，描述了要安装的属性的详细信息。

        代码示例【[完整代码示例](./glib-main/base_obj.c)】:
        ```c
        g_object_class_install_property(klass,PROP_TYPE_STRING,obj_properties[PROP_TYPE_STRING]);
        g_object_class_install_property(klass,PROP_TYPE_INT,obj_properties[PROP_TYPE_INT]);
        g_object_class_install_property(klass,PROP_TYPE_POINTER,obj_properties[PROP_TYPE_POINTER]);
        g_object_class_install_property(klass,PROP_TYPE_BOXED,obj_properties[PROP_TYPE_BOXED]);
        ```

3.  GParamSpec值的设置

    GParamSpec值的设置常见的有两个方法可以设置，g_object_set和g_object_set_property: 
    1.  **g_object_set**:
        -   类型: void ()(gpointer, const gchar *, ...)
        -   描述: 用于设置对象的属性值,可设置多个。
        -   参数:
            -   _object: 要设置属性的对象指针。
            -   first_property_name: 第一个属性的名称，类型为const gchar *。该函数可以接受多个属性名称和对应的属性值。
            -   ...: 属性名称和对应的属性值,在设置最后一个参数的值后，还需要加上NULL，不然g_object_set遍历参数的时候会报错。

        代码示例【[完整代码示例](./glib-main/base_obj.c)】:
        ```c
        // 设置int类型属性的值
        // g_object_set(G_OBJECT(obj),BASE_PROP_INT,89,NULL);

        // 设置指针类型属性的值
        MyStruct *prop_pointer = g_new(MyStruct,1);
        prop_pointer->value1 = 123;
        prop_pointer->value2 = "bbbb";
        // g_object_set(G_OBJECT(obj),BASE_PROP_POINTER,prop_pointer,NULL);

        // 设置结构体类型的值
        MyCustomPoint *point = g_new(MyCustomPoint,1);
        point->x = 51;
        point->y = 49;
        // g_object_set(G_OBJECT(obj),BASE_PROP_BOXED,point,NULL);

        // 用一个g_object_set设置三个属性的值
        g_object_set(G_OBJECT(obj),BASE_PROP_INT,89,BASE_PROP_POINTER,prop_pointer,BASE_PROP_BOXED,point,NULL);
        ```

    2.  **g_object_set_property**:
        -   类型: void ()(GObject *, const gchar *, const GValue *)
        -   描述: 用于设置对象的属性值。
        -   参数:
            -   object: 要设置属性的对象指针。
            -   property_name: 属性的名称，类型为const gchar *。
            -   value: 属性名称和对应的属性值,类型为const GValue *。
        ```c
        // 如果要使用g_object_set_property设置属性的值，需要将值转换成 GValue *类型，需要注意GValue的内存管理于释放。
        GValue gval_string = G_VALUE_INIT;
        g_value_init(&gval_string,G_TYPE_STRING);
        g_value_set_string(&gval_string,"asdajksda");
        g_object_set_property(G_OBJECT(obj),BASE_PROP_STRING,&gval_string);
        // GValue的释放
        g_value_unset(&gval_string);
        ```
4.  GParamSpec值的获取

    GParamSpec值的设置常见的有两个方法可以设置，g_object_get和g_object_get_property: 
    1.  g_object_get:
        -   类型: void ()(gpointer, const gchar *, ...)
        -   描述: 用于获取对象的属性值,可获取多个。
        -   参数:
            -   _object: 要获取属性的对象指针。
            -   first_property_name: 第一个属性的名称，类型为const gchar *。该函数可以接受多个属性名称和存放数据的地址。
            -   ...: 属性名称和对应的属性值,在设置最后一个参数的值后，还需要加上NULL，不然g_object_set遍历参数的时候会报错。
        
        ```c
        gchar * val_str;
        g_object_get(G_OBJECT(self),BASE_PROP_STRING,&val_str,NULL);
        g_log(domain, G_LOG_LEVEL_INFO, "base obj %s = %s", BASE_PROP_STRING,val_str);
        g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s = %s\n", BASE_PROP_STRING,self->priv->prop_string);

        gint val_int;
        g_object_get(G_OBJECT(self),BASE_PROP_INT,&val_int,NULL);
        g_log(domain, G_LOG_LEVEL_INFO, "base obj %s = %d", BASE_PROP_INT,val_int);
        g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s = %d\n", BASE_PROP_INT,self->priv->prop_int);
        ```
    2.  **g_object_get_property**
        -   类型: void ()(GObject *, const gchar *, GValue *)
        -   描述: 用于设置对象的属性值。
        -   参数:
            -   object: 要获取属性的对象指针。
            -   property_name: 属性的名称，类型为const gchar *。
            -   value: 获取属性后的存储地址,类型为GValue *。
        ```c
        // 如果要使用g_object_set_property设置属性的值，需要将值转换成 GValue *类型，需要注意GValue的内存管理于释放。
        GValue value = G_VALUE_INIT;
        g_value_init (&value, MY_BOXED_POINT_TYPE);
        g_object_get_property(G_OBJECT(self),BASE_PROP_BOXED,&value);
        MyCustomPoint prop_box = *(MyCustomPoint *)g_value_get_boxed(&value);
        g_log(domain, G_LOG_LEVEL_INFO, "base obj %s: [x=%d y=%d]", BASE_PROP_BOXED,prop_box.x,prop_box.y);
        g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s: [x=%d y=%d]\n", BASE_PROP_BOXED,self->priv->point.x,self->priv->point.y);
        // GValue的释放
        g_value_unset(&value);
        ```        
    
### 5.  signal的介绍和使用
1.  信号的创建

**g_signal_new**: 
-   类型: guint ()(const gchar  *, GType, guint, GSignalAccumulator, gpointer, GSignalCMarshaller, GType, guint, ...)
-   描述: 用于定义一个新的信号。它确定了信号的名称、发送者、标识、参数等属性，并将该信号注册到GObject类型系统中。
-   参数:
    -   signal_name: 信号的名称，用于在代码中标识信号。
    -   itype: 信号发送者的类型，通常使用G_TYPE_FROM_CLASS宏获取。
    -   signal_flags: 信号的标志，指定信号的属性，如运行阶段、累加器行为等。
    -   class_offset: 信号所属类的偏移量（通常为0）。
    -   accumulator: 信号的累加器函数，用于决定多个信号处理函数返回值的合并方式。如果不需要累加器，则可传递NULL。
    -   accu_data: 传递给累加器函数的用户数据。
    -   c_marshaller: 信号的C marshaller函数，用于在信号发射时管理信号参数的内存分配和释放。通常使用NULL。
    -   return_type: 信号的返回值类型。如果信号没有返回值，则使用G_TYPE_NONE。
    -   n_params: 信号的参数数量。
    -   ... : 信号的参数列表。

    代码示例【[完整代码示例](./glib-main/base_obj.c)】:
    ```c
    // 添加信号
    g_signal_new(
        BASE_SIGNAL_STRING_CHANGED,
        G_TYPE_FROM_CLASS(klass),
        G_SIGNAL_RUN_FIRST,
        0,
        NULL,
        NULL,
        NULL,
        G_TYPE_NONE,
        1,
        G_TYPE_STRING
    );
    ```

2.  连接信号处理函数

**g_signal_connect**:
-   类型: gulong ()(gpointer, const gchar *, GCallback, gpointer)
-   描述: 用于定义一个新的信号。它确定了信号的名称、发送者、标识、参数等属性，并将该信号注册到GObject类型系统中。
-   参数:
    -   instance: 指向对象的指针。
    -   detailed_signal: 详细信号的字符串表示。
    -   c_handler: 指向回调函数的指针。
    -   data: 传递给回调函数的用户数据。

    代码示例【[完整代码示例](./glib-main/base_obj.c)】:
    ```c
    // 连接信号处理函数
    static void base_signal_prop_changed(BaseObj *obj, int value, gpointer user_data) {
        // 处理信号的逻辑
        // ...
        BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(obj);
        switch (value)
        {
        case PROP_TYPE_STRING:
            g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to %s\n",BASE_PROP_STRING,priv->prop_string);
            break;
        case PROP_TYPE_INT:
            /* code */
            g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to %d\n",BASE_PROP_INT,priv->prop_int);
            break;
        case PROP_TYPE_POINTER:
            /* code */
            g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to [%d,%s]\n",BASE_PROP_POINTER,priv->prop_pointer->value1,priv->prop_pointer->value2);
            break;
        case PROP_TYPE_BOXED:
            /* code */
            g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to [%d,%d]\n",BASE_PROP_BOXED,priv->point.x,priv->point.y);
            break;
        default:
            break;
        }    
    }
    static void base_obj_init(BaseObj *self)
    {
        g_log(domain, G_LOG_LEVEL_INFO, "base obj init!");
        self->desc = "base";
        g_log(domain, G_LOG_LEVEL_INFO, "base obj init done!\n");
        // 连接信号处理函数
        // 
        int *user_data = g_new0(int,1);
        *user_data = 100;
        // 监听信号名称为"prop-changed"的所有信号，并将user_data传给处理函数
        g_signal_connect(self, BASE_SIGNAL_PROP_CHANGED , G_CALLBACK(base_signal_prop_changed), user_data);

        // 监听信号名称为"prop-changed"的特定的信号，并将user_data传给处理函数
        g_signal_connect(self, BASE_SIGNAL_PROP_CHANGED"::string" , G_CALLBACK(base_signal_prop_changed), user_data);
        g_signal_connect(self, BASE_SIGNAL_PROP_CHANGED"::int" , G_CALLBACK(base_signal_prop_changed), user_data);

    }
    ```

3.  信号的触发
    1.  **g_signal_emit**
    -   类型: void ()(gpointer, guint, GQuark, ...)
    -   描述: 用于定义一个新的信号。它确定了信号的名称、发送者、标识、参数等属性，并将该信号注册到GObject类型系统中。
    -   参数:
        -   instance: 指向对象的指针。
        -   signal_id: 信号的id，即g_signal_new的返回值。
        -   detail: 表示信号的详细信息，可以使用 g_quark_from_static_string 函数将字符串转换为 GQuark 类型。
        -   ... : 可变参数，用于传递信号的参数。
    
    代码示例【[完整代码示例](./glib-main/base_obj.c)】:
    ```c
    static void set_property(GObject *object,guint  property_id,const GValue   *value,GParamSpec    *pspec)
    {
        g_log(domain, G_LOG_LEVEL_INFO, "set_property property_id=%d, type_name=%s!",property_id,g_type_name (pspec->value_type));
        BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(object);
        switch (property_id)
        {
        case PROP_TYPE_STRING:
            priv->prop_string =  g_value_dup_string(value);
            g_signal_emit(object,base_signals[PROP_TYPE_STRING],0,priv->prop_string);
            break;
        ...
        }
    }
    ```


    2.  **g_signal_emit_by_name**
    -   类型: void ()(gpointer, const gchar *, ...)
    -   描述: 用于定义一个新的信号。它确定了信号的名称、发送者、标识、参数等属性，并将该信号注册到GObject类型系统中。
    -   参数:
        -   instance: 指向对象的指针。
        -   detailed_signal: 表示要发送的信号的名称和详细说明，格式为"signal-name::detail"。例如，"clicked"是GtkButton的一个常见信号，"clicked::right-button"表示右键点击事件。
        -   ... : 可变参数，用于传递信号的参数。

    代码示例【[完整代码示例](./glib-main/base_obj.c)】:
    ```c
    static void set_property(GObject *object,guint  property_id,const GValue   *value,GParamSpec    *pspec)
    {
        g_log(domain, G_LOG_LEVEL_INFO, "set_property property_id=%d, type_name=%s!",property_id,g_type_name (pspec->value_type));
        BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(object);
        switch (property_id)
        {
        case PROP_TYPE_STRING:
            priv->prop_string =  g_value_dup_string(value);
            g_signal_emit_by_name(object,BASE_SIGNAL_STRING_CHANGED,priv->prop_string);
            break;
        ...
        }
    }
    ```
