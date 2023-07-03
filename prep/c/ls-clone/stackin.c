#include <stdio.h>

char* first() {
    int i;
    char first[10];
    for (i = 0; i < 10; i++) {
        first[i] = 'o';
    }
    first[9] = '\0';
    return first;
}

int main() {
    char* ok = first();
    printf("%s\n", ok);
}
