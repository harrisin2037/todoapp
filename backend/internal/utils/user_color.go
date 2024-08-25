package utils

import (
	"fmt"
	"hash/fnv"
	"math"
)

func GenerateColor(id uint) string {

	hash := fnv.New32a()

	hash.Write([]byte(fmt.Sprintf("%d", id)))
	hue := math.Mod(float64(hash.Sum32()), 360)

	return fmt.Sprintf("hsl(%.0f, 70%%, 80%%)", hue)
}
