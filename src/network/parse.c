#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

#define EXPECTED_MAGIC_NUMBER 0xA1B2C3D4

//infile headers
void read_and_print_pcap_file_header(FILE* file);


uint32_t le(uint32_t num) {
  uint32_t new = 0;
  while (num) {
    new |= (num & 0xFF);
    // only shift left if there's another shift to do
    num = num >> 1;
    if (num) {
      new << 1;
    }
  }
  return new;
}

uint8_t swap8(uint8_t num) {
  return ((num & 0x0F) << 4) | (num >> 4);
}

long get_end_of_file(FILE* file) {
  // record current position
  long current = ftell(file);
  // go to the end of the file
  fseek(file, 0, SEEK_END);
  // record end of the file
  long end_of_file = ftell(file);
  // reset file position
  rewind(file);
  fseek(file, current, SEEK_CUR);
  return end_of_file;
}

typedef struct {
  uint32_t ts_sec;
  uint32_t ts_micro;
  uint32_t data_len; // in bytes
  uint32_t untruc_data_len;
} pcap_packet;

void read_pcap_packet(FILE* file, pcap_packet* p) {
  fread(&p->ts_sec, sizeof(uint32_t), 1, file);
  fread(&p->ts_micro, sizeof(uint32_t), 1, file);
  fread(&p->data_len, sizeof(uint32_t), 1, file);
  fread(&p->untruc_data_len, sizeof(uint32_t), 1, file);
}

typedef struct {
  char mac_dest[6];
  char mac_src[6];
  uint16_t ether_type;
} eth_packet;

void read_eth_packet(FILE* file, eth_packet* p) {
  fread(&p->mac_dest, sizeof(char), 6, file);
  fread(&p->mac_src, sizeof(char), 6, file);
  fread(&p->ether_type, sizeof(uint16_t), 1, file);
}

typedef struct {
  uint8_t version;
  uint8_t header_len;
  uint8_t dscp;
  uint8_t ecn;
  uint16_t total_len;
  uint16_t identification;
  uint8_t flags;
  uint16_t fragment_offset;
  uint8_t ttl;
  uint8_t protocol;
  uint16_t header_checksum;
  uint32_t src_ip;
  uint32_t dest_ip;
} ipv4_packet;

void read_ipv4_packet(FILE* file, ipv4_packet* p) {
  uint8_t version_and_header_len;
  fread(&version_and_header_len, sizeof(uint8_t), 1, file);
  // little endian
  // don't have to????
  //version_and_header_len = swap8(version_and_header_len);
  p->version = version_and_header_len >> 4;
  p->header_len = (version_and_header_len & 0x0F) * 4; // in bytes
  // sanity check
  if (p->header_len < 20) {
    printf("Should be greater than 20 bytes!!!");
    exit(1);
  }
  if (p->header_len > 20) {
    printf("Should be less than 60 bytes!!!");
    exit(1);
  }
  uint8_t dscp_and_ecn;
  fread(&dscp_and_ecn, sizeof(uint8_t), 1, file);
  p->dscp = dscp_and_ecn & 0x03
  p->ecn = dscp_and_ecn >> 2;
  fread(&p->total_len, sizeof(uint16_t), 1, file);
  fread(&p->identification, sizeof(uint16_t), 1, file);
  uint16_t flags_and_fragment_offset;
  fread(&flags_and_fragment_offest, sizeof(uint16_t), 1, file);
  p->flags = flags_and_fragment_offset >> 13;
  p-> 
}

typedef struct {
  uint32_t seq;
  uint32_t data_len;
  char* buffer;
} tcp_data;


int main() {
  // plan to write the code to parse two packets
  // pcap packet -> eth packet -> ipv4 packet -> tcp packet -> HTTP protocol

  FILE* file = fopen("./net.cap", "rb");
  read_and_print_pcap_file_header(file);
  pcap_packet pcap_pkt;
  eth_packet eth_pkt;
  ipv4_packet ipv4_pkt;
  
  
  read_pcap_packet(file, &pcap_pkt);
  read_eth_packet(file, &eth_pkt);
  printf("eth_type: %d\n", eth_pkt.ether_type);
  read_ipv4_packet(file, &ipv4_pkt);
  printf("ipv4 version: %d\n", ipv4_pkt.version);
  printf("ipv4 len: %d\n", ipv4_pkt.header_len);


  // Read first packet

  // read second packet

  fclose(file);
}
////  while (1) {
////    fseek(file, pcap_pkt.data_len, SEEK_CUR);
////    if (ftell(file) >= end_of_file) {
////      break;
////    }
////    read_packet(file, &pcap_pkt);
////    printf("packlen: %u\n", pcap_pkt.data_len);
////    printf("packlenuntruc: %u\n", pcap_pkt.untruc_data_len);
////    count++; 
////  }


