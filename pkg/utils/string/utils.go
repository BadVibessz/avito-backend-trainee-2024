package string

import (
	"strconv"
	"strings"
)

// FillIntSliceFromString returns int slice constructed from string of format '1,2,3...'
func FillIntSliceFromString(s string) ([]int, error) {
	slice := make([]int, 0, len(s))

	for _, intStr := range strings.Split(s, ",") {
		i, err := strconv.Atoi(intStr)
		if err != nil {
			return nil, err
		}

		slice = append(slice, i)
	}

	return slice, nil
}
