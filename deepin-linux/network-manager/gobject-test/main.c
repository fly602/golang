#include "Child1.h"  

static void test_hello_callback_connect(TESTChild1 *pBase,int i,gpointer user_data)
{
	int j = 0;
	printf("TESTChild1[%p] recv hello signal connect[%d].\n",pBase,i);
	for(j = 0;j < 2047483648;j++);
	sleep(3);
	printf("TESTChild1[%p] recv hello signal connect[%d] end.\n",pBase,i);
}

static void test_hello_callback_connect1(TESTChild1 *pBase,int i,gpointer user_data)
{
	int j = 0;
	printf("TESTChild1[%p] recv hello signal connect1[%d].\n",pBase,i);
	for(j = 0;j < 2047483648;j++);
	sleep(3);
	printf("TESTChild1[%p] recv hello signal connect1[%d] end.\n",pBase,i);
}

int main (void)  
{  
    g_type_init ();  
    int i;  
    TESTChild1 *P,*P1,*P2;
	// 只有实现了_init, _class_init才能创建object对象。
	//g_object_new(&i,NULL);
#if 1
	printf(">new first child1 object.\n");
	printf("Base Object Type=%ld\n",(unsigned long)TEST_TYPE_BASE);
	printf("Base Object Type=%ld\n",(unsigned long)TEST_TYPE_CHILD1);
	P = g_object_new (TEST_TYPE_CHILD1, NULL);  
	printf("P->parent=%x \n",P->parent);
	TESTChild1 *P3 = g_object_ref(P);
    g_object_unref (P); 
	printf("new second child1 object.\n");
	printf("Base Object Type=%ld\n",(unsigned long)TEST_TYPE_BASE);
	P1 = g_object_new (TEST_TYPE_CHILD1, NULL);  
    g_object_unref (P1); 
	printf("new Base object.\n");
	printf("Base Object Type=%ld\n",(unsigned long)TEST_TYPE_BASE);
	P2 = g_object_new(TEST_TYPE_BASE,NULL);
	g_object_unref (P2); 
	printf("unref Base object.\n");
	g_object_unref (P3);
#endif
#if 0
    P = g_object_new (TEST_TYPE_CHILD1, NULL);  
    P1 = g_object_new (TEST_TYPE_CHILD1, NULL);  
	P2 = g_object_new (TEST_TYPE_CHILD1, NULL);
	printf("g_object_new P[%p] P1[%p] P2[%p].\n",P,P1,P2);
	g_signal_connect(G_OBJECT(P), "hello",
                             G_CALLBACK (test_hello_callback_connect1),NULL);
	g_signal_connect(G_OBJECT(P), "hello",
                             G_CALLBACK (test_hello_callback_connect),NULL);
	
	printf("start send signal.\n");
    send_hello_signal(P,1);
    // send_hello_signal(P1,2);
    // send_hello_signal(P2,3);
	// send_basehello_signal(P1,4);
	printf("end send signal.\n");
	g_object_unref (P);  
	g_object_unref(P1);
	g_object_unref(P2);
#endif
    return 0;  
}  