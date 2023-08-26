section .text
global fib

fib:
  cmp rdi, 0
  je .base

  push rdi
  sub rdi, 1
  call fib
  pop rdi
  imul rax, rdi
  ret
.base
  mov rax, 1
  ret
