#include "base_obj.h"

static const gchar *domain = "BASE-OBJ";

G_DEFINE_TYPE(BaseObj, base_obj, G_TYPE_OBJECT)

// static void base_obj_init(BaseObj *self);
// static void base_obj_class_init(BaseObjClass *klass);
// static GType base_obj_get_type_once(void);
// static gpointer base_obj_parent_class = ((void *)0);
// static gint BaseObj_private_offset;
// static void base_obj_class_intern_init(gpointer klass)
// {
//     base_obj_parent_class = g_type_class_peek_parent(klass);
//     if (BaseObj_private_offset != 0)
//         g_type_class_adjust_private_offset(klass, &BaseObj_private_offset);
//     base_obj_class_init((BaseObjClass *)klass);
// }
// __attribute__((__unused__)) static inline gpointer base_obj_get_instance_private(BaseObj *self) { return (((gpointer)((guint8 *)(self) + (glong)(BaseObj_private_offset)))); }
// GType base_obj_get_type(void)
// {
//     static volatile gsize g_define_type_id__volatile = 0;
//     if ((__extension__({
//             typedef char _GStaticAssertCompileTimeAssertion_7[(sizeof *(&g_define_type_id__volatile) == sizeof(gpointer)) ? 1 : -1] __attribute__((__unused__));
//             (void)(0 ? (gpointer) * (&g_define_type_id__volatile) : 0);
//             (!(__extension__({
//                 typedef char _GStaticAssertCompileTimeAssertion_8[(sizeof *(&g_define_type_id__volatile) == sizeof(gpointer)) ? 1 : -1] __attribute__((__unused__));
//                 __sync_synchronize();
//                 (gpointer) * (&g_define_type_id__volatile);
//             })) &&
//              g_once_init_enter(&g_define_type_id__volatile));
//         })))
//     {
//         GType g_define_type_id = base_obj_get_type_once();
//         (__extension__({
//             typedef char _GStaticAssertCompileTimeAssertion_9[(sizeof *(&g_define_type_id__volatile) == sizeof(gpointer)) ? 1 : -1] __attribute__((__unused__));
//             (void)(0 ? *(&g_define_type_id__volatile) = (g_define_type_id) : 0);
//             g_once_init_leave((&g_define_type_id__volatile), (gsize)(g_define_type_id));
//         }));
//     }
//     return g_define_type_id__volatile;
// }
// __attribute__((noinline)) static GType base_obj_get_type_once(void)
// {
//     GType g_define_type_id = g_type_register_static_simple(((GType)((20) << (2))), g_intern_static_string("BaseObj"), sizeof(BaseObjClass), (GClassInitFunc)(void (*)(void))base_obj_class_intern_init, sizeof(BaseObj), (GInstanceInitFunc)(void (*)(void))base_obj_init, (GTypeFlags)0);
//     {
//         {
//             {};
//         }
//     }
//     return g_define_type_id;
// }

static void base_obj_init(BaseObj *self)
{
    self->desc = "base";
    g_log(domain, G_LOG_LEVEL_INFO, "base obj init!");
}

static void base_obj_class_init(BaseObjClass *klass)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class init!");
    GObjectClass *object_class = G_OBJECT_CLASS(klass);
    object_class->finalize = finalize;
    object_class->constructor = constructor;

    klass->base_hello = base_hello;
}

static void finalize(GObject *object)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class finalize!");
}

static GObject *constructor(GType type, guint n_construct_properties, GObjectConstructParam *construct_properties)
{
    GObject *parent = G_OBJECT_CLASS(base_obj_parent_class);
    GObject *obj = G_OBJECT_CLASS(base_obj_parent_class)->constructor(type, n_construct_properties, construct_properties);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class constructor!");
    return obj;
}

static void base_hello()
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class say hello!");
}

BaseObj *base_obj_new()
{
    BaseObj *obj = g_object_new(BASE_OBJ_TYPE, NULL);
    g_log(domain, G_LOG_LEVEL_INFO, "new base obj!");

    return obj;
}