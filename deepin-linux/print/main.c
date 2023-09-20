#include <stdio.h>
#include <stdarg.h>

void myprintf(const char* fmt, ...){
    va_list args;
    va_start(args, fmt);
    printf(fmt,args);
    va_end(args);
}