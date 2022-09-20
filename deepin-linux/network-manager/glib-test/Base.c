#include "Base.h"
#include <string.h>

typedef enum
{
	HELLO_BASE_SIGNAL = 0,
	LAST_BASE_SIGNAL
};

static int basesignals[LAST_BASE_SIGNAL];

G_DEFINE_TYPE(TESTBase,test_base,G_TYPE_OBJECT);

static void test_base_init(TESTBase *pBase)
{
	memcpy(pBase->szName,"baseJob",sizeof(pBase->szName));
	printf("test_base_init.\n");
	
}

static void test_base_class_destroy(GObject *object)
{
	TESTBase *pTESTBase = NULL;

	pTESTBase = TEST_BASE(object);
	G_OBJECT_CLASS(test_base_parent_class)->dispose(object);
	printf("test_base_class_destroy.\n");
}

static void test_base_class_finalize(GObject *object)
{
	TESTBase *pTESTBase = NULL;

	pTESTBase = TEST_BASE(object);
	
	printf("test_base_class_finalize.\n");
}

static void test_base_hello_signal(TESTBase *pBase,int i)
{
	
	printf("test_base_hello_signal[%p] [%d].\n",pBase,i);
}

static void test_base_class_init(TESTBaseClass *pBaseClass)
{
	pBaseClass->iAction = 1;
	
	G_OBJECT_CLASS(pBaseClass)->finalize = test_base_class_finalize;
	G_OBJECT_CLASS(pBaseClass)->dispose = test_base_class_destroy;
	pBaseClass->basehello  = test_base_hello_signal;
	basesignals[HELLO_BASE_SIGNAL] = g_signal_new ("basehello",
                      G_TYPE_FROM_CLASS (pBaseClass),
                      G_SIGNAL_RUN_LAST,
					  G_STRUCT_OFFSET (TESTBaseClass, basehello),
                      NULL,
                      NULL,
                      g_cclosure_marshal_VOID__INT,
                      G_TYPE_NONE, 1, G_TYPE_INT);
	printf("test_base_class_init.\n");
}
void send_basehello_signal(TESTBase *pBase,int i)
{
	g_signal_emit (pBase, basesignals[HELLO_BASE_SIGNAL], 0,i);
}