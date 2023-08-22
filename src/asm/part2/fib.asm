section .text
global fib

; convention is eax has return
func:
    pop r12
    ; if r12 < 2 return r12
    sub r12, 1
    push r12
    sub r12, 1
    push r12
    
fib:
    ; push rdi to the stack?
    push rdi
    mov r9, 10
    call func
    ret


; push n to the stack?
; fib(n-1) + fib(n - 2)

; factorial 2! = 2 * 1! = 2 * 1
