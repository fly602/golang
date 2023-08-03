#include "dbus_obj.h"

static const gchar *domain = "DBUS-OBJ";

// BASE_OBJ_TYPE是父类类型
G_DEFINE_TYPE(DbusObj, dbus_obj, G_TYPE_OBJECT)

static void dbus_obj_init (DbusObj *self)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj init!");
}

static void dbus_obj_class_init (DbusObjClass *klass)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class init!");
    GObjectClass *object_class = G_OBJECT_CLASS (klass);
    object_class->finalize = finalize;
    object_class->dispose = dispose;
    object_class->constructor = constructor;

    // override base_hello
}

static void dbus_hello (void)
{
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class say hello!");
}

static void finalize (GObject *object)
{
    G_OBJECT_CLASS(dbus_obj_parent_class)->finalize(object);
    g_log(domain, G_LOG_LEVEL_INFO, "dbus obj class finalize!");
}

static void dispose(GObject *object)
{
    G_OBJECT_CLASS(dbus_obj_parent_class)->dispose(object);
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

    return obj;
}