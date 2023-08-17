#include "common.h"
#include "chunk.h"
#include "debug.h"

#include <stdio.h>

int main(int argc, const char* argv[]) {
    printf("Start chunk\n");
    Chunk chunk;
    initChunk(&chunk);
    int constant = addConstant(&chunk, 1.2);
    writeChunk(&chunk, OP_CONSTANT);
    writeChunk(&chunk, constant);
    writeChunk(&chunk, OP_RETURN);
    disassembleChunk(&chunk, "test chunk");
    freeChunk(&chunk);
    printf("End chunk\n");
    return 0;
}
