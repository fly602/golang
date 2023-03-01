#include <stdio.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>

#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>
#include <glib.h>

#include <sys/eventfd.h>
#include <unistd.h>

static gchar control_path[128] = "/run/user/1000/gio/control";

static gboolean
control_input (GIOChannel *channel, GIOCondition cond, gpointer user_data);

static void
control_data_free (gpointer data)
{
	char *cdata = data;
	free(cdata);
}

static gboolean
control_accept (GIOChannel *channel, GIOCondition cond, gpointer callback_data)
{
	struct sockaddr_un addr;
	socklen_t addrlen;
	char *cdata;
	GIOChannel *new_channel;
	int fd, new_fd;
	int val;

	fd = g_io_channel_unix_get_fd (channel);

	addrlen = sizeof (addr);
	new_fd = accept (fd, (struct sockaddr *) &addr, &addrlen);
	if (new_fd < 0) {
		g_warning ("couldn't accept new control request: %s", g_strerror (errno));
		return TRUE;
	}

	val = fcntl (new_fd, F_GETFL, 0);
	if (val < 0) {
		g_warning ("can't get control request fd flags: %s", g_strerror (errno));
		close (new_fd);
		return TRUE;
	}

	if (fcntl (new_fd, F_SETFL, val | O_NONBLOCK) < 0) {
		g_warning ("can't set control request to non-blocking io: %s", g_strerror (errno));
		close (new_fd);
		return TRUE;
	}

	cdata = (char *)malloc(256);
	new_channel = g_io_channel_unix_new (new_fd);
	g_io_channel_set_close_on_unref (new_channel, TRUE);
	g_io_add_watch_full (new_channel, G_PRIORITY_DEFAULT, G_IO_IN | G_IO_HUP,
	                     control_input, cdata, control_data_free);
	g_io_channel_unref (new_channel);

	return TRUE;
}

gboolean
gkd_control_listen (void)
{
	struct sockaddr_un addr;
	GIOChannel *channel;
	int sock;

	unlink (control_path);

	sock = socket (AF_UNIX, SOCK_STREAM, 0);
	if (sock < 0) {
		g_warning ("couldn't open socket: %s", g_strerror (errno));
		return FALSE;
	}

	memset (&addr, 0, sizeof (addr));
	addr.sun_family = AF_UNIX;
	g_strlcpy (addr.sun_path, control_path, sizeof (addr.sun_path));
	if (bind (sock, (struct sockaddr*) &addr, sizeof (addr)) < 0) {
		g_warning ("couldn't bind to control socket: %s: %s", control_path, g_strerror (errno));
		close (sock);
		return FALSE;
	}

	if (listen (sock, 128) < 0) {
		g_warning ("couldn't listen on control socket: %s: %s", control_path, g_strerror (errno));
		close (sock);
		return FALSE;
	}

	channel = g_io_channel_unix_new (sock);
	g_io_add_watch (channel, G_IO_IN | G_IO_HUP, control_accept, NULL);
	g_io_channel_set_close_on_unref (channel, TRUE);
	// g_io_channel_unref(channel);

	return TRUE;
}

static gboolean
control_input (GIOChannel *channel, GIOCondition cond, gpointer user_data)
{
	char *cdata = user_data;
	guint32 packet_size = 0;
	gboolean finished = FALSE;
	int fd, res;
	pid_t pid;
	uid_t uid;
	char buffer[128];

	fd = g_io_channel_unix_get_fd (channel);

	if (cond & G_IO_IN) {
		if (errno != EAGAIN && errno != EINTR)
					finished = TRUE;
		int n = read(fd,buffer,128);
		if (n < 0){
			return FALSE;
		} else {
			g_warning("read buffer[%d]: %s",n,buffer);
		}
	}

	if (finished)
		cond |= G_IO_HUP;

	return (cond & G_IO_HUP) == 0;
}

void eventfd_test(){
	int efd = eventfd(0, EFD_NONBLOCK | EFD_CLOEXEC);
    eventfd_write(efd, 2);
    eventfd_t count;
    eventfd_read(efd, &count);
	// getchar();
	// close(efd);
}

int main(){
	g_warning("running...");
	GMainLoop *main_loop = g_main_loop_new(NULL, FALSE);

	for (int i = 0;i < 1019;i++){
		eventfd_test();
	}
	gkd_control_listen();

	// getchar();

	g_main_loop_run(main_loop);
}