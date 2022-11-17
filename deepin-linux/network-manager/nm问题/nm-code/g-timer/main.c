// test-gmainloop.c
#include <stdio.h>
#include <glib.h>

pthread_t thid1;
pthread_t thid2;

int handle_state = FALSE;

int *func1();
int *func2();

static gboolean dns_work (gpointer arg)
{
    sleep(1);
    printf("dns_work done!\n");
    handle_state = TRUE;
    
    return FALSE;
}


int *func1(){
    //线创建一个工作线程，模拟正在处理域名解析工作
    // 对应的是nm_connectivity_check_start 中的curl_multi_add_handle
    if(pthread_create(&thid1, NULL, (void *)dns_work, NULL) != 0) {
		printf("thread creation failed\n");
		exit(1);
	}
    // 接下来添加一个定时器，处理dns解析的结果
    g_timeout_add_seconds (5, func2, NULL); // 1s
    return NULL;
}

// 模拟nm的_timeout_cb函数，
// curl_multi_remove_handle函数是阻塞的，最终调用的就是pthread_join
int *func2(){
    printf("working for dns_work\n");
    if (handle_state== FALSE){
        pthread_join(thid1,NULL);
    }
    printf("work for dns_work done\n");
    return NULL;
}

int main()
{
    static GMainLoop *main_loop = NULL;
    main_loop = g_main_loop_new (NULL, FALSE);

    // 再创建一个线程模拟添加定时器，处理域名解析的结果
    if(pthread_create(&thid2, NULL, (void *)func1, NULL) != 0) {
		printf("thread creation failed\n");
		exit(1);
	}
    g_main_loop_run (main_loop);
    return 0;
}