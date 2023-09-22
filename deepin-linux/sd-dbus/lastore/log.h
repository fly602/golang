#ifndef __LOG_H__
#define __LOG_H__

#include <syslog.h>
#include <stdio.h>

// #define LOG(level,format,...) 
//     syslog(level,"%s:%d "format,__FILE__,__LINE__,##__VA_ARGS__)

#define LOG(level,format,...) \
    fprintf(stderr,"%s:%d "format"\n",__FILE__,__LINE__,##__VA_ARGS__)
#endif