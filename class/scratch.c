#include <stdio.h>
#include <stdlib.h>

int square(int n) {
  return n * n;
}

int main(int argc, char* argv[]) {
  int n = atoi(argv[1]);
  printf("%d\n", square(n));
}
