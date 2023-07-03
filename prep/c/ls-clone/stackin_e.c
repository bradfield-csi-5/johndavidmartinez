#include <stdio.h>

void first(char* buf) {
    int i;
    for(i = 0; i < 4; i++) {
        buf[i] = 'f';
    }
    buf[3] = '\0';
}

void second(char* buf) {
    int i;
    for(i = 0; i < 4; i++) {
        buf[i] = 's';
    }
    buf[3] = '\0';
}

int main() {
    char f[4], s[4];
    first(f);
    second(s);

    printf("%s\n", f);
    printf("%s\n", s);
}
