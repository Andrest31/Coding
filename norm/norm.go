package norm

import (
	"math"

	anotherRand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func GenerateNormalInt(min, max int, mean, stddev float64, src anotherRand.Source) int {
	normal := distuv.Normal{
		Mu:    mean,
		Sigma: stddev,
		Src:   src,
	}
	for {
		value := normal.Rand()
		intValue := int(math.Round(value))

		if intValue >= min && intValue <= max {

			return intValue

		}
	}
}
