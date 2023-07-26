#ifndef BASE_OBJ_H
#define BASE_OBJ_H

#include <stdio.h>
#include <gio/gio.h>
#include "glib_micros.h"
#include "my_type.h"

typedef struct  _BaseObj BaseObj;
typedef struct  _BaseObjClass BaseObjClass;
typedef struct  _BaseObjPriv BaseObjPriv;

struct  _BaseObj
{
    /* data */
    GObject parent;
    gchar * desc;

    /* 私有数据*/
    BaseObjPriv *priv;
};

struct _BaseObjClass
{
    /* data */
    GObjectClass parent_class;
    void	     (*base_hello)		(void);
};
struct  _BaseObjPriv
{
  int prop_int;
  char * prop_string;
  MyStruct *prop_pointer;
  MyCustomPoint point;
};

#define BASE_OBJ_TYPE (base_obj_get_type()) 
#define BASE_OBJ(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), BASE_OBJ_TYPE, BaseObj))
#define BASE_OBJ_CLASS(klass)     (G_TYPE_CHECK_CLASS_CAST ((klass),  BASE_OBJ_TYPE, BaseObjClass))
#define OBJ_IS_BASE(obj)          (G_TYPE_CHECK_INSTANCE_TYPE ((obj), BASE_OBJ_TYPE))
#define CLASS_IS_BASE_CLASS(klass)  (G_TYPE_CHECK_CLASS_TYPE ((klass),  BASE_OBJ_TYPE))
#define BASE_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  BASE_OBJ_TYPE, BaseObjClass))

// 定义获取私有数据的宏
#define BASE_OBJ_GET_PRIVATE(obj) \
  (G_TYPE_INSTANCE_GET_PRIVATE((obj), BASE_OBJ_TYPE, BaseObjPriv))

// 属性
#define BASE_PROP_STRING "prop-string"
#define BASE_PROP_INT "prop-int"
#define BASE_PROP_POINTER "prop-pointer"
#define BASE_PROP_BOXED "prop-boxed"

// 信号
#define BASE_SIGNAL_PROP_CHANGED "prop-changed"

GType base_obj_get_type (void);

BaseObj *base_obj_new();


static void finalize (GObject *object);
static void dispose(GObject *object);

static GObject* constructor(GType type,guint n_construct_properties, GObjectConstructParam *construct_properties);
static void set_property(GObject *object,guint  property_id,const GValue   *value,GParamSpec    *pspec);
static void get_property(GObject *object,guint  property_id,const GValue   *value,GParamSpec    *pspec);
static void base_signal_prop_changed(BaseObj *obj, int value, gpointer user_data);

static void base_hello();
void base_obj_print_priv(BaseObj *self);
void base_obj_set_prop(BaseObj *obj);


#endif