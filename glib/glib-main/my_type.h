#ifndef __MY_TYPE_H__
#define __MY_TYPE_H__

#include <stdio.h>
#include <gio/gio.h>
#include "glib_micros.h"

typedef struct  _MyStruct  MyStruct;
typedef struct  _MyCustomPoint  MyCustomPoint;

typedef struct _MyCustomPoint{
    GBoxedCopyFunc parent;
    gint x;
    gint y;
};

// 自定义结构体，作为GObject的属性
typedef struct _MyStruct
{
  int value1;
  char * value2;
};

GType my_custom_point_get_type(void);
static void my_custom_point_free(MyCustomPoint* point);
static MyCustomPoint* my_custom_point_copy(const MyCustomPoint* src);

#define MY_BOXED_POINT_TYPE my_custom_point_get_type()

#endif