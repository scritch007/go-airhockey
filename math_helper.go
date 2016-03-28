package main

import (
	"math"
)

func deg_to_radf(deg float32) float32 {
	return float32(deg * math.Pi / 180.0)
}