//  long packet_start;
//  long packet_end;
//  long end_of_file = get_end_of_file(file);
//  pcap_packet pcap_pkt;
//  eth_packet eth_pkt;
//  long count = 0;
//  //read 
//  int read = 12 + 2;
//  while (1) {
//    packet_start = ftell(file);
//    printf("Packet start: %ld\n", packet_start);
//    if (packet_start >= end_of_file) {
//      break;
//    }
//    read_pcap_packet(file, &pcap_pkt); 
//    //read_eth_packet(file, &eth_pkt); 
//    //printf("ethpacket type: %d\n", eth_pkt.ether_type);
//    printf("Datalen: %ld\n", pcap_pkt.data_len);
//    packet_end = packet_start + pcap_pkt.data_len;
//    printf("Pkt End: %ld\n", packet_end);
//    count++;
//    // jump to next packet
//    fseek(file, pcap_pkt.data_len, SEEK_CUR);
//  }
//  printf("Packet count: %ld\n", count);



void read_and_print_pcap_file_header(FILE* file) {
  uint32_t magic_num;
  uint16_t major_version;
  uint16_t minor_version;
  uint32_t offset;
  uint32_t ts_accuracy;
  uint32_t snapshot_len;
  uint32_t link_layer_header_type;
  fread(&magic_num, sizeof(uint32_t), 1, file);
  fread(&major_version, sizeof(uint16_t), 1, file);
  fread(&minor_version, sizeof(uint16_t), 1, file);
  fread(&offset, sizeof(uint32_t), 1, file);
  fread(&ts_accuracy, sizeof(uint32_t), 1, file);
  fread(&snapshot_len, sizeof(uint32_t), 1, file);
  fread(&link_layer_header_type, sizeof(uint32_t), 1, file);
  // end per file headers
  
  if (magic_num == EXPECTED_MAGIC_NUMBER) {
    printf("Swap not necessary\n");
  } else {
    printf("Swap necessary but not supported. Exiting!\n");
    exit(1);
  }
  printf("Major version: %u\n", major_version);
  printf("Minor version: %u\n", minor_version);
  if (offset == 0) {
    printf("Offset is zero as expected\n");
  } else {
    printf("Offset is non-zero. Unexpected. Exiting!\n");
    exit(1);
  }
  if (ts_accuracy == 0) {
    printf("Timestamp accuracy is zero as expected\n");
  } else {
    printf("Timestampt accuracy is non-zero. Unexpected. Exiting!\n");
    exit(1);
  }
  printf("Snapshot Length: %u\n", snapshot_len);
  printf("Link layer header type: %u (1 is Ethernet)\n", link_layer_header_type);
}

////  // count packets
////  int count = 2;
////  int end_of_file = get_end_of_file(file);
////  printf("end: %ul\n", end_of_file);
////
////  printf("Packet count: %d\n", count);
////1049
//
////  ipv4_packet ipv4_pkt;
////  read_ipv4_packet(file, &ipv4_pkt);
////  printf("version: %d\n", ipv4_pkt.version);
////  printf("len: %d\n", ipv4_pkt.header_len);
//
//  // second packet
//  fseek(file, pcap_pkt.data_len, SEEK_CUR);
//  read_packet(file, &pcap_pkt);
//  printf("TS sec: %u\n", pcap_pkt.ts_sec);
//  printf("TS micro: %u\n", pcap_pkt.ts_micro);
//  printf("TS dat_len: %u\n", pcap_pkt.data_len);
//  printf("TS truc_data_len: %u\n", pcap_pkt.untruc_data_len);
//
//
//
//
//
////  uint32_t ts_sec;
////  uint32_t ts_micro;
////  uint32_t data_len; // in bytes
////  uint32_t untruc_data_len;
////  fread(&ts_sec, sizeof(uint32_t), 1, file);
////  fread(&ts_micro, sizeof(uint32_t), 1, file);
////  fread(&data_len, sizeof(uint32_t), 1, file);
////  fread(&untruc_data_len, sizeof(uint32_t), 1, file);
////  printf("TS sec: %u\n", ts_sec);
////  printf("TS micro: %u\n", ts_micro);
////  printf("TS data_len: %u\n", data_len);
////  printf("TS untruct_data_len: %u\n", untruc_data_len);
//  
//
//
