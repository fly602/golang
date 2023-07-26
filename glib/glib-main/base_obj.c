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

GOBJECT_PROPERTIES_DEFINE_BASE(
    PROP_TYPE_STRING,      // 字符串类型的属性
    PROP_TYPE_INT,          // int类型的属性
    PROP_TYPE_POINTER,          // 结构体指针类型的属性
    PROP_TYPE_BOXED,          // 结构体类型的属性
);

static guint base_signals[_PROPERTY_ENUMS_LAST] = {0};

static void finalize(GObject *object)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class finalize!");
}

static void dispose(GObject *object)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class dispose!");
    BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(object);
    g_free(priv->prop_pointer);
}

static void set_property(GObject *object,guint  property_id,const GValue   *value,GParamSpec    *pspec)
{
    g_log(domain, G_LOG_LEVEL_INFO, "set_property property_id=%d, type_name=%s!",property_id,g_type_name (pspec->value_type));
    BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(object);
    switch (property_id)
    {
    case PROP_TYPE_STRING:
        priv->prop_string =  g_value_dup_string(value);
        g_signal_emit_by_name(object,BASE_SIGNAL_PROP_CHANGED "::string",PROP_TYPE_STRING);
        // g_signal_emit(object,base_signals[PROP_TYPE_STRING],0,priv->prop_string);
        break;
    case PROP_TYPE_INT:
        priv->prop_int =  g_value_get_int(value);
        g_signal_emit_by_name(object,BASE_SIGNAL_PROP_CHANGED "::int",PROP_TYPE_INT);
        break;
    case PROP_TYPE_POINTER:
        priv->prop_pointer =  (MyStruct *)g_value_get_pointer(value);
        g_signal_emit_by_name(object,BASE_SIGNAL_PROP_CHANGED "::pointer",PROP_TYPE_POINTER);
        break;
    case PROP_TYPE_BOXED:
        priv->point =  *(MyCustomPoint *)g_value_get_boxed(value);
        g_signal_emit_by_name(object,BASE_SIGNAL_PROP_CHANGED "::boxed",PROP_TYPE_BOXED);
        break;
    default:
        break;
    }
    g_log(domain, G_LOG_LEVEL_INFO, "set_property %d done!\n",property_id);
}

static void get_property(GObject *object,guint  property_id,const GValue   *value,GParamSpec    *pspec)
{
    g_log(domain, G_LOG_LEVEL_INFO, "get_property property_id=%d, type_name=%s!",property_id,g_type_name (pspec->value_type));
    BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(object);
    switch (property_id)
    {
    case PROP_TYPE_STRING:
        g_value_set_string(value,priv->prop_string);
        break;
    case PROP_TYPE_INT:
        g_value_set_int(value,priv->prop_int);
        break;
    case PROP_TYPE_POINTER:
        g_value_set_pointer(value,priv->prop_pointer);
        break;
    case PROP_TYPE_BOXED:
        // 注意：g_value_set_boxed会增加gobject的引用计数，使用g_value_take_boxed
        g_value_take_boxed(value,&priv->point);
        // g_value_set_boxed(value,&priv->point);

        break;
    default:
        break;
    }
    g_log(domain, G_LOG_LEVEL_INFO, "get_property property_id=%d done!\n",property_id);
}

static GObject *constructor(GType type, guint n_construct_properties, GObjectConstructParam *construct_properties)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class constructor!");
    GObject *parent = G_OBJECT_CLASS(base_obj_parent_class);
    GObject *obj = G_OBJECT_CLASS(base_obj_parent_class)->constructor(type, n_construct_properties, construct_properties);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class constructor done!\n");
    return obj;
}

// base类的方法
static void base_hello()
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class say hello!");
}

BaseObj *base_obj_new()
{
    g_log(domain, G_LOG_LEVEL_INFO, "new base obj!");
    BaseObj *obj = g_object_new(BASE_OBJ_TYPE, NULL);
    g_log(domain, G_LOG_LEVEL_INFO, "new base obj done!\n");

    base_obj_set_prop(obj);
    return obj;
}

static void base_obj_init(BaseObj *self)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj init!");
    self->desc = "base";
    g_log(domain, G_LOG_LEVEL_INFO, "base obj init done!\n");
    // 连接信号处理函数
    int *user_date = g_new0(int,1);
    *user_date = 100;

    // 监听信号名称为"prop-changed"的所有信号
    // g_signal_connect(self, BASE_SIGNAL_PROP_CHANGED , G_CALLBACK(base_signal_prop_changed), user_date);

    // 监听信号名称为"prop-changed"的特定的信号
    g_signal_connect(self, BASE_SIGNAL_PROP_CHANGED"::string" , G_CALLBACK(base_signal_prop_changed), user_date);
    g_signal_connect(self, BASE_SIGNAL_PROP_CHANGED"::int" , G_CALLBACK(base_signal_prop_changed), user_date);

}

