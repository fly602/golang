OBJS=main.o

TARGET=main

CFLAGS = -g -Wall -Wpedantic -Wno-padded  -lpthread -O $(shell pkg-config --cflags --libs glib-2.0)

$(TARGET): $(OBJS)
	gcc $(CFLAGS) $(OBJS) -o $@
%.o: %.c
	gcc $(CFLAGS) -c $^ -o $@

clean:
	-rm -rf $(TARGET) $(OBJS)