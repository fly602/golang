#include "sd_bus_method.h"

#define SD_BUS_ARG_INFO_FIELD(type_,contents_) \
	((type_info *) (&((const type_info) { \
		.type = type_, \
		.contents = contents_, \
	})))

#define SD_BUS_ARG_INFOS(...) \
	((type_info **) ((const type_info *[]) { \
		__VA_ARGS__ \
		NULL, \
	}))

sd_bus_method bus_methods[BUS_METHOD_MAX] = {
	// 注意，顺序同enum BUS_METHOD枚举中的顺序不要乱
	{
		BUS_METHOD_LOG_REPORT,
		BUS_DAEMON_EVENTLOG_NAME,
		BUS_DAEMON_EVENTLOG_PATH,
		BUS_DAEMON_EVENTLOG_IF_NAME,
		"ReportLog",
		"s",
		.types = SD_BUS_ARG_INFOS (
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
					),
	},
	{
		BUS_METHOD_NOTIFY_CLOSE,
		BUS_OSD_NOTIFICATION_NAME,
		BUS_OSD_NOTIFICATION_PATH,
		BUS_OSD_NOTIFICATION_IF_NAME,
		"CloseNotification",
		"s",
		.types = SD_BUS_ARG_INFOS (
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
					),
	},
	{
		BUS_METHOD_GET_CONNECTION_USER,
		BUS_FREEDESKTOP_BUS_NAME,
		BUS_FREEDESKTOP_BUS_PATH,
		BUS_FREEDESKTOP_BUS_IF_NAME,
		"GetConnectionUnixUser",
		"s",
		.types = SD_BUS_ARG_INFOS (
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
					),
	},
	{
		BUS_METHOD_NETWORK_GET_PROXYMETHOD,
		BUS_DAEMON_NETWORK_NAME,
		BUS_DAEMON_NETWORK_PATH,
		BUS_DAEMON_NETWORK_IF_NAME,
		"GetProxyMethod",
		"",
	},
	{
		BUS_METHOD_NETWORK_GET_PROXY,
		BUS_DAEMON_NETWORK_NAME,
		BUS_DAEMON_NETWORK_PATH,
		BUS_DAEMON_NETWORK_IF_NAME,
		"GetProxy",
		"s",
		.types = SD_BUS_ARG_INFOS (
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
					),
	},
	{
		BUS_METHOD_NETWORK_GET_PROXY_AUTH,
		BUS_DAEMON_NETWORK_NAME,
		BUS_DAEMON_NETWORK_PATH,
		BUS_DAEMON_NETWORK_IF_NAME,
		"GetProxyAuthentication",
		"s",
		.types = SD_BUS_ARG_INFOS (
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
					),
	},
	{
		BUS_METHOD_WM_ACTIVEWINDOW,
		BUS_DAEMON_WM_NAME,
		BUS_DAEMON_WM_PATH,
		BUS_DAEMON_WM_IF_NAME,
		"ActiveWindow",
		"",
	},
	{
		BUS_METHOD_NOTIFY_NOTIFY,
		BUS_OSD_NOTIFICATION_NAME,
		BUS_OSD_NOTIFICATION_PATH,
		BUS_OSD_NOTIFICATION_IF_NAME,
		"Notify",
		"susssasa{sv}i",
		.types = SD_BUS_ARG_INFOS (
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
						SD_BUS_ARG_INFO_FIELD ('u',      NULL),
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
						SD_BUS_ARG_INFO_FIELD ('s',      NULL),
						SD_BUS_ARG_INFO_FIELD ('a',      "s"),
						SD_BUS_ARG_INFO_FIELD ('a',      "{sv}"),
						SD_BUS_ARG_INFO_FIELD ('i',      NULL),
					),
	},
};

