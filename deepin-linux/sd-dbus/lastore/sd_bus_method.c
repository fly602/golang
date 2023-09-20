#include "sd_bus_method.h"
#include <stdarg.h>

static sd_bus_method bus_methods[BUS_METHOD_MAX] = {
	// 注意，顺序同enum BUS_METHOD枚举中的顺序不要乱
	{
		BUS_METHOD_LOG_REPORT,
		BUS_DAEMON_EVENTLOG_NAME,
		BUS_DAEMON_EVENTLOG_PATH,
		BUS_DAEMON_EVENTLOG_IF_NAME,
		"ReportLog",
		"s",
	},
	{
		BUS_METHOD_NOTIFY_CLOSE,
		BUS_OSD_NOTIFICATION_NAME,
		BUS_OSD_NOTIFICATION_PATH,
		BUS_OSD_NOTIFICATION_IF_NAME,
		"CloseNotification",
		"s",
	},
	{
		BUS_METHOD_GET_CONNECTION_USER,
		BUS_FREEDESKTOP_BUS_NAME,
		BUS_FREEDESKTOP_BUS_PATH,
		BUS_FREEDESKTOP_BUS_IF_NAME,
		"GetConnectionUnixUser",
		"s",
	},
};

int bus_syslastore_registerAgent(struct Agent *agent,char *path){
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *m = NULL;
	/* Issue the method call and store the respons message in m */
	int r = sd_bus_call_method(agent->sys_bus,
							   BUS_SYSLASTORE_NAME,	/* service to contact */
							   BUS_SYSLASTORE_PATH,	/* object path */
							   BUS_SYSLASTORE_IF_NAME, /* interface name */
							   "RegisterAgent",			/* method name */
							   &error,					/* object to return error in */
							   &m,						/* return message on success */
							   "s",
							   path); /* second argument */
	if (r < 0)
	{
		fprintf(stderr, "Failed to issue method call: %s\n", error.message);
		goto finish;
	}
finish:
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}

int check_caller_auth(sd_bus_message *m, void *userdata){
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	uint32_t uid = 0;
	struct Agent *agent = NULL;
	const u_int32_t rootUid = 0;

	if (userdata == NULL){
		fprintf(stderr, "userdata nil\n");
		return EXIT_FAILURE;
	}
	const char *sender = sd_bus_message_get_sender(m);
	if (sender == NULL){
		syslog(LOG_ERR, "sender nil\n");
		return EXIT_FAILURE;
	}
	
	/* Issue the method call and store the respons message in m */
	agent = (struct Agent *)userdata; 
	int r = bus_session_call_method(agent->sys_bus,&bus_methods[BUS_METHOD_GET_CONNECTION_USER],&reply,sender);
	if (r < 0)
	{
		fprintf(stderr,  "Failed to issue method call: %s\n", error.message);
		goto finish;
	}

	r = sd_bus_message_read(reply, "u", &uid);
	if (r < 0)
	{
		fprintf(stderr,  "Failed to read method call reply: %s\n", error.message);
		goto finish;
	}
	syslog(LOG_INFO, "GetConnectionUnixUser uid: %d\n", uid);
	if (uid != rootUid) {
		fprintf(stderr, "not allow %s call this method\n", sender);
		goto finish;
	}
finish:
	syslog(LOG_INFO, "GetConnectionUnixUser ret: %d\n", r);
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}

int bus_session_call_method(sd_bus *bus, sd_bus_method *method,sd_bus_message **reply,...){
	sd_bus_error error = SD_BUS_ERROR_NULL;

	va_list args;
    va_start(args, reply);
		printf("method: id=%d, name=%s, path=%s,if_name=%s, method=%s in_args=%s\n",
								method->id,
								method->bus_name,	
							   method->bus_path,	
							   method->if_name, 
							   method->method_name,
							   method->in_args);
	/* Issue the method call and store the respons message in m */
	int r = sd_bus_call_method(bus,
							   method->bus_name,	
							   method->bus_path,	
							   method->if_name, 
							   method->method_name,			
							   &error,					
							   reply,						
							   method->in_args,
							   args);
	va_end(args);
	if (r < 0)
	{
		fprintf(stderr, "to here Failed to issue method call: %s,method: %s\n", error.message,method->method_name);
		goto finish;
	}
finish:
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}

int bus_eventlog_reportlog(sd_bus_message *m, void *userdata){
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	Agent *agent = NULL;
	char *msg = NULL;

	if (userdata == NULL){
		fprintf(stderr, "userdata nil\n");
		return EXIT_FAILURE;
	}
	int r = sd_bus_message_read(m, "s", &msg);
	if (r < 0)
	{
		printf("%s %d:=====to here\n",__FILE__,__LINE__);
		fprintf(stderr,  "Failed to read msg: %s\n", error.message);
		goto finish;
	}
	syslog(LOG_DEBUG,"ReportLog: %s",msg);
	agent = (Agent *)userdata;
	/* Issue the method call and store the respons message in m */

	r = bus_session_call_method(agent->session_bus,&bus_methods[BUS_METHOD_LOG_REPORT],&reply,msg);
	if (r == EXIT_FAILURE){
		printf("%s %d:=====to here\n",__FILE__,__LINE__);
		r = -1;
	}
	printf("%s %d:=====to here\n",__FILE__,__LINE__);
finish:
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}

int bus_notification_close(sd_bus_message *m, void *userdata){
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	struct Agent *agent = NULL;
	uint32_t id = 0;

	if (userdata == NULL){
		fprintf(stderr, "userdata nil\n");
		return EXIT_FAILURE;
	}
	int r = sd_bus_message_read(m, "u", &id);
	if (r < 0)
	{
		fprintf(stderr,  "Failed to read msg: %s\n", error.message);
		goto finish;
	}
	syslog(LOG_DEBUG,"ReportLog: %d",id);
	agent = (struct Agent *)userdata;
	/* Issue the method call and store the respons message in m */
	r = bus_session_call_method(agent->session_bus,&bus_methods[BUS_METHOD_NOTIFY_CLOSE],&reply,id);
	if (r < 0)
	{
		fprintf(stderr, "Failed to issue method call: %s\n", error.message);
		goto finish;
	}
finish:
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}



int CloseNotification(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	syslog(LOG_DEBUG,"CloseNotification");
	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}
    return bus_notification_close(m,userdata);
}

int GetManualProxy(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	syslog(LOG_DEBUG,"GetManualProxy");
	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}
    return bus_notification_close(m,userdata);
}

int ReportLog(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	syslog(LOG_DEBUG,"ReportLog");
	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}
    return bus_eventlog_reportlog(m,userdata);
}

int SendNotify(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	syslog(LOG_DEBUG,"SendNotify");
	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}
    return sd_bus_reply_method_return(m, "b", queryVFSAvailable("/") > 1024 * 1024 * 10);
}