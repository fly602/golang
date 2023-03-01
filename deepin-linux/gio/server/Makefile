OBJS=main.o

TARGET=main

CFLAGS = -g -Wall -Wpedantic -Wno-padded -O $(shell pkg-config --cflags --libs glib-2.0 gobject-2.0)

$(TARGET): $(OBJS)
	gcc $(CFLAGS) $(OBJS) -o $@
%.o: %.c
	gcc $(CFLAGS) -c $^ -o $@

clean:
	-rm -rf $(TARGET) $(OBJS)