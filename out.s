.data
newline: .asciiz "\n"
.text
.globl main
main:
li $t0, 25
move $s0, $t0
li $t0, 2
move $t1, $t0
li $t0, 4
move $t2, $t0
add $t0, $t1, $t2
move $s1, $t0
mul $t0, $s0, $s1
move $a0, $t0
li $v0, 1
syscall
la $a0, newline
li $v0, 4
syscall
li $v0, 10
syscall
