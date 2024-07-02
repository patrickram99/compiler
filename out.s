.data
newline: .asciiz "\n"
true_str: .asciiz "true"
false_str: .asciiz "false"
.text
.globl main
main:
move $fp, $sp
sw $ra, 0($sp)
addiu $sp, $sp, -4
str_0: .asciiz "a,jsdadsklkjlads"
la $t0, str_0
move $a0, $t0
move $a1, $t0
jal concat_strings
move $t1, $v0
lw $ra, 4($sp)
addiu $sp, $sp, 4
li $v0, 10
syscall
