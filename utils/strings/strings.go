package strings

import (
	"strings"
)

func Delete(str, substr string) string {
	return strings.ReplaceAll(str, substr, "")
}
