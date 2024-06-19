package util

func Clamp(f, min, max float32) float32 {
	if f > max {
		f = max
	}
	if f < min {
		f = min
	}
	return f
}
