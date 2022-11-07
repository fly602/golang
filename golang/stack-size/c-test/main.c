#include <stdio.h>

// 测试栈溢出，一般情况内存中栈段的大小是固定的8M，超出就会报崩溃
#define STACK_SIZE (8 * 1024 * 1024)

int main(){
    int a[2 * STACK_SIZE];
    a[0] = 1;
    a[2 * STACK_SIZE -1] = 2;
}