#include "my_type.h"

G_DEFINE_BOXED_TYPE(MyCustomPoint, my_custom_point, my_custom_point_copy, my_custom_point_free)

// 自定义类型的拷贝函数
static MyCustomPoint *my_custom_point_copy(const MyCustomPoint *src)
{
    MyCustomPoint *dest = g_new(MyCustomPoint, 1);
    dest->x = src->x;
    dest->y = src->y;
    return dest;
}

// 自定义类型的释放函数
static void my_custom_point_free(MyCustomPoint *point)
{
    // g_free(point);
}