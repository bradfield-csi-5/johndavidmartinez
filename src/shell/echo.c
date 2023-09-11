#include <stdio.h>

int main() {
  int c;
  while(1) {
    c = getchar();
    printf("val: %d\n", c);
    putchar(c);
  }
}
