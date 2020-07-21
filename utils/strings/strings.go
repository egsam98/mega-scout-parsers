package strings

import (
	"strconv"
	"strings"
)

func Delete(str, substr string) string {
	return strings.ReplaceAll(str, substr, "")
}

func ToInt(str string, ignore bool) int {
	integer, err := strconv.Atoi(str)
	if !ignore && err != nil {
		panic(err)
	}
	return integer
}
