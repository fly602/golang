#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <systemd/sd-bus.h>
#include <string.h>
#include <syslog.h>
#include "agent.h"
#include "log.h"

#define PROG_NAME "lastore-agent"

int main(int argc, char *argv[])
{
	// 初始化日志系统，指定程序名称和选项
	openlog(PROG_NAME, LOG_PID, LOG_USER);

	// 初始化lastore
	struct Agent *agent = agent_init();
	if (agent == NULL) {
		LOG(LOG_ERR, "Init %s err",PROG_NAME);
		return -1;
	}

	agent_loop(agent);
	return 0;
}