#include <stdio.h>
#include <locale.h>
#include <gio/gio.h>
#include "glib_micros.h"
#include "base_obj.h"
#include "dbus_obj.h"

static const gchar *domain = "glib_main";

void test_g_new0(){
    // 使用g_new0分配内存
    // 分配一个大小为 10 的 int 数组，并初始化为 0，
    int *i = g_new0 (int, 10);
    i[0]= 1;
    i[3]= 11;
    g_free(i);

}

gboolean timeout_callback(gpointer data) {
    int* counter = (int*)data;
    (*counter)++;

    g_log(domain, G_LOG_LEVEL_INFO, "Timeout: %d", *counter);

    if (*counter == 5) {
        // 返回 FALSE 停止重复触发定时器
        return FALSE;
    }
    
    // 返回 TRUE 继续重复触发定时器
    return TRUE;
}


GMainLoop *loop;

int main(){
    // 设置中文环境
    setlocale(LC_ALL, "");

    // 设置env开启debug
    setenv("G_MESSAGES_DEBUG", "all", 1);

    test_g_new0();
    loop = g_main_loop_new (NULL, TRUE);
    // 创建base_obj
    BaseObj *base =  base_obj_new();
    g_log(domain, G_LOG_LEVEL_INFO, "base is dbus? %s!",OBJ_IS_DBUS(base)?"true":"false");

    // 获取base类对象
    BaseObjClass *base_class =  BASE_OBJ_GET_CLASS(base);
    // 验证base_class是否是base类对象
    g_log(domain, G_LOG_LEVEL_INFO, "base_class is base class? %s!",CLASS_IS_BASE_CLASS(base_class)?"true":"false");
    // 验证base_class是否是dbus类对象
    g_log(domain, G_LOG_LEVEL_INFO, "base_class is dbus class? %s!",CLASS_IS_DBUS_CLASS(base_class)?"true":"false");

    base_obj_print_priv(base);
    // 释放资源
    g_object_unref(base);


    // 创建dbus_obj
    DbusObj *dbus = dbus_obj_new();
    // 验证dbus是否是base对象
    g_log(domain, G_LOG_LEVEL_INFO, "dbus is base? %s!",OBJ_IS_BASE(dbus)?"true":"false");

    // 获取dbus类对象
    DbusObjClass *dbus_class = DBUS_OBJ_GET_CLASS(dbus);
    g_log(domain, G_LOG_LEVEL_INFO, "dbus_class is base class? %s!",CLASS_IS_BASE_CLASS(dbus_class)?"true":"false");

    // 将dbus类对象 类型转换成 base类对象
    BaseObjClass *base_dbus_class = BASE_OBJ_GET_CLASS(dbus);

    // 执行base类对象的方法
    base_dbus_class->base_hello();

    base_obj_set_prop(BASE_OBJ(dbus));
    base_obj_print_priv(BASE_OBJ(dbus));

    // 增加dbus的引用计数
    gpointer p1 = g_object_ref(dbus);

    // 释放资源
    g_object_unref(dbus);
    g_object_unref(p1);


    g_main_loop_run (loop);
}