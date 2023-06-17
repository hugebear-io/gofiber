package fabric

import "strings"

func EmptyString(value string) bool {
	return strings.TrimSpace(value) == EMPTY_STRING
}
