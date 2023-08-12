section .text
global pangram
pangram:
  xor r9, r9
  xor rax, rax
  jmp get_char
increment_idx:
  inc r9
get_char:
  movzx r11, byte[rdi + r9]
  cmp r11, 0
  je return
  jmp normalize_to_bit_char
set_bit:
  cmp r12, 0xFFFFFF
  je set_return_true
  jmp increment_idx
lowercase_char:
  add r11, 0x20
normalize_to_bit_char:
  cmp r11, 0x61
  jl lowercase_char
  sub r11, 0x61
  xor r12, r12
  or r12, 0x01
  shl r12, r11
  jmp set_bit
set_return_true:
  mov rax, 0x01
return:
  ret
