.data
newline: .asciiz "\n"
true_str: .asciiz "true"
false_str: .asciiz "false"
.text
.globl main
main:
li $t0, 5
li $t1, 5
seq $t2, $t0, $t1
beqz $t2, label_1
li $t3, 2
li $t4, 5
mul $t5, $t3, $t4
li $t6, 45
sub $t7, $t5, $t6
move $a0, $t7
li $v0, 1
syscall
li $v0, 4
la $a0, newline
syscall
j label_2
label_1:
li $t8, 0
beqz $t8, label_3
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
label_2:
move $t9, $t7
li $v0, 10
syscall
