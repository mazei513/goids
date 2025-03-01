#include "textflag.h"

// func Sqrt32(x float32) float32
TEXT ·Sqrt32(SB),NOSPLIT,$0
	MOVSS x+0(FP), X0
	SQRTSS X0, X0
	MOVSS X0, ret+8(FP)
	RET

// func Mag32(a, b Boid32) (dx, dy, mag float32)
TEXT ·Mag32(SB),NOSPLIT,$0
	MOVUPS a+0(FP), X0
	MOVUPS b+16(FP), X1
	SUBPS X1, X0
	MOVUPS X0, ret+32(FP)
	MULPS X0, X0
	HADDPS X0, X0
	SQRTSS X0, X0
	MOVSS X0, ret+40(FP)
	RET

// func Mags32(a Boid32, b []Boid32, d []Vec32, mag []float32)
TEXT ·Mags32(SB),NOSPLIT,$0
	MOVLPS a+0(FP), X0

	MOVQ b_base+16(FP), AX
	MOVQ b_len+24(FP), BX

	MOVQ dx_base+40(FP), CX
	MOVQ mag_base+64(FP), DX

loop:
	MOVLPS (AX), X1
	MOVHPS 16(AX), X1
	SUBPS X1, X0
	MOVUPS X0, (CX)
	MULPS X0, X0
	HADDPS X0, X0
	SQRTSS X0, X0
	MOVSS X0, (DX)
	UNPCKHPS X0, X0
	MOVSS X0, 4(DX)

	ADDQ $32, AX
	ADDQ $16, CX
	ADDQ $8, DX
	SUBQ $2, BX
	JNZ loop
	
	RET
