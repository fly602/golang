package main

import (
	"fmt"
)

func calcMaxScaleFactor(width, height uint16) float64 {
	scaleFactors := []float64{1.0, 1.25, 1.5, 1.75, 2.0, 2.25, 2.5, 2.75, 3.0}

	if width < height {
		width, height = height, width
	}
	maxWScale := float64(width) / 1024.0
	maxHScale := float64(height) / 768.0

	maxValue := 0.0
	if maxWScale < maxHScale {
		maxValue = maxWScale
	} else {
		maxValue = maxHScale
	}

	if maxValue > 3.0 {
		maxValue = 3.0
	}

	maxScale := 1.0
	for idx := 0; (float64(idx)*0.25 + 1.0) <= maxValue; idx++ {
		maxScale = scaleFactors[idx]
	}

	return maxScale
}

func main() {

	fmt.Println(calcMaxScaleFactor(1920, 1080))
}
