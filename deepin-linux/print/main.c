#include <stdio.h>
#include <stdarg.h>

void myprintf(const char* fmt, ...){
    va_list args;
    va_start(args, fmt);
    printf(fmt,args);
    va_end(args);
}

char *arr[] = {"1","2","3"};

int main(){
    myprintf("=====%s %s\n","aa","bb");
    printf("%s %d""===>>arr len=\n",__FILE__,__LINE__);
}