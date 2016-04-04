package main

import (
	"math"
)

func deg_to_radf(deg float32) float32 {
	return float32(deg * math.Pi / 180.0)
}

/*
static float clamp(float value, float min, float max) {
	return fmin(max, fmax(value, min));
}
*/
func clamp(value, min, max float32) float32 {
	return float32(math.Min(float64(max), math.Max(float64(value), float64(min))))
}
