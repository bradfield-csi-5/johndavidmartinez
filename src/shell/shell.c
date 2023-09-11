#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define INPUT_SIZE 200
#define RETURN 10

void pln(char* arr, int size) {
  for (int i = 0; i < size; i++) {
    putchar(arr[i]);
  }
  putchar('\n');
}

void process(char* arr, int size) {
  // extract first 'word'
  // case against builtin
  //   if builtin execute builtin
  // need to practice string separation with tokens first

  pln(arr, size);
}

void readeval() {
  char input[INPUT_SIZE];
  int c;
  int idx = 0;

  while (1) {
    printf("> ");
    while ((c = getchar()) != RETURN) {
      if (c == EOF) {
        printf("\nShell Exit...\n");
        exit(0);
      }
      input[idx] = c;
      idx++;
    }
    process(input, idx);
    idx = 0;
  }
}

int main(int argc, char* argv[]) {
  if (argc > 1) {
    if (strcmp(argv[1], "-c") == 0) {
      int len = strlen(argv[2]);
      process(argv[2], len);
    } else {
      printf("Usage\n");
      exit(1);
    }
  } else {
    readeval();
  }
}
