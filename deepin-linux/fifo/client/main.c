#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/socket.h>
#include <sys/un.h>

#define SOCK_PATH "/home/uos/dde-go/src/github.com/linuxdeepin/golang/deepin-linux/fifo/mysocket"

int main() {
    // 创建 UNIX 域套接字
    int sockfd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (sockfd == -1) {
        perror("socket");
        exit(EXIT_FAILURE);
    }

    // 绑定地址
    struct sockaddr_un addr;
    addr.sun_family = AF_UNIX;
    strcpy(addr.sun_path, SOCK_PATH);
    bind(sockfd, (struct sockaddr*)&addr, sizeof(addr));

    // 监听连接
    if (listen(sockfd, 1) == -1) {
        perror("listen");
        exit(EXIT_FAILURE);
    }
    printf("listen done\n");
    // 接受连接
    int newsockfd = accept(sockfd, NULL, NULL);
    if (newsockfd == -1) {
        perror("accept");
        exit(EXIT_FAILURE);
    }
    printf("accept done\n");
    // 接收文件描述符
    // 接收数据
    struct msghdr msg = {0};
    int buffer;
    struct iovec iov = { .iov_base = &buffer, .iov_len = sizeof(buffer) };
    msg.msg_iov = &iov;
    msg.msg_iovlen = 1;
    ssize_t bytes_received = recvmsg(newsockfd, &msg, 0);
    if (bytes_received == -1) {
        perror("recvmsg");
        exit(EXIT_FAILURE);
    }

    printf("Received %zd bytes: %d\n", bytes_received, buffer);

    // 关闭连接和文件描述符
    close(newsockfd);
    close(sockfd);

    return 0;
}
