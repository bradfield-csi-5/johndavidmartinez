#include <stdio.h>

char* first() {
    int i;
    char ok[4];
    for (i = 0; i < 4; i++) {
        ok[i] = 'z';
    }
    ok[3] = '\0';
    return ok;
}

void nest() {
    printf("%s\n", first());
}

int main() {
    nest();
}
