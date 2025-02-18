package util

import (
	"math"
	"strings"
)

func Round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func Precision(target float64, source string) float64 {
	points := strings.SplitAfter(source, ".")

	if len(points) == 2 {
		return Round(target, len(points[1]))
	}

	return target
}