int sd_bus_read_dict(sd_bus_message *msg,GHashTable **map)
{
	sd_bus_error err = SD_BUS_ERROR_NULL;
	int error;

	*map = g_hash_table_new(g_str_hash, g_str_equal);
	error = sd_bus_message_enter_container(msg, SD_BUS_TYPE_ARRAY, "{sv}");
	if (error < 0)
			return error;

	while ((error = sd_bus_message_enter_container(msg, SD_BUS_TYPE_DICT_ENTRY, "sv")) > 0)
	{
		const char *key;
		const char *value;
		const char *contents;

		error = sd_bus_message_read_basic(msg, SD_BUS_TYPE_STRING, &key);
		if (error < 0)
				return error;

		error = sd_bus_message_peek_type(msg, NULL, &contents);
		if (error < 0)
				return error;

		error = sd_bus_message_enter_container(msg, SD_BUS_TYPE_VARIANT, contents);
		if (error < 0)
				return error;

		error = sd_bus_message_read_basic(msg, SD_BUS_TYPE_STRING, &value);
		if (error < 0)
				return error;
		g_hash_table_insert(*map, (gpointer)key, (gpointer)value);
		error = sd_bus_message_exit_container(msg);
		if (error < 0)
				return error;

		error = sd_bus_message_exit_container(msg);
		if (error < 0)
				return error;
	}
	if (error < 0)
			return error;

	error = sd_bus_message_exit_container(msg);
	if (error < 0)
			return error;

	if (err._need_free != 0)
	{
			LOG(LOG_DEBUG,"%d \n", error);
			LOG(LOG_DEBUG,"returned error: %s\n", err.message);
	}
	else
	{
	}
	sd_bus_error_free(&err);

	return 0;
}

int sd_bus_message_get_data(sd_bus_message *msg,...){
	va_list ap;
	va_start(ap, msg);
	int r = sd_bus_message_get_datav(msg,ap);
	va_end(ap);
	return r;
}

