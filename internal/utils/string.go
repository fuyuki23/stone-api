package utils

import "bytes"

func AppendString(strs ...string) string {
	var b bytes.Buffer
	for _, str := range strs {
		b.WriteString(str)
	}
	return b.String()
}
