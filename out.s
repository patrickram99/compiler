.data
newline: .asciiz "\n"
true_str: .asciiz "true"
false_str: .asciiz "false"
str_0: .asciiz "Holaaa"
.text
.globl main
main:
li $t0, 2
li $t1, 1
li $t2, 1
li.s $f0, 0.500000
mtc1 $t2, $f1
cvt.s.w $f1, $f1
mul.s $f2, $f1, $f0
mtc1 $t1, $f3
cvt.s.w $f3, $f3
add.s $f4, $f3, $f2
mtc1 $t0, $f5
cvt.s.w $f5, $f5
c.eq.s $f5, $f4
li $t3, 1
bc1t float_true_0
li $t3, 0
float_true_0:
beq $t3, $zero, label_2
li $t4, 1
beq $t4, $zero, label_4
la $a0, true_str
j label_5
label_4:
la $a0, false_str
label_5:
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
j label_3
label_2:
li $t5, 0
beq $t5, $zero, label_6
la $a0, true_str
j label_7
label_6:
la $a0, false_str
label_7:
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
label_3:
la $t0, str_0
move $a0, $t0
li $v0, 4
syscall
li $v0, 4
la $a0, newline
syscall
li $v0, 10
syscall
