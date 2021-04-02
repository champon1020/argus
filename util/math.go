package util

// MinInt returns minimum integer of two arguments.
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
