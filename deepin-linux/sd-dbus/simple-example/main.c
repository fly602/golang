#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <systemd/sd-bus.h>
#include<string.h>
#include <gio/gio.h>
#include <syslog.h>

#define OBJECT_PATH "/net/poettering/Calculator"
#define INTERFACE_NAME "net.poettering.Calculator"
#define GS_SCHEMA_POWER     "com.deepin.dde.power"
#define GS_UPDATING_LOW_POWER_PERCENT "low-power-percent-in-updating-notify"

struct UserData
{
    /* data */
    sd_bus *bus;
    char *name;
    int age;
        GSettings *gs_power;
    double low_power_percent ;
};


static int method_multiply(sd_bus_message *m, void *userdata, sd_bus_error *ret_error) {
        int64_t x, y;
        int r;

        /* Read the parameters */
        r = sd_bus_message_read(m, "xx", &x, &y);
        if (r < 0) {
                fprintf(stderr, "Failed to parse parameters: %s\n", strerror(-r));
                return r;
        }

        /* Reply with the response */
        return sd_bus_reply_method_return(m, "x", x * y);
}

static int method_divide(sd_bus_message *m, void *userdata, sd_bus_error *ret_error) {
        int64_t x, y;
        int r;

        /* Read the parameters */
        r = sd_bus_message_read(m, "xx", &x, &y);
        if (r < 0) {
                fprintf(stderr, "Failed to parse parameters: %s\n", strerror(-r));
                return r;
        }

        /* Return an error on division by zero */
        if (y == 0) {
                sd_bus_error_set_const(ret_error, "net.poettering.DivisionByZero", "Sorry, can't allow division by zero.");
                return -EINVAL;
        }

        return sd_bus_reply_method_return(m, "x", x / y);
}

static int get_property(sd_bus *bus, const char *path, const char *interface, const char *property, sd_bus_message *reply, void *userdata, sd_bus_error *ret_error){
        struct UserData *u = userdata;
        int r;


        if (strcmp("name",property)==0){
            printf("property get for %s called, returning \"%s\".\n", property, u->name);
            r = sd_bus_message_append(reply, "s", u->name);
        } else if (strcmp("age",property)==0){
            printf("property get for %s called, returning \"%s\".\n", property, u->name);
            r = sd_bus_message_append(reply, "i", u->age);
        }
        // assert_se(r >= 0);

        return 1;
}

static int set_property(sd_bus *bus, const char *path, const char *interface, const char *property, sd_bus_message *value, void *userdata, sd_bus_error *ret_error){
        struct UserData *u = userdata;
        const char *s;
        int r;

        printf("property set for %s called\n", property);

        if (strcmp("name",property)==0){
            char *n = NULL;
            r = sd_bus_message_read(value, "s", &s);
            n = strdup(s);
            // assert_se(n);
            free(u->name);
            u->name = n;
        } else if (strcmp("age",property)==0){
            int age = 0;
            r = sd_bus_message_read(value, "i", &age);
            // assert_se(r >= 0);
            u->age = age;
        }

        return 1;
}

void hand_gs_changed(GSettings *settings, const gchar *key, gpointer user_data)
{
    struct UserData *u = user_data;
    syslog(LOG_INFO,"low power percent changed:%s",key);
    u->low_power_percent = (double)g_settings_get_double(u->gs_power, GS_UPDATING_LOW_POWER_PERCENT);
}

/* The vtable of our little object, implements the net.poettering.Calculator interface */
static const sd_bus_vtable calculator_vtable[] = {
        SD_BUS_VTABLE_START(0),
        // SD_BUS_METHOD参数说明：
        // arg0: 函数名称;  arg1： 入参类型及个数，例如："s"，有一个字符串类型入参，"ss"有两个；
        // arg2：出参，用法同入参；arg3: 接口的实现函数； arg4： flags
        SD_BUS_METHOD("Multiply", "xx", "x", method_multiply, SD_BUS_VTABLE_UNPRIVILEGED),
        SD_BUS_METHOD("Divide",   "xx", "x", method_divide,   SD_BUS_VTABLE_UNPRIVILEGED),
        SD_BUS_WRITABLE_PROPERTY("name","s",get_property,set_property,0,0),
        SD_BUS_WRITABLE_PROPERTY("age","i",get_property,set_property,0,0),
        SD_BUS_VTABLE_END
};

int main(int argc, char *argv[]) {
        sd_bus_slot *slot = NULL;
        sd_bus *bus = NULL;
        int r;
        struct UserData *ud = (struct UserData *)malloc(sizeof(struct UserData));
        memset(ud, 0, sizeof(struct UserData));
        ud->name = "fuleyi";
        ud->age = 123;
        /* Connect to the user bus this time */
        r = sd_bus_open_user(&bus);
        if (r < 0) {
                fprintf(stderr, "Failed to connect to system bus: %s\n", strerror(-r));
                goto finish;
        }
        const char *unique_name = NULL;
        r = sd_bus_get_unique_name(bus, &unique_name);
        if (r < 0) {
            // 处理错误
            printf("Unique name err\n");
        }
        printf("Unique name: %s\n", unique_name);

        ud->bus = bus;
        /* Install the object */
        r = sd_bus_add_object_vtable(bus,
                                     &slot,
                                     OBJECT_PATH,  /* object path */
                                     INTERFACE_NAME,   /* interface name */
                                     calculator_vtable,
                                     ud);
        if (r < 0) {
                fprintf(stderr, "Failed to issue method call: %s\n", strerror(-r));
                goto finish;
        }

        /* Take a well-known service name so that clients can find us */
        // r = sd_bus_request_name(bus, "net.poettering.Calculator", 0);
        // if (r < 0) {
        //         fprintf(stderr, "Failed to acquire service name: %s\n", strerror(-r));
        //         goto finish;
        // }

        for (;;) {
                /* Process requests */
                r = sd_bus_process(bus, NULL);
                if (r < 0) {
                        fprintf(stderr, "Failed to process bus: %s\n", strerror(-r));
                        goto finish;
                }
                if (r > 0) /* we processed a request, try to process another one, right-away */
                        continue;

                /* Wait for the next request to process */
                r = sd_bus_wait(bus, (uint64_t) -1);
                if (r < 0) {
                        fprintf(stderr, "Failed to wait on bus: %s\n", strerror(-r));
                        goto finish;
                }
        }

        // 监听gsetting
        GSettings *gs = g_settings_new(GS_SCHEMA_POWER);
        ud->gs_power = gs;
        ud->low_power_percent = (double)g_settings_get_double(gs, GS_UPDATING_LOW_POWER_PERCENT);
        char key[128] = {};
        sprintf(key,"changed::%s",GS_UPDATING_LOW_POWER_PERCENT);
        g_signal_connect(gs,key,G_CALLBACK(hand_gs_changed), (void *)ud);

        GMainLoop *loop = g_main_loop_new(NULL, FALSE);
        g_main_loop_run(loop);

finish:
        sd_bus_slot_unref(slot);
        sd_bus_unref(bus);

        return r < 0 ? EXIT_FAILURE : EXIT_SUCCESS;
}