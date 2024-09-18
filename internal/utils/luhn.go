package utils

import "unicode"

func LuhnCheck(number string) bool {
	var sum int
	double := false

	// Идем с конца строки
	for i := len(number) - 1; i >= 0; i-- {
		digit := number[i]

		if !unicode.IsDigit(rune(digit)) {
			return false
		}

		// quick ascii->int conversion: "5" - "0" == 53 - 48  == 5
		n := int(digit - '0')

		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		double = !double
	}

	return sum%10 == 0
}
