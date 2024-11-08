package utils

import "bytes"

func AppendString(str1, str2 string) string {
	var b bytes.Buffer
	b.WriteString(str1)
	b.WriteString(str2)
	return b.String()
}
