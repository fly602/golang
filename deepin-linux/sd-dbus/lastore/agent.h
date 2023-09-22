#ifndef __LASTORE_H__
#define __LASTORE_H__

#include <systemd/sd-bus.h>
#include <sys/statvfs.h>
#include <syslog.h>
#include <stdio.h>
#include <stdlib.h>

struct Agent
{
	sd_bus *session_bus;
    sd_bus *sys_bus;
    sd_bus_slot *slot;
    int is_wayland_session;
};

typedef struct Agent Agent;

#define OBJECT_PATH         "/com/deepin/lastore/agent"
#define INTERFACE_NAME      "com.deepin.lastore.Agent"

struct Agent * agent_init();
void agent_loop(struct Agent *agent);
uint64_t queryVFSAvailable(char *path);
#endif