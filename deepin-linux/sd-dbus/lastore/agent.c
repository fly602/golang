#include "agent.h"
#include "sd_bus_method.h"

struct Agent *agent = NULL;

uint64_t queryVFSAvailable(char *path)
{
    struct statvfs fs_info;
    int res = statvfs(path, &fs_info);
    if (!res)
    {
        LOG(LOG_ERR, "statvfs err");
        return -1;
    }
    return (uint64_t)fs_info.f_bavail * (uint64_t)fs_info.f_bsize;
}



// 添加dbus接口函数和属性PROPERTY
static const sd_bus_vtable agent_vtable[] = {
    SD_BUS_VTABLE_START(0),
    SD_BUS_METHOD("CloseNotification", "u", "", CloseNotification, SD_BUS_VTABLE_UNPRIVILEGED),
    SD_BUS_METHOD("GetManualProxy", "", "a{ss}", GetManualProxy, SD_BUS_VTABLE_UNPRIVILEGED),
    SD_BUS_METHOD("ReportLog", "s", "", ReportLog, SD_BUS_VTABLE_UNPRIVILEGED),
    SD_BUS_METHOD("SendNotify", "susssasa{sv}i", "u", SendNotify, SD_BUS_VTABLE_UNPRIVILEGED),
    SD_BUS_VTABLE_END};

// 初始化lastore
struct Agent *agent_init()
{
    agent = (struct Agent *)malloc(sizeof(struct Agent));
    memset(agent, 0, sizeof(struct Agent));

    if (strcmp(getenv("XDG_SESSION_TYPE"),"wayland") == 0){
        agent->is_wayland_session = true;
    }

    // 创建sd-bus
    int r = sd_bus_open_user(&agent->session_bus);
    if (r < 0)
    {
        LOG(LOG_ERR, "failed to connect to system bus: %s", strerror(-r));
        goto out;
    }

    r = sd_bus_open_system(&agent->sys_bus);
    if (r < 0) {
        LOG(LOG_ERR, "Failed to connect to system bus: %s", strerror(-r));
        goto out;
    }
    const char *unique_name = NULL;
    r = sd_bus_get_unique_name(agent->sys_bus, &unique_name);
    if (r < 0)
    {
        // 处理错误
        LOG(LOG_ERR, "unique name err");
        goto out;
    }
    LOG(LOG_INFO, "unique name: %s", unique_name);

    // 注册dbus函数
    r = sd_bus_add_object_vtable(agent->sys_bus,
                                 &agent->slot,
                                 OBJECT_PATH,    /* object path */
                                 INTERFACE_NAME, /* interface name */
                                 agent_vtable,
                                 agent);
    if (r < 0)
    {
        LOG(LOG_ERR, "failed to issue method call: %s", strerror(-r));
        goto out;
    }
    r = bus_syslastore_register_agent(agent,OBJECT_PATH);
out:
    if (r < 0)
    {
        LOG(LOG_ERR, "failed to register lastore agent: %s", strerror(-r));
        agent_close(agent);
    }
    return r < 0 ? NULL : agent ;
}

// 资源释放
void agent_close(Agent *agent){
    if (agent->slot)
        sd_bus_slot_unref(agent->slot);
    
    if (agent->session_bus)
        sd_bus_unref(agent->session_bus);

    if (agent->sys_bus)
        sd_bus_unref(agent->sys_bus);

    if (agent)
        free(agent);
}

// 启动dbus loop
void agent_loop(struct Agent *agent)
{
    int r = 0;
    for (;;)
    {
        /* Process requests */
        r = sd_bus_process(agent->sys_bus, NULL);
        if (r < 0)
        {
            LOG(LOG_ERR, "failed to process bus: %s", strerror(-r));
            goto finish;
        }
        if (r > 0) /* we processed a request, try to process another one, right-away */
            continue;

        /* Wait for the next request to process */
        r = sd_bus_wait(agent->sys_bus, (uint64_t)-1);
        if (r < 0)
        {
            LOG(LOG_ERR, "failed to wait on bus: %s", strerror(-r));
            goto finish;
        }
    }

finish:
    agent_close(agent);
}
