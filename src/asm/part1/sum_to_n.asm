; loop variant
; section .text
; global sum_to_n
; sum_to_n:
;   ; rdi has input
;   ; rax is return
;   mov rax, 0
; loop:
;   add rax, rdi
;   dec rdi
;   cmp rdi, 0
;   jg loop
;   ret

; fast variant
section .text
global sum_to_n
sum_to_n:
  ; rdi has input
  ; rax is return
  mov rax, 0
  add rax, rdi   ; 0 + n
  imul rdi, rdi  ; n * n
  add rax, rdi   ; n + (n * n)
  shr rax, 1     ; (n + (n * n)) / 2
  ret