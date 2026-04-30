package test

import (
	"chain_simulation/modules/fast_selir"
	"fmt"
	"testing"
)

func TestFastSelir(t *testing.T) {
	//result := fast_selir.CalculateFastSelirBFBits(1, 0.00001)
	//fmt.Println(result)
	result := fast_selir.CalculateFastSelirBFBits(2, 0.00001)
	fmt.Println(result)
}
