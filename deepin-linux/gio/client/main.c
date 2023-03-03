#include <stdio.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>

#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>
#include <glib.h>

static gchar control_path[128] = "/run/user/1000/gio/control";

void connect_daemon(){
    int sock;
    struct sockaddr_un addr;

    sock = socket (AF_UNIX, SOCK_STREAM, 0);
	if (sock < 0) {
		g_warning("couldn't create control socket: %s", strerror (errno));
		return;
	}

    /* close on exec */
	fcntl (sock, F_SETFD, 1);

    memset (&addr, 0, sizeof (addr));
	addr.sun_family = AF_UNIX;
	g_strlcpy (addr.sun_path, control_path, sizeof (addr.sun_path));

	if (connect (sock, (struct sockaddr *)&addr, sizeof (addr)) < 0) {
		if (errno == ECONNREFUSED) {
			close (sock);
			return;
		}
		g_warning ("couldn't connect to gnome-keyring-daemon socket at: %s: %s",
		        addr.sun_path, strerror (errno));
		close (sock);
		return;
	}
    char buffer[1024]= "hello world";
    int n = write(sock,buffer,strlen(buffer));
    g_warning ("client send buffer[%d]: %s",n,buffer);
    getchar();
    close (sock);
	return;
}

void fork_pid(){
   int i = 0;
   pid_t pid;

    for (i = 0;i < 2000;i++){
        pid = fork();
        if (pid == 0){
            connect_daemon();
        }
    }
}

int main(){
    fork_pid();
    getchar();
}