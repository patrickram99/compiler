.data
newline: .asciiz "\n"
.text
.globl main
main:
li $t0, 2
li.s $f0, 2.500000
li $t1, 2
mtc1 $t1, $f1
cvt.s.w $f1, $f1
mul.s $f2, $f0, $f1
mtc1 $t0, $f3
cvt.s.w $f3, $f3
sub.s $f4, $f3, $f2
li $t2, 25
mtc1 $t2, $f5
cvt.s.w $f5, $f5
sub.s $f6, $f4, $f5
li $t3, 45
mtc1 $t3, $f7
cvt.s.w $f7, $f7
add.s $f8, $f6, $f7
li $t4, 2
li $t5, 4
mul $t6, $t4, $t5
li $t7, 5
sub $t8, $t6, $t7
mtc1 $t8, $f9
cvt.s.w $f9, $f9
sub.s $f10, $f8, $f9
mov.s $f12, $f10
li $v0, 2
syscall
li $v0, 4
la $a0, newline
syscall
li $v0, 10
syscall
