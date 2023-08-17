#include "common.h"
#include "chunk.h"
#include "debug.h"

#include <stdio.h>

int main(int argc, const char* argv[]) {
    printf("Start chunk\n");
    Chunk chunk;
    initChunk(&chunk);
    int constant = addConstant(&chunk, 1.2);
    writeChunk(&chunk, OP_CONSTANT, 123);
    writeChunk(&chunk, constant, 123);
    writeChunk(&chunk, OP_RETURN, 123);
    disassembleChunk(&chunk, "test chunk");
    freeChunk(&chunk);
    printf("End chunk\n");
    return 0;
}
