package mat4

import (
	"fmt"
	"math"
)

func Create() []float64 {
	var out = make([]float64, 16)
	for i := range out {
		if i%5 == 0 {
			fmt.Println(i)
			out[i] = 1
		} else {
			out[i] = 0
		}
	}
	return out
}

func Perspective(out []float64, fovy float64, aspect float64, near float64, far float64) []float64 {
	f := 1.0 / math.Tan(fovy/2)
	nf := 1 / (near - far)
	out[0] = f / aspect
	out[1] = 0
	out[2] = 0
	out[3] = 0
	out[4] = 0
	out[5] = f
	out[6] = 0
	out[7] = 0
	out[8] = 0
	out[9] = 0
	out[10] = (far + near) * nf
	out[11] = -1
	out[12] = 0
	out[13] = 0
	out[14] = 2 * far * near * nf
	out[15] = 0
	return out
}
