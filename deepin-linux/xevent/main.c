#include <X11/Xlib.h>
#include <X11/extensions/Xrandr.h>
#include <stdio.h>

int main() {
    Display *display;
    Window root;
    XEvent event;

    // 打开与X服务器的连接
    display = XOpenDisplay(NULL);
    if (display == NULL) {
        fprintf(stderr, "无法打开X服务器连接\n");
        return 1;
    }

    // 获取根窗口
    root = DefaultRootWindow(display);

    // 打开Xrandr扩展
    int opcode, event_base, error_base;
    if (XQueryExtension(display, "RANDR", &opcode, &event_base, &error_base)) {
        // 选择监听RRScreenChangeNotify事件
        XRRSelectInput(display, root, RRScreenChangeNotifyMask);

        // 进入事件循环
        while (1) {
            XNextEvent(display, &event);

            // 将事件转换为XRRNotifyEvent
            XRRNotifyEvent *xrrEvent = (XRRNotifyEvent *)&event;

            // 处理RRScreenChangeNotify事件
            if (xrrEvent->type == event_base + RRScreenChangeNotify) {
                printf("屏幕插拔事件：屏幕状态发生变化,event_base=%d\n",event_base);

                // 区分不同的事件子类型
                switch (xrrEvent->subtype) {
                    case RRNotify_CrtcChange:
                        printf("CRTC变化事件\n");
                        // 处理CRTC变化事件的代码
                        break;
                    case RRNotify_OutputChange:
                        printf("输出设备变化事件\n");
                        // 处理输出设备变化事件的代码
                        break;
                    case RRNotify_OutputProperty:
                        printf("输出设备属性变化事件\n");
                        // 处理输出设备属性变化事件的代码
                        break;
                    case RRNotify_ProviderChange:
                        printf("提供者变化事件\n");
                        // 处理提供者变化事件的代码
                        break;
                    default:
                        printf("未知的子类型:%d\n",xrrEvent->subtype);
                        // 可能有其他处理方式，视情况而定
                        break;
                }

                // 在这里可以添加处理屏幕插拔事件的代码
            }
        }
    } else {
        fprintf(stderr, "Xrandr扩展不可用\n");
    }

    // 关闭与X服务器的连接
    XCloseDisplay(display);

    return 0;
}
