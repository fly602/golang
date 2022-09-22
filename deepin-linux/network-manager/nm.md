#   network-manager 源码解读

##  NetworkManager介绍
-   可以参考wikidev文档：[NetworkManager及相关概念介绍](https://wikidev.uniontech.com/NetworkManager及相关概念介绍)
-   源码目录说明可以参考wikidev文档：[NetworkManager源码目录说明](https://wikidev.uniontech.com/NetworkManager源码目录说明)

补充说明：NetworkManager代码大量用到了Glib-2.0和Gobject-2.0库函数,相关用法说明可以参考一些网上的资料：
-   **Gnome Api参考**：[gnome api参考](https://developer-old.gnome.org/references)， 包含下面所有的api接口，有比较详细的讲解，阅览这一篇就够了！

>   其他Gnome资料：
-   Gobject-2.0：[Gobject-2.0](https://docs.gtk.org/gobject/index.html)
-   Glib-2.0：[Glib-2.0](https://docs.gtk.org/glib/)
-   GModule-2.0：[GModule – 2.0](https://docs.gtk.org/gmodule/index.html)
-   Gio – 2.0：[Gio – 2.0](https://docs.gtk.org/gio/)



##  nm 中gnome底层库api接口的使用
>   nm是一个庞大的工程，包含了许多模块，并且在工程中大量使用了gnome底层库的api接口和宏定义，代码读起来晦涩难懂。可以说是nm的骨架就是这些接口搭建起来的。俗话说，画龙先画骨，我们先来了解一下nm中使用到的这些api接口的功能和用法，才能看懂上层的实现逻辑。

至此，让我们回到nm的代码，看看nm是如何实现的，我们先来看看main函数，在最开始初始化的时候会看到如下代码：
```
int
main (int argc, char *argv[])
{
    ... ...
    g_type_ensure (NM_TYPE_DBUS_MANAGER);
    ... ...
}
```
g_type_ensure：手册上的说明是确认声明的类型已经被注册生效，并且有相关函数_class_init()的生成。详见[g_type_ensure](https://developer-old.gnome.org/gobject/stable/gobject-Type-Information.html#g-type-ensure)。

我们主要来看看NM_TYPE_DBUS_MANAGER这个宏，代码能找到的声明是在src/nm-dbus-manager.h中：
```
#define NM_TYPE_DBUS_MANAGER (nm_dbus_manager_get_type ())
```
然而并没有找到nm_dbus_manager_get_type实现，于是在src/nm-dbus-manager.c中找到：
```
G_DEFINE_TYPE(NMDBusManager, nm_dbus_manager, G_TYPE_OBJECT)
```
nm_dbus_manager_get_type 果然是在G_DEFINE_TYPE中实现了。


