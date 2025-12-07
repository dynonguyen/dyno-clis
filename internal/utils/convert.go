package utils

import "strings"

var layoutMapping = map[string]string{
	"Y": "2006",
	"M": "01",
	"D": "02",
	"h": "15",
	"m": "04",
	"s": "05",
	"f": "000",
}

func ConvertDateLayout(format string) string {
	for old, new := range layoutMapping {
		format = strings.ReplaceAll(format, old, new)
	}
	return format
}
