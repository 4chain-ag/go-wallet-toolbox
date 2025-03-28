package validate

import (
	"fmt"
	"regexp"
)

func checkStrLength[T ~string](s T, min, max int) error {
	if len(s) < min {
		return fmt.Errorf("at least %d length", min)
	}
	if len(s) > max {
		return fmt.Errorf("no more than %d length", max)
	}
	return nil
}

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$")

func checkHexString[T ~string](hex T) error {
	if len(hex)%2 != 0 {
		return fmt.Errorf("even length, not %d", len(hex))
	}

	if !hexRegex.MatchString(string(hex)) {
		return fmt.Errorf("hexadecimal string")
	}
	return nil
}

const MaxSatoshis = 2100000000000000

func checkSatoshis[T ~uint](satoshis T) error {
	if satoshis > MaxSatoshis {
		return fmt.Errorf("less than %d", MaxSatoshis)
	}
	return nil
}
