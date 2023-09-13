#include <stdio.h>
#include <stdint.h>
#include <unistd.h>
#include <stdlib.h>

#define BUFFER_SIZE 1

int main() {
  // why doesn't it work on the stack
  char buffer[BUFFER_SIZE]; 

  FILE* file = fopen("./num", "rb");
  FILE* copy = fopen("./copy", "wb");
  int bytecount;
  char* start = buffer;
  while (feof(file) == 0) {
    printf("HERE\n");
    bytecount = fread(buffer, sizeof(char), BUFFER_SIZE, file);
    printf("bytecount: %d\n", bytecount);
    fwrite(start, sizeof(char), bytecount, copy);
  }
  fclose(file);
  fclose(copy);
}
