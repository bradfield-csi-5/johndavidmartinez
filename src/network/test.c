#include <stdio.h>
#include <stdint.h>


uint8_t swap8(uint8_t num) {
  return ((num & 0x0F) << 4) | (num >> 4);
}

int main() {
  uint8_t num = 17;
  printf("%d\n", num);
  printf("%d\n", swap8(num));
}
