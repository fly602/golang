#ifndef __SD_BUS_METHOD__
#define __SD_BUS_METHOD__

#include "agent.h"

#define BUS_SYSLASTORE_NAME "com.deepin.lastore"
#define BUS_SYSLASTORE_PATH "/com/deepin/lastore"
#define BUS_SYSLASTORE_IF_NAME "com.deepin.lastore.Manager"

#define BUS_FREEDESKTOP_BUS_NAME "org.freedesktop.DBus"
#define BUS_FREEDESKTOP_BUS_PATH "/org/freedesktop/DBus"
#define BUS_FREEDESKTOP_BUS_IF_NAME "org.freedesktop.DBus"

#define BUS_DAEMON_EVENTLOG_NAME "com.deepin.daemon.EventLog"
#define BUS_DAEMON_EVENTLOG_PATH "/com/deepin/daemon/EventLog"
#define BUS_DAEMON_EVENTLOG_IF_NAME "com.deepin.daemon.EventLog"

#define BUS_OSD_NOTIFICATION_NAME "org.freedesktop.Notifications"
#define BUS_OSD_NOTIFICATION_PATH "/org/freedesktop/Notifications"
#define BUS_OSD_NOTIFICATION_IF_NAME "org.freedesktop.Notifications"

struct sd_bus_method
{
    uint32_t id;
    char *bus_name;
    char *bus_path;
    char *if_name;
    char *method_name;
    char *in_args;
};

typedef struct sd_bus_method sd_bus_method;

enum BUS_METHOD{
    BUS_METHOD_LOG_REPORT,
    BUS_METHOD_NOTIFY_CLOSE,
    BUS_METHOD_GET_CONNECTION_USER,
    BUS_METHOD_MAX,
};

// sd_bus接口调用的封装
#define bus_call_method(bus, method, reply,ret, ...) \
do{	\
	sd_bus_error error = SD_BUS_ERROR_NULL;	\
	int r = sd_bus_call_method(bus, \
							   method.bus_name, \
							   method.bus_path,	\
							   method.if_name, \
							   method.method_name,	\
							   &error,		\
							   reply,	\
							   method.in_args,\
							   ##__VA_ARGS__);\
	if (r < 0){\
		fprintf(stderr, "to here Failed to issue method call: %s,method: %s\n", error.message,method.method_name);\
	}\
	sd_bus_error_free(&error);\
	*ret = r;\
}while (0)

// system lastore RegisterAgent接口
int bus_syslastore_registerAgent(struct Agent *agent,char *path);

// 校验是否是系统调用
int check_caller_auth(sd_bus_message *m, void *userdata);
// dde-daemon reportlog接口
int bus_eventlog_reportlog(sd_bus_message *m, void *userdata);

// sd-bus接口
int CloseNotification(sd_bus_message *m, void *userdata, sd_bus_error *ret_error);
int GetManualProxy(sd_bus_message *m, void *userdata, sd_bus_error *ret_error);
int ReportLog(sd_bus_message *m, void *userdata, sd_bus_error *ret_error);
int SendNotify(sd_bus_message *m, void *userdata, sd_bus_error *ret_error);



#endif