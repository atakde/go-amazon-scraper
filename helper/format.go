package helper

import (
	"regexp"
)

func FormatPrice(price string) string {
	re := regexp.MustCompile(`(?m)([0-9.,]+)`)
	return re.FindString(price)
}
