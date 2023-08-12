section .text
global binary_convert
binary_convert:
  xor rax, rax
  xor r9, r9
  jmp calc_value
shift_bit:
  shl rax, 1
  inc r9
calc_value:
  movzx r10, byte[rdi + r9]
  cmp r10, 0
  je return
  cmp r10, 48 ; 48 ascii for 0
  je shift_bit
  or rax, 0x01 ; set lowest bit
  jmp shift_bit
return:
  shr rax, 1
  ret





; solution finding strlen first
;binary_convert:
;  xor r9, r9 ;r9 will be string length
;  xor rax, rax ; return value
;find_length:
;  movzx r10, byte[rdi + r9] ; read first character
;  inc r9
;  cmp r10, 0
;  jne find_length
;calc_setup:
;  dec r9 ; r9 is length + 1
;  xor r11, r11
;  mov r11, 1 ; power
;  jmp calc_value
;calc_increase_power:
;  imul r11, 2
;calc_value:
;  dec r9
;  cmp r9, -1
;  je return
;
;  movzx r10, byte[rdi + r9]
;  cmp r10, 48
;  je calc_increase_power
;  add rax, r11
;  jmp calc_increase_power
;return:
;  ret
