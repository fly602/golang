#ifndef _BASE_H_
#define _BASE_H_
//#include <gtk/gtk.h>

#include <glib-object.h>
#include <stdio.h>

typedef struct _TESTBase 	 TESTBase;
typedef struct _TESTBaseClass TESTBaseClass;

#define TEST_TYPE_BASE    (test_base_get_type())

#define TEST_BASE(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), TEST_TYPE_BASE, TESTBase))

#define TEST_IS_BASE(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TEST_TYPE_BASE))

#define TEST_BASE_CLASS(klass) (G_TYPE_CHECK_CLASS_CAST ((klass), TEST_TYPE_BASE, TESTBaseClass))

#define TEST_IS_BASE_CLASS(klass) (G_TYPE_CHECK_CLASS_TYPE ((klass), TEST_TYPE_BASE))

#define TEST_BASE_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), TEST_TYPE_BASE, TESTBaseClass))


struct _TESTBase
{
	GObject parent;
	char 	szName[32];
};

struct _TESTBaseClass
{
	GObjectClass classparent;
	int 		 iAction;
	void (*basehello)(TESTBase *obj,int i);
};


void send_basehello_signal(TESTBase *pBase,int i);


#endif