int sd_bus_message_get_datav(sd_bus_message *msg, va_list ap)
{
	char type;
	const char *contents;
	int r;
	for (;;)
	{
		r = sd_bus_message_peek_type(msg, &type, &contents);
		if (r < 0)
			return r;
		switch (type)
		{
		case SD_BUS_TYPE_STRING:
		{
			char **s = va_arg(ap, char **);
			r = sd_bus_message_read_basic(msg, type, s);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_BOOLEAN:
		{
			bool *b = va_arg(ap, bool *);

			r = sd_bus_message_read_basic(msg, type, b);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_INT64:
		{
			int64_t *i64 = va_arg(ap, int64_t *);

			r = sd_bus_message_read_basic(msg, type, i64);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_INT32:
		{
			int32_t *i = va_arg(ap, int32_t *);

			r = sd_bus_message_read_basic(msg, type, i);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_UINT32:
		{
			uint32_t *i = va_arg(ap, uint32_t *);

			r = sd_bus_message_read_basic(msg, type, i);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_OBJECT_PATH:
		{
			char **p = va_arg(ap, char **);

			r = sd_bus_message_read_basic(msg, type, p);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_ARRAY:
			if (strcmp(contents, "s") == 0)
			{
				char ***str = va_arg(ap, char ***);
				r = sd_bus_message_read_strv(msg, str);
				if (r < 0)
					goto finish;
			}
			else if (strcmp(contents, "{sv}") == 0)
			{
				GHashTable **map = va_arg(ap, GHashTable **);
				r = sd_bus_read_dict(msg, map);
				if (r < 0)
					goto finish;
			}

			break;
		default:
			LOG(LOG_DEBUG, "get type:%c contents:%s not define,TODO...", type, contents);
			goto finish;
		}
	}
finish:
	return r;
}

int sd_bus_set_dict(sd_bus_message *msg, GHashTable *map)
{
	// 打开一个 a{sv} 的容器
	gpointer key, value;
	GHashTableIter iter;

	g_hash_table_iter_init(&iter, map);
	int r = sd_bus_message_open_container(msg, 'a', "{sv}");
	if (r < 0)
	{
		LOG(LOG_ERR, "Failed to open container: %s\n", strerror(-r));
		goto finish;
	}

	while (g_hash_table_iter_next(&iter, &key, &value))
	{
		// 打开 sv 的容器
		r = sd_bus_message_open_container(msg, 'e', "sv");
		if (r < 0)
		{
			LOG(LOG_ERR, "Failed to open container: %s\n", strerror(-r));
			goto finish;
		}
		LOG(LOG_INFO, "get kv: %s: %s\n", (char *)key, (char *)value);
		r = sd_bus_message_append_basic(msg, 's', (char *)key);
		if (r < 0)
		{
			LOG(LOG_ERR, "Failed to apend kv to container: %s\n", strerror(-r));
			continue;
		}
		// 打开 v 的容器
		r = sd_bus_message_open_container(msg, 'v', "s");
		if (r < 0)
		{
			LOG(LOG_ERR, "Failed to open container: %s\n", strerror(-r));
			goto finish;
		}
		r = sd_bus_message_append_basic(msg, 's', (char *)value);
		if (r < 0)
		{
			LOG(LOG_ERR, "Failed to apend kv to container: %s\n", strerror(-r));
			continue;
		}
		// 关闭 v 的容器
		r = sd_bus_message_close_container(msg);
		if (r < 0)
		{
			LOG(LOG_ERR, "Failed to close container: %s\n", strerror(-r));
			goto finish;
		}
		// 关闭 sv 的容器
		r = sd_bus_message_close_container(msg);
		if (r < 0)
		{
			LOG(LOG_ERR, "Failed to close container: %s\n", strerror(-r));
			goto finish;
		}
	}

	// 关闭 a{sv} 的容器
	r = sd_bus_message_close_container(msg);
	if (r < 0)
	{
		LOG(LOG_ERR, "Failed to close container: %s\n", strerror(-r));
		goto finish;
	}
finish:
	return r;
}

int sd_bus_set_datav(sd_bus_message *msg,sd_bus_method *bus_method, va_list ap){
	int r = 0;
	int i = 0;
	char type;
	char *contents = NULL;
	if (bus_method->types == NULL){
		return 0;
	}
	for(;;){
		if (bus_method->types[i] == NULL){
			return 0;
		}
		type = bus_method->types[i]->type;
		LOG(LOG_DEBUG,"======>>>>>>>>get type:%c",type);
		switch (type)
		{
		
		case SD_BUS_TYPE_STRING:
		{
			char *s = va_arg(ap, char *);
			r = sd_bus_message_append_basic(msg, type, s);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_BOOLEAN:
		{
			int b = va_arg(ap, int);

			r = sd_bus_message_append_basic(msg, type, &b);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_INT64:
		{
			int64_t i64 = va_arg(ap, int64_t);

			r = sd_bus_message_append_basic(msg, type, &i64);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_INT32:
		{
			int32_t i = va_arg(ap, int32_t);

			r = sd_bus_message_append_basic(msg, type, &i);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_UINT32:
		{
			uint32_t i = va_arg(ap, uint32_t);

			r = sd_bus_message_append_basic(msg, type, &i);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_OBJECT_PATH:
		{
			char *p = va_arg(ap, char *);

			r = sd_bus_message_append_basic(msg, type, p);
			if (r < 0)
				goto finish;
			break;
		}

		case SD_BUS_TYPE_ARRAY:
			contents = bus_method->types[i]->contents;
			if (strcmp(contents, "s") == 0)
			{
				char **str = va_arg(ap, char **);
				r = sd_bus_message_append_strv(msg, str);
				if (r < 0)
					goto finish;
			}
			else if (contents[0] == '{')
			{
				GHashTable *map = va_arg(ap, GHashTable *);
				r = sd_bus_set_dict(msg, map);
				if (r < 0)
					goto finish;
			}

			break;
		default:
			LOG(LOG_DEBUG, "get type:%c contents:%s not define,TODO...", type, contents);
			goto finish;
		}
		i++;
	}
finish:
	return r;
}

int sd_bus_set_data(sd_bus_message *msg,sd_bus_method *bus_method,...){
	va_list ap;
	va_start(ap, bus_method);
	int r = sd_bus_set_datav(msg,bus_method,ap);
	va_end(ap);
	return r;
}

int bus_syslastore_register_agent(struct Agent *agent,char *path){
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
		LOG (LOG_ERR, "Failed to issue method call: %s\n", error.message);
		goto finish;
	}
finish:
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}

int bus_call_method(sd_bus *bus, sd_bus_method *bus_method, sd_bus_message **reply, ...)
{
	sd_bus_error error = SD_BUS_ERROR_NULL;
	// 需要重新构造dicts a{sv}
	sd_bus_message *msg = NULL;
	// gpointer key, value;
	int r = sd_bus_message_new_method_call(bus, &msg, bus_method->bus_name, bus_method->bus_path, bus_method->if_name, bus_method->method_name);
	if (r < 0)
	{
		LOG(LOG_ERR, "Failed to new mehod call: %s\n", strerror(-r));
		goto finish;
	}
	// 读取参数
	va_list ap;
	va_start(ap, reply);
	// int r = sd_bus_set_datav(msg,BUS_METHOD_NOTIFY_NOTIFY,app_name,replaces_id,app_icon,summary,body,actions_array,hints_dict,expire_timeout);
	r = sd_bus_set_datav(msg, bus_method, ap);
	if (r < 0)
	{
		va_end(ap);
		LOG(LOG_ERR, "Failed to set data: %s\n", strerror(-r));
		goto finish;
	}
	va_end(ap);
	// 调用方法
	r = sd_bus_call(bus, msg, 0, &error, reply);
	if (r < 0)
	{
		LOG(LOG_ERR, "Failed to method call: %s\n", strerror(-r));
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
	int r = 0;

	if (userdata == NULL){
		LOG (LOG_ERR, "userdata nil\n");
		return EXIT_FAILURE;
	}
	const char *sender = sd_bus_message_get_sender(m);
	if (sender == NULL){
		LOG(LOG_ERR, "sender nil\n");
		return EXIT_FAILURE;
	}
	
	/* Issue the method call and store the respons message in m */
	agent = (Agent *)userdata; 
	bus_call_method(agent->sys_bus,&bus_methods[BUS_METHOD_GET_CONNECTION_USER],&reply,sender);
	if (r < 0)
	{
		LOG (LOG_ERR,  "Failed to issue method call: %s\n", error.message);
		goto finish;
	}

	r = sd_bus_message_read(reply, "u", &uid);
	if (r < 0)
	{
		LOG (LOG_ERR,  "Failed to read method call reply: %s\n", error.message);
		goto finish;
	}
	LOG(LOG_INFO, "GetConnectionUnixUser uid: %d\n", uid);
	if (uid != rootUid) {
		LOG (LOG_ERR, "not allow %s call this method\n", sender);
		goto finish;
	}
finish:
	LOG(LOG_INFO, "GetConnectionUnixUser ret: %d\n", r);
	sd_bus_error_free(&error);
	return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}

int CloseNotification(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	struct Agent *agent = NULL;
	uint32_t id = 0;

	LOG(LOG_DEBUG,"CloseNotification");
	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}

	if (userdata == NULL){
		LOG (LOG_ERR, "userdata nil\n");
		return EXIT_FAILURE;
	}
	int r = sd_bus_message_read(m, "u", &id);
	if (r < 0)
	{
		LOG (LOG_ERR,  "Failed to read msg: %s\n", error.message);
		goto finish;
	}
	LOG(LOG_DEBUG,"ReportLog: %d",id);
	agent = (struct Agent *)userdata;
	/* Issue the method call and store the respons message in m */
	bus_call_method(agent->session_bus,&bus_methods[BUS_METHOD_NOTIFY_CLOSE],&reply,id);
	if (r < 0)
	{
		LOG (LOG_ERR, "Failed to issue method call: %s\n", error.message);
		goto finish;
	}
	return sd_bus_reply_method_return(m, NULL);
finish:
	sd_bus_error_free(&error);
	return r;
}

int GetManualProxy(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
    sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	struct Agent *agent = NULL;
	int r = 0;
	char *method = NULL;
	char key[256] = {0};
	char val[256] = {0};
	sd_bus_message *dict_array_msg = NULL;

	LOG(LOG_DEBUG,"GetManualProxy");
	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}

	if (userdata == NULL){
		LOG (LOG_ERR, "userdata nil\n");
		return EXIT_FAILURE;
	}

	agent = (struct Agent *)userdata;
	/* Issue the method call and store the respons message in m */
	bus_call_method(agent->session_bus,&bus_methods[BUS_METHOD_NETWORK_GET_PROXYMETHOD],&reply);
	if (r < 0 || reply == NULL)
	{
		LOG (LOG_ERR, "Failed to issue method call: %s\n", error.message);
		goto finish;
	}

	r = sd_bus_message_read(reply, "s", &method);
	LOG(LOG_DEBUG,"get porxy method: %s",method);
	if (r < 0 || method == NULL)
	{
		LOG (LOG_ERR,  "Failed to read msg: %s\n", error.message);
		goto finish;
	}

	if (strcmp(method,"manual") != 0){
		LOG (LOG_INFO,  "only support manual proxy\n");
		return sd_bus_error_setf(ret_error, SD_BUS_ERROR_FAILED, "only support manual proxy.");
	}
	char *proxy_types[] = {
		PROXY_TYPE_HTTP,PROXY_TYPE_HTTPS,PROXY_TYPE_FTP,PROXY_TYPE_SOCKS
	};

	r = sd_bus_message_new_method_return(m, &dict_array_msg);
    if (r < 0) {
        LOG (LOG_ERR, "Failed to create array: %s\n", strerror(-r));
        return r;
    }

	// 打开数组容器
	r = sd_bus_message_open_container(dict_array_msg,'a',"{ss}");
	if (r < 0 )
	{
		LOG(LOG_ERR, "Unable to open array container, %s\n", error.message);
		goto finish;
	}

	for (int i=0;i<sizeof(proxy_types) / sizeof(proxy_types[0]);i++){
		// dbus调用network getproxy
		bus_call_method(agent->session_bus,&bus_methods[BUS_METHOD_NETWORK_GET_PROXY],&reply,proxy_types[i]);
		if (r < 0 )
		{
			LOG (LOG_ERR,  "Failed to call method, %s\n", error.message);
			goto finish;
		}
		char *host = NULL;
		char *port = NULL;
		// 解析dbus调用结果
		r = sd_bus_message_read(reply, "ss", &host,&port);
		if (r < 0 )
		{
			LOG(LOG_ERR, "Failed to get reply, %s\n", error.message);
			continue;
		}
		// dbus调用network GetProxyAuthentication
		bus_call_method(agent->session_bus,&bus_methods[BUS_METHOD_NETWORK_GET_PROXY_AUTH],&reply,proxy_types[i]);
		if (r < 0 )
		{
			LOG (LOG_ERR,  "Failed to call method, %s\n", error.message);
			goto finish;
		}
		char *usr = NULL;
		char *pwd = NULL;
		int enable = 0;
		// 解析dbus调用结果
		r = sd_bus_message_read(reply, "ssb", &usr,&pwd,&enable);
		if (r < 0 )
		{
			LOG(LOG_ERR, "Failed to get reply, %s\n", error.message);
			continue;
		}
		// 添加键值对到 Dict 中
		memset(key,0,sizeof(key));
		if (strcmp(proxy_types[i],PROXY_TYPE_SOCKS) == 0){
			sprintf(key,"%s", proxy_types[i]);
		} else {
			sprintf(key,"%s_proxy", proxy_types[i]);
		}
		// 添加键和值到字典：key
		memset(val,0,sizeof(val));
		if (enable) {
			sprintf(val, "%s://%s:%s@%s:%s", PROXY_TYPE_HTTP, usr, pwd, host, port);
		} else {
			sprintf(val, "%s://%s:%s", PROXY_TYPE_HTTP, host, port);
		}
		// 添加键和值到字典：value
		r = sd_bus_message_append(dict_array_msg, "{ss}", key,val);
		if (r < 0) {
			LOG (LOG_ERR, "Failed to append value1: %s\n", error.message);
			continue;;
		}
	}
	// 关闭数组容器
	r = sd_bus_message_close_container(dict_array_msg);
	if (r < 0) {
			LOG (LOG_ERR, "Failed to append value1: %s\n", error.message);
			goto finish;
	}
	// 响应成功，并将 a{ss} 数据结构作为返回值
    return sd_bus_send(NULL, dict_array_msg,NULL);
finish:
	sd_bus_error_free(&error);
	return r;
}

int ReportLog(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	Agent *agent = NULL;
	char *msg = NULL;

	if (userdata == NULL){
		LOG (LOG_ERR, "userdata nil\n");
		return EXIT_FAILURE;
	}

	// 读取入参
	int r = sd_bus_message_read(m, "s", &msg);
	if (r < 0 || msg == NULL)
	{
		LOG (LOG_ERR,  "Failed to read msg: %s\n", error.message);
		goto finish;
	}
	LOG(LOG_DEBUG,"report log, msg: %s",msg);

	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}
	
	agent = (Agent *)userdata;
	/* Issue the method call and store the respons message in m */
	bus_call_method(agent->session_bus,&bus_methods[BUS_METHOD_LOG_REPORT],&reply,msg);
	if (r == EXIT_FAILURE){
		r = -1;
	}
	r = sd_bus_reply_method_return(m, NULL);
finish:
	sd_bus_error_free(&error);
	return r;
}

int SendNotify(sd_bus_message *m, void *userdata, sd_bus_error *ret_error)
{
	sd_bus_error error = SD_BUS_ERROR_NULL;
	sd_bus_message *reply = NULL;
	Agent *agent = NULL;

	// 接口入参
	char *app_name = NULL;
	uint32_t replaces_id = 0;
	const char *app_icon = NULL,*summary = NULL,*body = NULL;
	char **actions_array = NULL;
	GHashTable *hints_dict = NULL;
	int32_t expire_timeout;

	LOG(LOG_DEBUG, "SendNotify");

	if (userdata == NULL)
	{
		LOG(LOG_ERR, "userdata nil\n");
		return EXIT_FAILURE;
	}
	agent = (Agent *)userdata;

	if (check_caller_auth(m,userdata) != EXIT_SUCCESS){
		return EXIT_FAILURE;
	}
	int r = sd_bus_message_get_data(m,&app_name,&replaces_id,&app_icon,&summary,&body,&actions_array,&hints_dict,&expire_timeout);
	if (r < 0) {
        fprintf(stderr, "Failed to get data: %s\n", strerror(-r));
        goto finish;
    }
	LOG(LOG_INFO,"receive notify from lastore daemon, app name: %s", app_name);

	int need_send = 1;
	if (strcmp (app_name,UPDATE_NOTIFY_SHOW_OPTIONAL) == 0){
		memset(app_name,0,strlen(app_name));
		strcpy(app_name,UPDATE_NOTIFY_SHOW);
		// 只有当控制中心获取焦点,且控制中心当前为更新模块时,不发通知
		if (agent->is_wayland_session) {
			bus_call_method(agent->session_bus,&bus_methods[BUS_METHOD_WM_ACTIVEWINDOW],&reply);
			if (r < 0)
			{
				LOG (LOG_ERR,  "Failed to call method: %s\n", error.message);
				goto finish;
			}
			uint32_t win_id = 0;
			int r = sd_bus_message_read(m, "u", &win_id);
			if (r < 0)
			{
				LOG (LOG_ERR,  "Failed to read msg: %s\n", error.message);
				goto finish;
			}
			char win_path[128] = {0};
			sprintf(win_path, "%s_%d", BUS_DAEMON_WM_WININFO_PATH, win_id);
			sd_bus_method bus_method = {
				-1,
				BUS_DAEMON_WM_NAME,
				win_path,
				BUS_DAEMON_WM_WININFO_IF_NAME,
				"AppId",
				""};
			r = bus_call_method(agent->session_bus, 
							   &bus_method,		
							   &reply);
			if (r < 0){
				LOG (LOG_ERR, "to here Failed to issue method call: %s,method: %s\n", error.message,"AppId");
				goto finish;
			}
			char *win_name = NULL;
			r = sd_bus_message_read(reply, "s", &win_name);
			if (r < 0)
			{
				LOG (LOG_ERR,  "Failed to read msg: %s\n", error.message);
				goto finish;
			}
			if (strstr(win_name,"dde-control-center") != NULL){
				// 焦点在控制中心上,需要判断是否为更新模块
				char *cur_mod = NULL;
				r = sd_bus_get_property_string(agent->session_bus,
								BUS_CONTROL_CENTER_NAME,
								BUS_CONTROL_CENTER_PATH,
								BUS_CONTROL_CENTER_IF_NAME,
								"CurrentModule",
								&error,
								&cur_mod);
				if (r < 0)
				{
					LOG (LOG_ERR, "to here Failed to issue get property: %s,property: %s\n", error.message,"CurrentModule");
					goto finish;
				}
				if (strcmp(cur_mod,"update") == 0){
					LOG(LOG_INFO,"update module of dde-control-center is in the foreground, don't need send notify");
					need_send = 0;
				}
			} else if (strstr(win_name,"dde-lock") != NULL){
				// 前台应用在模态更新界面时,不发送通知(TODO: 如果后台更新时发生了锁屏，需要增加判断是否发通知)
				need_send = 0;
			}
		} else {
			const char *command = "xprop -id $(xprop -root _NET_ACTIVE_WINDOW | cut -d ' ' -f 5) WM_CLASS";
			char buffer[1024];
			// 使用 popen 执行外部命令并获取输出
			FILE *fp = popen(command, "r");
			if (fp == NULL)
			{
				LOG (LOG_ERR, "Failed to run command: %s\n", error.message);
				goto finish;
			}

			// 读取命令输出到缓冲区
			if (fgets(buffer, sizeof(buffer), fp) != NULL)
			{
				// 检查输出中是否包含 "dde-control-center"
				if (strstr(buffer, "dde-control-center") != NULL)
				{
					// 焦点在控制中心上,需要判断是否为更新模块
					char *cur_mod = NULL;
					r = sd_bus_get_property_string(agent->session_bus,
												   BUS_CONTROL_CENTER_NAME,
												   BUS_CONTROL_CENTER_PATH,
												   BUS_CONTROL_CENTER_IF_NAME,
												   "CurrentModule",
												   &error,
												   &cur_mod);
					if (r < 0)
					{
						LOG (LOG_ERR, "to here Failed to issue get property: %s,property: %s\n", error.message, "CurrentModule");
						pclose(fp);
						goto finish;
					}
					if (strcmp(cur_mod, "update") == 0)
					{
						LOG(LOG_INFO, "update module of dde-control-center is in the foreground, don't need send notify");
						need_send = 0;
					}
				}
				else if (strstr(buffer, "dde-lock") != NULL)
				{
					// 前台应用在模态更新界面时,不发送通知(TODO: 如果后台更新时发生了锁屏，需要增加判断是否发通知)
					need_send = 0;
				}
			}

			// 关闭文件指针
			pclose(fp);
		}
	}
	uint32_t id = 0;
	if (need_send) {
		bus_call_method(agent->session_bus,
						&bus_methods[BUS_METHOD_NOTIFY_NOTIFY],
						&reply,
						app_name,
						replaces_id,
						app_icon,summary,
						body,
						actions_array,
						hints_dict,
						expire_timeout);
		r = sd_bus_message_read(reply, "u", &id);
		if (r < 0)
		{
			LOG (LOG_ERR,  "Failed to read msg: %s\n", strerror(-r));
			goto finish;
		}
		r =  sd_bus_reply_method_return(m,"u",id);
	}
finish:
	if (hints_dict != NULL) {
		g_hash_table_destroy(hints_dict);
	}
	sd_bus_error_free(&error);
	return r;
}