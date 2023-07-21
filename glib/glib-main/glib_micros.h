#ifndef __GLIB_MACROS_H__
#define __GLIB_MACROS_H__

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>

#include <gio/gio.h>

#define GOBJECT_PROPERTIES_DEFINE_BASE(...) \
typedef enum { \
	PROP_0, \
	__VA_ARGS__ \
	_PROPERTY_ENUMS_LAST, \
} _PropertyEnums; \
static GParamSpec *obj_properties[_PROPERTY_ENUMS_LAST] = { NULL, }

// 封装日志函数，格式化
#define g_log(log_domain,log_level,format, ...) \
    g_log(log_domain, log_level, "%s:%d: " format, __FILE__, __LINE__, ##__VA_ARGS__)

#endif