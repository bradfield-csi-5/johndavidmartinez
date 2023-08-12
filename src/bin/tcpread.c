#include <stdio.h>
#include <sys/types.h>

// TCP flags
enum {
    FIN = 0x01,
    SYN = 0x02,
    RST = 0x04,
    PSH = 0x08,
    ACK = 0x10,
    URG = 0x20,
    ECE = 0x40,
    CWR = 0x80
};

typedef struct {
    uint16_t src_port;
    uint16_t des_port;
    uint32_t seq_number;
    uint32_t ack_number;
    char data_offset;
    char reserved;
    uint8_t flags;
    uint16_t win_size;
    uint16_t checksum;
    uint16_t urgent_ptr;
} tcp_header;

uint16_t swap16(uint16_t i) {
    return (i << 8) | (i >> 8);
}

uint32_t swap32(uint32_t i) {
    return (
	    ((i >> 24) & 0xFF) |
	    ((i << 8) & 0xFF0000) |
	    ((i >> 8) & 0xFF00) |
	    ((i << 24) & 0xFF000000)
    );
}

int main() {
    tcp_header tcph;
    FILE* file = fopen("./tcpheader", "rb");
    fread(&tcph.src_port, sizeof(tcph.src_port), 1, file);
    fread(&tcph.des_port, sizeof(tcph.des_port), 1, file);
    fread(&tcph.seq_number, sizeof(tcph.seq_number), 1, file);
    fread(&tcph.ack_number, sizeof(tcph.ack_number), 1, file);
    fread(&tcph.data_offset, sizeof(tcph.data_offset), 1, file);
    fread(&tcph.reserved, sizeof(tcph.reserved), 1, file);
    fread(&tcph.flags, sizeof(tcph.flags), 1, file);
    // to LE
    tcph.src_port = swap16(tcph.src_port);
    tcph.des_port = swap16(tcph.des_port);
    tcph.seq_number = swap32(tcph.seq_number);
    tcph.ack_number = swap32(tcph.ack_number);

    printf("Source Port: %u\n", tcph.src_port);
    printf("Destination Port: %u\n", tcph.des_port);
    printf("Sequence Number: %u\n", tcph.seq_number);
    printf("Acknowledgement Number: %u\n", tcph.ack_number);
    printf("Data Offset: %c\n", tcph.data_offset);
    printf("Reserved: %c\n", tcph.reserved);
    if ((tcph.flags & FIN) == FIN) {
	printf("FIN set\n");
    }
    if ((tcph.flags & SYN) == SYN) {
	printf("SYN set\n");
    }
    if ((tcph.flags & RST) == RST) {
	printf("RST set\n");
    }
    if ((tcph.flags & PSH) == PSH) {
	printf("PSH set\n");
    }
    if ((tcph.flags & ACK) == ACK) {
	printf("ACK set\n");
    }
    if ((tcph.flags & URG) == URG) {
	printf("URG set\n");
    }
    if ((tcph.flags & ECE) == ECE) {
	printf("ECE set\n");
    }
    if ((tcph.flags & CWR) == CWR) {
	printf("CWR set\n");
    }
    fclose(file);
}




