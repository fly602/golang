#include <X11/Xlib.h>
#include <X11/extensions/scrnsaver.h>
#include <stdio.h>


int main() {
    Display *display = XOpenDisplay(NULL);
    if (!display) {
        fprintf(stderr, "Unable to open display\n");
        return 1;
    }

    Window root = DefaultRootWindow(display);
    int event_base, error_base;

    if (XScreenSaverQueryExtension(display, &event_base, &error_base)) {
        // 设置事件掩码
        XScreenSaverSelectInput(display, root, ScreenSaverNotifyMask | ScreenSaverCycleMask);

        XEvent ev;
        while (1) {
            XNextEvent(display, &ev);
            printf("Screen saver recv event:%d\n",ev.type);
            switch (ev.type)
            {
            case ScreenSaverNotify:
                printf("Screen saver state changed\n");
                break;
            case ScreenSaverCycleMask:
                printf("Screen saver Cycle changed\n");
                break;
            default:
                printf("Screen saver unknown event:%d\n",ev.type);
                break;
            }
        }

    } else {
        fprintf(stderr, "Screensaver extension not available\n");
    }

    XCloseDisplay(display);
    return 0;
}
