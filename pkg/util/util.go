package util

import (
	"path"
	"strings"
)

func GetServiceNameFromFullMethod(fm string) string {
	return strings.Trim(path.Dir(fm), "/")
}
