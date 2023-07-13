
# 源文件目录
SRC_DIR := .
# 获取源文件列表
SRCS := $(wildcard $(SRC_DIR)/*.c)
# 包含的头文件
INCLUDES := -I.

OBJS=$(SRCS:.c=.o)

TARGET=glib_main

CFLAGS = -g -Wall -Wpedantic -Wno-padded -O $(shell pkg-config --cflags --libs glib-2.0 gobject-2.0)

$(TARGET): $(OBJS)
	echo $(TARGET)  $(OBJS)
	gcc $(CFLAGS) $^ -o $@
%.o: %.c
	gcc $(CFLAGS) $(INCLUDES) -c $^ -o $@

clean:
	-rm -rf $(TARGET) $(OBJS)