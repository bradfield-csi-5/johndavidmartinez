#include <stdio.h>
#include <stdint.h>
#include <unistd.h>

int main() {
  FILE* file = fopen("./num", "wb");
  uint32_t number = 34;
  fwrite(&number, sizeof(uint32_t), 1, file);
  fclose(file);

  FILE* reopen = fopen("./num", "rb");
  uint32_t readnumber = 0;
  fread(&readnumber, sizeof(uint32_t), 1, reopen);
  printf("number: %d\n", readnumber);
  fclose(reopen); 
}

