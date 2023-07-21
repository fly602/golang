#include "dbus_obj.h"

static const gchar *domain = "DBUS-OBJ";

enum {
	NOTIFY_SIGNAL,

	LAST_SIGNAL
};

#define DBUS_NOTIFY "dbus-notify"

static guint signals[LAST_SIGNAL] = { 0 };

// BASE_OBJ_TYPE是父类类型
G_DEFINE_TYPE(DbusObj, dbus_obj, BASE_OBJ_TYPE)

static void dbus_obj_init (DbusObj *self)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj init!");
    BaseObj *base = BASE_OBJ(self);
}

static void dbus_obj_class_init (DbusObjClass *klass)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class init!");
    GObjectClass *object_class = G_OBJECT_CLASS (klass);
    object_class->finalize = finalize;
    object_class->dispose = dispose;
    object_class->constructor = constructor;

    // override base_hello
    BaseObjClass *base_class = BASE_OBJ_CLASS(klass);
    base_class->base_hello = dbus_hello;
}

static void dbus_hello (void)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class say hello!");
}

static void finalize (GObject *object)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class finalize!");
}

static void dispose(GObject *object)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class dispose!");
}

static GObject* constructor(GType type,guint n_construct_properties, GObjectConstructParam *construct_properties)
{
    GObject *parent = G_OBJECT_CLASS(dbus_obj_parent_class);
    GObject* obj = G_OBJECT_CLASS(dbus_obj_parent_class)->constructor(type, n_construct_properties, construct_properties);
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class constructor!");
    GObjectClass *klass;
    return obj;
}

DbusObj *dbus_obj_new()
{
    DbusObj *obj = g_object_new(DBUS_OBJ_TYPE, NULL);
    g_log(domain, G_LOG_LEVEL_INFO, "new dbus obj!");
    g_log(domain, G_LOG_LEVEL_INFO, "g_type_name(DBUS_OBJ_TYPE) =%s!",g_type_name(DBUS_OBJ_TYPE));
    g_log(domain, G_LOG_LEVEL_INFO, "g_type_name(G_TYPE_FROM_INSTANCE(obj))=%s!",g_type_name(G_TYPE_FROM_INSTANCE(obj)));

    BaseObj *base = BASE_OBJ(obj);
    g_log(domain, G_LOG_LEVEL_INFO, "obj=0x%p base=0x%p!",obj,base);
    return obj;
}