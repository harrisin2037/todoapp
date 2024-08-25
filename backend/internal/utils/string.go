package utils

import "fmt"

func StringToUint(s string) (uint, error) {
	var i uint
	if _, err := fmt.Sscanf(s, "%d", &i); err != nil {
		return 0, err
	}
	return i, nil
}
