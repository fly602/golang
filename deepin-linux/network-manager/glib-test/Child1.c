#include "Child1.h"


G_DEFINE_TYPE(TESTChild1,test_child1,TEST_TYPE_BASE);

typedef enum
{
	HELLO_SIGNAL = 0,
	LAST_SIGNAL
};

static int signals[LAST_SIGNAL];

static void test_child1_class_destroy(GObject *object)
{
	TESTBase *pTESTBase = NULL;

	pTESTBase = TEST_BASE(object);
	
	printf("test_child1_class_destroy.\n");
	
	G_OBJECT_CLASS(test_child1_parent_class)->dispose(object);
}

static void test_child1_class_finalize(GObject *object)
{
	TESTChild1 *pTESTChild1 = NULL;

	pTESTChild1 = TEST_BASE(object);
	
	printf("test_child1_class_finalize.\n");
	G_OBJECT_CLASS(test_child1_parent_class)->finalize(object);
	
    //GTK_OBJECT_CLASS (test_child1_parent_class)->destroy (object);
}


static void test_child1_init(TESTChild1 *pBase)
{
	printf("test_child1_parent_class [%x].\n",G_OBJECT_CLASS(test_child1_parent_class));
	printf("test_child1_init [%s].\n",pBase->parent.szName);
}
static void test_hello_callback(TESTChild1 *pBase,int i)
{
	printf("TESTChild1[%p] recv hello signal[%d].\n",pBase,i);
}

static void test_child1_class_init(TESTChild1Class *pBaseClass)
{
	printf("test_child1_class_init [%d].\n",pBaseClass->classparent.iAction);
	
	//GTK_OBJECT_CLASS(pBaseClass)->destroy = test_child1_class_destroy;
	G_OBJECT_CLASS(pBaseClass)->dispose = test_child1_class_destroy;
	G_OBJECT_CLASS(pBaseClass)->finalize = test_child1_class_finalize;
	pBaseClass->helloCb = test_hello_callback;

	signals[HELLO_SIGNAL] = g_signal_new ("hello",
                      G_TYPE_FROM_CLASS (pBaseClass),
                      G_SIGNAL_RUN_LAST,
                      G_STRUCT_OFFSET (TESTChild1Class, helloCb),
                      NULL,
                      NULL,
                      g_cclosure_marshal_VOID__INT,
                      G_TYPE_NONE, 1, G_TYPE_INT);
		
}

void send_hello_signal(TESTChild1 *pBase,int i)
{
	g_signal_emit (pBase, signals[HELLO_SIGNAL], 0,i);
}