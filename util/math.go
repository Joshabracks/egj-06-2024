package util

func Clamp(f, min, max float64) float64 {
	if f > max {
		f = max
	}
	if f < min {
		f = min
	}
	return f
}