static void base_obj_class_init(BaseObjClass *klass)
{
    g_log(domain, G_LOG_LEVEL_INFO, "base obj class init!");
    GObjectClass *object_class = G_OBJECT_CLASS(klass);

    // 添加私有成员
    g_type_class_add_private(klass,sizeof(BaseObjPriv));

    object_class->finalize = finalize;
    object_class->dispose = dispose;
    object_class->constructor = constructor;
    object_class->set_property = set_property;
    object_class->get_property = get_property;

    klass->base_hello = base_hello;

    // 添加信号
    base_signals[PROP_TYPE_STRING] = g_signal_new(
        BASE_SIGNAL_PROP_CHANGED,
        G_TYPE_FROM_CLASS(klass),
        G_SIGNAL_RUN_FIRST | G_SIGNAL_DETAILED,
        0,
        NULL,
        NULL,
        NULL,
        G_TYPE_NONE,
        1,
        G_TYPE_POINTER
    );

    // 添加属性
    // 添加字符串类型的属性
    obj_properties[PROP_TYPE_STRING] = 
        g_param_spec_string(
            BASE_PROP_STRING,  // 属性名称
            "string type",  // nick
            "prop type is string",  // blurb
            "",  // 默认值
            G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
    );

    // 添加INT类型的属性
    obj_properties[PROP_TYPE_INT] = 
        g_param_spec_int(
            BASE_PROP_INT,  // 属性名称
            "int type",  // nick
            "prop type is int",  // blurb
            0,  // 最小值
            100,  // 最大值
            50,  // 默认值
            G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
    );
    // 添加结构体指针类型的属性
    obj_properties[PROP_TYPE_POINTER] = 
        g_param_spec_pointer(
            BASE_PROP_POINTER,  // 属性名称
            "pointer type",  // nick
            "prop type is pointer",  // blurb
            G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
    );

        // 添加结构体类型的属性
    obj_properties[PROP_TYPE_BOXED] = 
        g_param_spec_boxed(
            BASE_PROP_BOXED,  // 属性名称
            "boxed type",  // nick
            "prop type is boxed",  // blurb
            MY_BOXED_POINT_TYPE,     // GType，将结构体转成指针存放
            G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS // 属性的读写标志
    );

    // 一次安装多个属性
    g_object_class_install_properties(klass,_PROPERTY_ENUMS_LAST,obj_properties);

    // 单个安装属性
    // g_object_class_install_property(klass,PROP_TYPE_STRING,obj_properties[PROP_TYPE_STRING]);
    // g_object_class_install_property(klass,PROP_TYPE_INT,obj_properties[PROP_TYPE_INT]);
    // g_object_class_install_property(klass,PROP_TYPE_POINTER,obj_properties[PROP_TYPE_POINTER]);
    // g_object_class_install_property(klass,PROP_TYPE_BOXED,obj_properties[PROP_TYPE_BOXED]);

    // 遍历已安装的属性
    guint n_params, i;
    GParamSpec **params = g_object_class_list_properties(klass, &n_params);
    for (i = 0; i < n_params; i++) {
        GParamSpec *param_spec = params[i];
        const gchar *name = g_param_spec_get_name(param_spec);
        const gchar *nick = g_param_spec_get_nick(param_spec);
        const gchar *blurb = g_param_spec_get_blurb(param_spec);
        
        g_log(domain, G_LOG_LEVEL_INFO, "Property name: %s", name);
        g_log(domain, G_LOG_LEVEL_INFO, "Property nickname: %s", nick);
        g_log(domain, G_LOG_LEVEL_INFO, "Property description: %s\n", blurb);
    }
    g_free(params);

    g_log(domain, G_LOG_LEVEL_INFO, "base obj class init done!\n");
}

