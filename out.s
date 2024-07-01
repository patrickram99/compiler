.text
.globl main
main:
    li $t0, 12
    move $t1, $t0
    li $t0, 15
    move $t2, $t0
    add $t0, $t1, $t2
    li $v0, 10
    syscall
