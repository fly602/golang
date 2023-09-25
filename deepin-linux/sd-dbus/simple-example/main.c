#include <stdio.h>
#include <stdlib.h>
#include <systemd/sd-bus.h>
#include <stdbool.h>
#include <string.h>

int bus_print_property(const char *name, sd_bus_message *property)
{
        char type;
        const char *contents;
        int r;

        r = sd_bus_message_peek_type(property, &type, &contents);
        if (r < 0)
                return r;

        switch (type)
        {

        case SD_BUS_TYPE_STRING:
        {
                const char *s;

                r = sd_bus_message_read_basic(property, type, &s);
                if (r < 0)
                        return r;

                printf("%s=%s\n", name, s);

                return 1;
        }

        case SD_BUS_TYPE_BOOLEAN:
        {
                bool b;

                r = sd_bus_message_read_basic(property, type, &b);
                if (r < 0)
                        return r;

                printf("%s\n", name);

                return 1;
        }

        case SD_BUS_TYPE_INT64:
        {
                int64_t i64;

                r = sd_bus_message_read_basic(property, type, &i64);
                if (r < 0)
                        return r;

                printf("%s=%i\n", name, (int)i64);
                return 1;
        }

        case SD_BUS_TYPE_INT32:
        {
                int32_t i;

                r = sd_bus_message_read_basic(property, type, &i);
                if (r < 0)
                        return r;

                printf("%s=%i\n", name, (int)i);
                return 1;
        }

        case SD_BUS_TYPE_OBJECT_PATH:
        {
                const char *p;

                r = sd_bus_message_read_basic(property, type, &p);
                if (r < 0)
                        return r;

                printf("%s=%s\n", name, p);

                return 1;
        }

        case SD_BUS_TYPE_DOUBLE:
        {
                double d;

                r = sd_bus_message_read_basic(property, type, &d);
                if (r < 0)
                        return r;

                printf("%s=%g\n", name, d);
                return 1;
        }

        case SD_BUS_TYPE_ARRAY:
                if (strcmp(contents, "s") == 0)
                {
                        bool first = true;
                        const char *str;

                        r = sd_bus_message_enter_container(property, SD_BUS_TYPE_ARRAY, contents);
                        if (r < 0)
                                return r;

                        while ((r = sd_bus_message_read_basic(property, SD_BUS_TYPE_STRING, &str)) > 0)
                        {
                                if (first)
                                        printf("%s=", name);

                                printf("%s%s", first ? "" : " ", str);

                                first = false;
                        }
                        if (r < 0)
                                return r;

                        if (first)
                                printf("%s=", name);
                        if (!first)
                                puts("");

                        r = sd_bus_message_exit_container(property);
                        if (r < 0)
                                return r;

                        return 1;
                }
                else
                {
                        printf("array unreadable");
                        return 0;
                }

                break;
        }

        return 0;
}

int main()
{
        sd_bus *bus = NULL;
        sd_bus_error err = SD_BUS_ERROR_NULL;
        sd_bus_message *msg = NULL;
        int error;

        sd_bus_default_user(&bus);

        sd_bus_get_property(bus,
                            "org.mpris.MediaPlayer2.plasma-browser-integration",
                            "/org/mpris/MediaPlayer2",
                            "org.mpris.MediaPlayer2.Player",
                            "Metadata",
                            &err, &msg, "a{sv}");
        error = sd_bus_message_enter_container(msg, SD_BUS_TYPE_ARRAY, "{sv}");
        if (error < 0)
                return error;

        while ((error = sd_bus_message_enter_container(msg, SD_BUS_TYPE_DICT_ENTRY, "sv")) > 0)
        {
                const char *name;
                const char *contents;

                error = sd_bus_message_read_basic(msg, SD_BUS_TYPE_STRING, &name);
                if (error < 0)
                        return error;

                error = sd_bus_message_peek_type(msg, NULL, &contents);
                if (error < 0)
                        return error;

                error = sd_bus_message_enter_container(msg, SD_BUS_TYPE_VARIANT, contents);
                if (error < 0)
                        return error;

                error = bus_print_property(name, msg);
                if (error < 0)
                        return error;
                if (error == 0)
                {

                        printf("%s=[unprintable]\n", name);
                        /* skip what we didn't read */
                        error = sd_bus_message_skip(msg, contents);
                        if (error < 0)
                                return error;
                }

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
                printf("%d \n", error);
                printf("returned error: %s\n", err.message);
        }
        else
        {
        }
        sd_bus_error_free(&err);
        sd_bus_unref(bus);

        return 0;
}