void base_obj_print_priv(BaseObj *self)
{
    gchar * val_str;
    g_object_get(G_OBJECT(self),BASE_PROP_STRING,&val_str,NULL);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj %s = %s", BASE_PROP_STRING,val_str);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s = %s\n", BASE_PROP_STRING,self->priv->prop_string);

    gint val_int;
    g_object_get(G_OBJECT(self),BASE_PROP_INT,&val_int,NULL);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj %s = %d", BASE_PROP_INT,val_int);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s = %d\n", BASE_PROP_INT,self->priv->prop_int);

    // 打印指针类型属性的内容
    gpointer val_pointer;
    g_object_get(G_OBJECT(self),BASE_PROP_POINTER,&val_pointer,NULL);
    MyStruct *prop_pointer = val_pointer;
    g_log(domain, G_LOG_LEVEL_INFO, "base obj %s: [value1=%d value2=%s]", BASE_PROP_POINTER,prop_pointer->value1,prop_pointer->value2);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s: [value1=%d value2=%s]\n",BASE_PROP_POINTER,self->priv->prop_pointer->value1,self->priv->prop_pointer->value2); 

    // // 获取并打印结构体属性的内容
    // gpointer val_boxed;
    // g_object_get(G_OBJECT(self),BASE_PROP_BOXED,&val_boxed,NULL);
    // MyCustomPoint prop_box = *(MyCustomPoint *)val_boxed;
    // g_log(domain, G_LOG_LEVEL_INFO, "base obj %s: [x=%d y=%d]", BASE_PROP_BOXED,prop_box.x,prop_box.y);
    // g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s: [x=%d y=%d]\n", BASE_PROP_BOXED,self->priv->point.x,self->priv->point.y);

    // 使用g_object_get_property的方式获取属性
    GValue value = G_VALUE_INIT;
    g_value_init (&value, MY_BOXED_POINT_TYPE);
    g_object_get_property(G_OBJECT(self),BASE_PROP_BOXED,&value);
    MyCustomPoint prop_box = *(MyCustomPoint *)g_value_get_boxed(&value);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj %s: [x=%d y=%d]", BASE_PROP_BOXED,prop_box.x,prop_box.y);
    g_log(domain, G_LOG_LEVEL_INFO, "base obj self %s: [x=%d y=%d]\n", BASE_PROP_BOXED,self->priv->point.x,self->priv->point.y);
    g_value_unset(&value);

}

void base_obj_set_prop(BaseObj *obj)
{
    obj->priv =  BASE_OBJ_GET_PRIVATE(obj);

    // 设置字符串类型属性的值
    // g_object_set(G_OBJECT(obj),BASE_PROP_STRING,"asdajksda",NULL);
    GValue gval_string = G_VALUE_INIT;
    g_value_init(&gval_string,G_TYPE_STRING);
    g_value_set_string(&gval_string,"asdajksda");
    g_object_set_property(G_OBJECT(obj),BASE_PROP_STRING,&gval_string);
    g_value_unset(&gval_string);
    // 设置int类型属性的值
    // g_object_set(G_OBJECT(obj),BASE_PROP_INT,89,NULL);

    // 设置指针类型属性的值
    MyStruct *prop_pointer = g_new(MyStruct,1);
    prop_pointer->value1 = 123;
    prop_pointer->value2 = "bbbb";
    // g_object_set(G_OBJECT(obj),BASE_PROP_POINTER,prop_pointer,NULL);

    // 设置结构体类型的值
    MyCustomPoint *point = g_new(MyCustomPoint,1);
    point->x = 51;
    point->y = 49;
    // g_object_set(G_OBJECT(obj),BASE_PROP_BOXED,point,NULL);

    // 用一个g_object_set设置三个属性的值
    g_object_set(G_OBJECT(obj),BASE_PROP_INT,89,BASE_PROP_POINTER,prop_pointer,BASE_PROP_BOXED,point,NULL);
}

static void base_signal_prop_changed(BaseObj *obj, int value, gpointer user_data) {
    // 处理信号的逻辑
    // ...
    BaseObjPriv *priv = BASE_OBJ_GET_PRIVATE(obj);
    switch (value)
    {
    case PROP_TYPE_STRING:
        g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to %s\n",BASE_PROP_STRING,priv->prop_string);
        break;
    case PROP_TYPE_INT:
        /* code */
        g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to %d\n",BASE_PROP_INT,priv->prop_int);
        break;
    case PROP_TYPE_POINTER:
        /* code */
        g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to [%d,%s]\n",BASE_PROP_POINTER,priv->prop_pointer->value1,priv->prop_pointer->value2);
        break;
    case PROP_TYPE_BOXED:
        /* code */
        g_log(domain, G_LOG_LEVEL_INFO, "base handle signal %s change to [%d,%d]\n",BASE_PROP_BOXED,priv->point.x,priv->point.y);
        break;
    default:
        break;
    }    
}