.data
newline: .asciiz "\n"
true_str: .asciiz "true"
false_str: .asciiz "false"
str_0: .asciiz "Holaa"
.text
.globl main
main:
la $t0, str_0
move $a0, $t0
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
li $t0, 5
li $t1, 5
seq $t1, $t0, $t1
beq $t1, $zero, label_1
li $t2, 1
beq $t2, $zero, label_3
la $a0, true_str
j label_4
label_3:
la $a0, false_str
label_4:
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
j label_2
label_1:
li $t3, 0
beq $t3, $zero, label_5
la $a0, true_str
j label_6
label_5:
la $a0, false_str
label_6:
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
label_2:
li $v0, 10
syscall
