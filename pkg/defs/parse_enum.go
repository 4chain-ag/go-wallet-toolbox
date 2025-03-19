package defs

import (
	"fmt"
	"strings"
)

func parseEnumCaseInsensitive[T ~string](value string, enums ...T) (T, error) {
	for _, e := range enums {
		if strings.EqualFold(value, string(e)) {
			return e, nil
		}
	}
	return "", fmt.Errorf("invalid value: %s", value)
}
