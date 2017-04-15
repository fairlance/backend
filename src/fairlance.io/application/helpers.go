package application

import (
	"math"
)

func round(val float64, roundOn float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	return round / pow
}

func removeFromUINTSlice(slice []uint, value uint) []uint {
	for index, val := range slice {
		if val == value {
			slice[index] = slice[len(slice)-1]
			slice = slice[:len(slice)-1]
		}
	}
	return slice
}
