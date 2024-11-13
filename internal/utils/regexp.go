package utils

import "regexp"

var IsDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
