CC=gcc
CFLAGS=-I.
DEPS = chunk.h memory.h common.h debug.h value.h
OBJ = main.o chunk.o memory.o debug.o value.o

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

vm: $(OBJ)
	$(CC) -o $@ $^ $(CFLAGS)
