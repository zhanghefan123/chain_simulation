package fast_selir

import "math"

func CalculateFastSelirBFBits(numberOfHops int, fpr float64) int {
	m := math.Ceil(-float64(numberOfHops) * math.E * math.Log(fpr))
	return int(m)
}
