.data
newline: .asciiz "\n"
true_str: .asciiz "true"
false_str: .asciiz "false"
.text
.globl main
main:
move $fp, $sp
sw $ra, 0($sp)
addi $sp, $sp, -4
.data
omg: .word 0
.text
sw $t0, omg
lw $ra, 4($sp)
addi $sp, $sp, 4
jr $ra
