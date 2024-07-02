.data
newline: .asciiz "\n"
true_str: .asciiz "true"
false_str: .asciiz "false"
str_0: .asciiz "Holaa"
.text
.globl main
main:
move $fp, $sp
sw $ra, 0($sp)
addiu $sp, $sp, -4
li $t0, 5
li $t1, 2
mul $t0, $t0, $t1
sw $t0, -4($sp)
lw $t1, -4($sp)
move $a0, $t1
li $v0, 1
syscall
li $v0, 4
la $a0, newline
syscall
la $t2, str_0
move $a0, $t2
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
li $t2, 5
lw $t3, -4($sp)
seq $t4, $t2, $t3
beq $t4, $zero, label_1
li $t3, 1
beq $t3, $zero, label_3
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
li $t4, 0
beq $t4, $zero, label_5
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
lw $ra, 4($sp)
addiu $sp, $sp, 4
li $v0, 10
syscall
