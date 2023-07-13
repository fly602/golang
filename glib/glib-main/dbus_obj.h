#ifndef DBUS_OBJ_H
#define DBUS_OBJ_H

#include <stdio.h>
#include <gio/gio.h>
#include "base_obj.h"

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

#define DBUS_OBJ_TYPE (dbus_obj_get_type()) 
#define DBUS_OBJ(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), DBUS_OBJ_TYPE, DbusObj))
#define DBUS_OBJ_CLASS(klass)     (G_TYPE_CHECK_CLASS_CAST ((klass),  DBUS_OBJ_TYPE, DbusObjClass))
#define OBJ_IS_DBUS(obj)          (G_TYPE_CHECK_INSTANCE_TYPE ((obj), DBUS_OBJ_TYPE))
#define CLASS_IS_DBUS_CLASS(klass)  (G_TYPE_CHECK_CLASS_TYPE ((klass),  DBUS_OBJ_TYPE))
#define DBUS_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  DBUS_OBJ_TYPE, DbusObjClass))


GType dbus_obj_get_type (void);

DbusObj *dbus_obj_new();

static void finalize (GObject *object);

static GObject* constructor(GType type,guint n_construct_properties, GObjectConstructParam *construct_properties);

static void dbus_hello (void);

#define g_log(log_domain,log_level,format, ...) \
    g_log(log_domain, log_level, "%s:%d: " format, __FILE__, __LINE__, ##__VA_ARGS__)

#endif