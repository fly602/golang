OBJS=main.o


TARGET=main

CFLAGS = -g -Wall -Wpedantic -Wno-padded -O `pkg-config --cflags --libs libsystemd glib-2.0 gio-2.0`

$(TARGET): $(OBJS)
	gcc $(CFLAGS) $(OBJS) -o $@
%.o: %.c
	gcc $(CFLAGS) -c $^ -o $@

clean:
	-rm -rf $(TARGET) $(OBJS)