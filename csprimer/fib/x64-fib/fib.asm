section .text
global fib

fib:
  cmp rdi, 1
  jle .base

  push rdi
  sub rdi, 1
  call fib
  pop rdi
  push rax
  sub rdi, 2
  call fib
  pop rcx
  add rax, rcx
  ret

.base
  mov rax, rdi
  ret



;fib:
;  cmp rdi, 0
;  je .base
;
;  push rdi
;  sub rdi, 1
;  call fib
;  pop rdi
;  imul rax, rdi
;  ret
;.base
;  mov rax, 1
;  ret
