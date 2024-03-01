#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <sys/types.h>
#include <sys/stat.h>

#define FIFO_PATH "/run/systemd/inhibit/15.ref"
#define SOCK_PATH "/home/uos/dde-go/src/github.com/linuxdeepin/golang/deepin-linux/fifo/mysocket"

int main() {
    // 创建一个 FIFO
    mkfifo(FIFO_PATH, 0666);

    // 打开 FIFO 以写入数据
    int fd = open(FIFO_PATH, O_RDONLY|O_CLOEXEC|O_NONBLOCK);
    if (fd == -1) {
        perror("open");
        exit(EXIT_FAILURE);
    }
    printf("open done\n");
    // 创建 UNIX 域套接字
    int sockfd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (sockfd == -1) {
        perror("socket");
        exit(EXIT_FAILURE);
    }
    printf("socket done\n");
    // 设置套接字地址
    struct sockaddr_un addr;
    memset(&addr, 0, sizeof(struct sockaddr_un));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, SOCK_PATH, sizeof(addr.sun_path) - 1);

    // 连接到服务器端
    if (connect(sockfd, (struct sockaddr *)&addr, sizeof(struct sockaddr_un)) == -1) {
        perror("connect");
        exit(EXIT_FAILURE);
    }

    // 准备要发送的数据
    int buffer = fd;

    // 使用 sendmsg() 发送数据
    struct msghdr msg = {0};
    struct iovec iov = { .iov_base = &buffer, .iov_len = sizeof(buffer) };
    msg.msg_iov = &iov;
    msg.msg_iovlen = 1;

    if (sendmsg(sockfd, &msg, 0) < 0) {
        perror("sendmsg");
        exit(EXIT_FAILURE);
    }
    printf("socket done\n");

    // 关闭 FIFO 和套接字
    close(fd);
    close(sockfd);

    return 0;
}
