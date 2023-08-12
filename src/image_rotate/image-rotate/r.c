#include <stdio.h>
#include <sys/types.h>

typedef struct {
    char head[3];
    uint32_t size;
    uint32_t reserved;
    uint16_t offset;
    uint32_t header_size;
    uint32_t width;
    uint32_t height;
} bmp_file;

// convert to little endian
uint32_t le(uint32_t i) {
    return (i << 16) | (i >> 16);
}

int main() {
    bmp_file b;
    char *name = NULL;
    char *extra = NULL;

    FILE* file = fopen("./teapot.bmp", "rb");
    // read header
    fread(b.head, sizeof(char), 2, file);
    b.head[2] = '\0';

    fread(&b.size, sizeof(b.size), 1, file);
    fread(&b.reserved, sizeof(b.reserved), 1, file);
    fread(&b.offset, sizeof(b.offset), 1, file);
    fread(&b.header_size, sizeof(b.header_size), 1, file);
    fread(&b.width, sizeof(b.width), 1, file);
    fread(&b.height, sizeof(b.height), 1, file);

    printf("Header: %s\n", b.head);
    printf("Size in Bytes: %d\n", b.size);
    printf("Offset: %d\n", b.offset);
    printf("Header size %d\n", b.header_size);
    printf("Dimensions: %dx%d\n", b.width, b.height);

    // jump to

    fclose(file);
}












