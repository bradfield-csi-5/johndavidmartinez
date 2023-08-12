#include <stdio.h>
#include <sys/types.h>

int main() {
    uint16_t i = 1;
    while (i) {
	printf("%u\n", i);
	i = i << 1 | i;
    }
}

