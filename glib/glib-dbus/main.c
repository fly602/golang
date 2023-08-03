#include <stdio.h>
#include <locale.h>
#include "dbus_obj.h"
#include "glib_micros.h"

GMainLoop *loop;

int main(){
    // 设置中文环境
    setlocale(LC_ALL, "");

    // 设置env开启debug
    setenv("G_MESSAGES_DEBUG", "all", 1);

    g_type_init();

    loop = g_main_loop_new (NULL, TRUE);
    DbusObj * gdbus = dbus_obj_new();

    g_object_unref(gdbus);
    
    g_main_loop_run (loop);
}