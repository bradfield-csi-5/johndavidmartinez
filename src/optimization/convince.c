#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <strings.h>

int main() {
    uint64_t a = ffs(8);
    uint64_t b = ffs(2);
    printf("a: %llu\n", a);
    printf("b: %llu\n", b);
    printf("%llu\n", 1LLU << (a - b));
}
