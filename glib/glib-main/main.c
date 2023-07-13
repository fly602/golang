#include <stdio.h>
#include <gio/gio.h>
#include "base_obj.h"
#include "dbus_obj.h"

static const gchar *domain = "glib_main";

#define g_log(log_domain,log_level,format, ...) \
    g_log(log_domain, log_level, "%s:%d: " format, __FILE__, __LINE__, ##__VA_ARGS__)

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
    test_g_new0();
    loop = g_main_loop_new (NULL, TRUE);
    // 创建base_obj
    BaseObj *base =  base_obj_new();
    g_log(domain, G_LOG_LEVEL_INFO, "base is dbus? %s!",OBJ_IS_DBUS(base)?"true":"false");

    g_object_unref(base);


    // 创建dbus_obj
    DbusObj *dbus = dbus_obj_new();
    g_log(domain, G_LOG_LEVEL_INFO, "dbus is base? %s!",OBJ_IS_BASE(dbus)?"true":"false");
    DbusObjClass *dbus_class = DBUS_OBJ_GET_CLASS(dbus);
    BaseObjClass *base_dbus_class = BASE_OBJ_CLASS(dbus_class);
    base_dbus_class->base_hello();
    g_object_unref(dbus);
    g_main_loop_run (loop);
}