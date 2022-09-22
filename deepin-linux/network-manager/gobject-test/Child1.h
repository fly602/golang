#ifndef _CHILD1_H_
#define _CHILD1_H_
#include "Base.h"

typedef struct _TESTChild1 	 TESTChild1;
typedef struct _TESTChild1Class TESTChild1Class;

#define TEST_TYPE_CHILD1    (test_child1_get_type())
//
#define TEST_CHILD1(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), TEST_TYPE_CHILD1, TESTChild1))
//
#define TEST_IS_CHILD1(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TEST_TYPE_CHILD1))
//
#define TEST_CHILD1_CLASS(klass) (G_TYPE_CHECK_CLASS_CAST ((klass), TEST_TYPE_CHILD1, TESTChild1Class))
//
#define TEST_IS_CHILD1_CLASS(klass) (G_TYPE_CHECK_CLASS_TYPE ((klass), TEST_TYPE_CHILD1))
//
#define TEST_CHILD1_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), TEST_TYPE_CHILD1, TESTChild1Class))

struct _TESTChild1
{
	TESTBase parent;
	char 	 szChild1Name[32];
};

struct _TESTChild1Class
{
	TESTBaseClass classparent;
	int 		  iChild1Action;
	void		  (*helloCb)(TESTChild1 *pBase,int i);
};


void send_hello_signal(TESTChild1 *pBase,int i);

#endif