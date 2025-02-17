package util

import (
	"math"
	"strings"

	"github.com/spf13/cast"
)

func Round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func Precision(target, source float64) float64 {
	src := cast.ToString(source)

	points := strings.SplitAfter(src, ".")

	if len(points) < 2 {
		return target
	}

	return Round(target, len(points[1]))
}
