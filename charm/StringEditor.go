package charm

import "strings"

func ReplaceAll(str string, list []string) string {
	for _, val := range list {
		str = strings.ReplaceAll(str, val, "")
	}
	return str
}
