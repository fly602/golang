#ifndef BASE_OBJ_H
#define BASE_OBJ_H

#include <stdio.h>
#include <gio/gio.h>

typedef struct  _BaseObj BaseObj;
typedef struct  _BaseObjClass BaseObjClass;


struct  _BaseObj
{
    /* data */
    GObject parent;
    gchar * desc;
};

struct _BaseObjClass
{
    /* data */
    GObjectClass parent_class;
    void	     (*base_hello)		(void);

};

#define BASE_OBJ_TYPE (base_obj_get_type()) 
#define BASE_OBJ(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), BASE_OBJ_TYPE, BaseObj))
#define BASE_OBJ_CLASS(klass)     (G_TYPE_CHECK_CLASS_CAST ((klass),  BASE_OBJ_TYPE, BaseObjClass))
#define OBJ_IS_BASE(obj)          (G_TYPE_CHECK_INSTANCE_TYPE ((obj), BASE_OBJ_TYPE))
#define CLASS_IS_BASE_CLASS(klass)  (G_TYPE_CHECK_CLASS_TYPE ((klass),  BASE_OBJ_TYPE))
#define BASE_OBJ_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  BASE_OBJ_TYPE, BaseObjClass))

GType base_obj_get_type (void);

BaseObj *base_obj_new();


static void finalize (GObject *object);

static GObject* constructor(GType type,guint n_construct_properties, GObjectConstructParam *construct_properties);

static void base_hello();

#define g_log(log_domain,log_level,format, ...) \
    g_log(log_domain, log_level, "%s:%d: " format, __FILE__, __LINE__, ##__VA_ARGS__)

#endif