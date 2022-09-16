#   network-manager 源码解读

##  NetworkManager介绍
-   可以参考wikidev文档：[NetworkManager及相关概念介绍](https://wikidev.uniontech.com/NetworkManager及相关概念介绍)
-   源码目录说明可以参考wikidev文档：[NetworkManager源码目录说明](https://wikidev.uniontech.com/NetworkManager源码目录说明)

补充说明：NetworkManager代码大量用到了Glib-2.0和Gobject-2.0库函数,相关用法说明可以参考一些网上的资料：
-   Gobject-2.0：[Gobject-2.0](https://docs.gtk.org/gobject/index.html)
-   Glib-2.0：[Glib-2.0](https://docs.gtk.org/glib/)
-   GModule-2.0：[GModule – 2.0](https://docs.gtk.org/gmodule/index.html)
-   Gio – 2.0：[Gio – 2.0](https://docs.gtk.org/gio/)

nm 起始main

