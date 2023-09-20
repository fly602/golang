OBJS=main.o \
	agent.o \
	sd_bus_method.o


TARGET=lastore-agent

CFLAGS = -g -Wall -Wpedantic -Wno-padded -O `pkg-config --cflags --libs libsystemd`

$(TARGET): $(OBJS)
	gcc $(CFLAGS) $(OBJS) -o $@
%.o: %.c
	gcc $(CFLAGS) -c $^ -o $@

clean:
	-rm -rf $(TARGET) $(OBJS)