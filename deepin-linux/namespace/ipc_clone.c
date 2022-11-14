#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <unistd.h>
#include <sched.h>
#define STACK_SIZE (1024*1024)

static char container_stack[STACK_SIZE];
const char * args  = "/bin/bash";

int contain_func(void * arg)
{
        printf("this is in %s, and pid : %d \n", __func__, getpid());

        sethostname("alexander", 10);
        system("mount -t proc proc /proc");
        execv(args, arg);
        // printf("this is in %s end\n", __func__);
        return 1;
}


int main(void)
{
        int clone_pid = clone(contain_func, container_stack + STACK_SIZE,
                              CLONE_NEWPID | CLONE_NEWUTS| CLONE_NEWIPC|CLONE_NEWNS |SIGCHLD , NULL);

        waitpid(clone_pid, NULL, 0);
        printf("this is in %s\n", __func__);
        return 0;
}