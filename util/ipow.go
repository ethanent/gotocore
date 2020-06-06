package util

// Ipow returns a^b
func Ipow(a, b int) int {
	v := a

	for i := 1; i < b; i++ {
		v *= a
	}

	return v
}
