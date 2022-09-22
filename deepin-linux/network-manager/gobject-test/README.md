#   Gobject example

##  例子代码的编译与执行：
-   编译：make
-   运行：./main

##  简述
GObject，亦称 Glib 对象系统，是一个程序库，它可以帮助我们使用 C 语言编写面向对象程序；它提供了一个通用的动态类型系统（ GType ）、一个基本类型的实现集（如整型、枚举等）、一个基本对象类型 - Gobject 、一个信号系统以及一个可扩展的参数 / 变量体系。

##   主要了解Gobject常用接口功能：
-   GType：通用类型，typedef gsize GType; 底层是unsigned log类型的。用户声明的类和对象注册到gobject后会返回一个GType，后面使用GType就代表了这个，这个GType是由一个全局的链表保存static GSList   *g_once_init_list。
-   g_object_new：创建一个类型type_name的对象对象，并执行type_name_init,type_name_class_init的初始化；
-   g_object_unref：减少对象 的引用计数。当其引用计数降至 0 时，对象将完成（即释放其内存）。
-   G_DEFINE_TYPE：用于类型实现的便捷宏，它声明类初始化函数、对象初始化函数和指向父类的静态变量。
-   g_signal_new：创建新signal。
-   g_signal_connect：将 GC 回调函数连接到特定对象的信号。该处理程序将在信号的默认处理程序之前调用。
-   g_signal_emit：发出信号。

###  来看看G_DEFINE_TYPE的实现：
### [G_DEFINE_TYPE()](https://developer-old.gnome.org/gobject/stable/gobject-Type-Information.html#G-DEFINE-TYPE:CAPS)
```
#define G_DEFINE_TYPE(TN, t_n, T_P)			    G_DEFINE_TYPE_EXTENDED (TN, t_n, T_P, 0, {})
参数说明：
TN：新的类名。
t_n：新的类名，用于和_init,__class_init生成新的函数。
T_P：父类型。
```
用于类型默认初始化的宏，它正在的实现是G_DEFINE_TYPE_EXTENDED这个宏，G_DEFINE_TYPE_EXTENDED实现了一些通用的接口。
nm中的使用：
```
G_DEFINE_TYPE(NMDBusManager, nm_dbus_manager, G_TYPE_OBJECT)
```
参数说明：
-   NMDBusManager：新的类型；
-   nm_dbus_manager：将生成新的接口，例如nm_dbus_manager_init、nm_dbus_manager_class_init；
-   G_TYPE_OBJECT：父类型是GObject。

### [G_DEFINE_TYPE_EXTENDED](https://developer-old.gnome.org/gobject/stable/gobject-Type-Information.html#G-DEFINE-TYPE-EXTENDED:CAPS)
```
gobject中的实现：
#define G_DEFINE_TYPE_EXTENDED(TN, t_n, T_P, _f_, _C_)	    _G_DEFINE_TYPE_EXTENDED_BEGIN (TN, t_n, T_P, _f_) {_C_;} _G_DEFINE_TYPE_EXTENDED_END()

/* This was defined before we had G_DEFINE_TYPE_WITH_CODE_AND_PRELUDE, it's simplest
 * to keep it.
 */
#define _G_DEFINE_TYPE_EXTENDED_BEGIN(TypeName, type_name, TYPE_PARENT, flags) \
  _G_DEFINE_TYPE_EXTENDED_BEGIN_PRE(TypeName, type_name, TYPE_PARENT) \
  _G_DEFINE_TYPE_EXTENDED_BEGIN_REGISTER(TypeName, type_name, TYPE_PARENT, flags) \

/* 最后得到如下声明类型type_name实现 */
static void     type_name##_init              (TypeName        *self); \ 需要type_name自己实现
static void     type_name##_class_init        (TypeName##Class *klass); \ 需要type_name自己实现
static gpointer type_name##_parent_class = NULL; \
static void     type_name##_class_intern_init (gpointer klass) \    宏定义已经实现
GType type_name##_get_type (void)
```
在 GObject 世界里，类是两个结构体的组合，一个是对象结构体，另一个是类结构体。即**TypeName**是对象的结构体，**TypeName##Class**是类结构体。它们组合起来便是一个**TypeName**类。

### 好了，为什么要用GObject？
-   基于引用计算器管理内存g_object_ref/g_object_unref
-   对象的构造函数与析构函数
-   可设置对象属性的 g_object_get/set_property 函数
-   易于使用的信号机制

上述代码，我们当前来看只需要实现_init() 和_class_init()方法即可，这是GObject的命名约定。对G_DEFINE_TYPE宏进行展开后的说明如下：
- type_name##_init：type_name的对象的初始化，G_DEFINE_TYPE只做了声明，没有具体的实现，需要程序员自己实现。
- type_name##_class_init：type_name类的初始化
- type_name##_parent_class：type_name 类的父类，在type_name##_class_intern_init中初始化，指向**TypeName##Class**结构体中的parent
- type_name##_class_intern_init：type_name 类的父类初始化函数。



