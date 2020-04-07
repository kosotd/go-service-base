package utils

import (
	"strings"
)

func Equals(str1, str2 string) bool {
	str1 = strings.ToLower(strings.Trim(str1, " "))
	str2 = strings.ToLower(strings.Trim(str2, " "))
	return strings.Compare(str1, str2) == 0
}

func NotEquals(str1, str2 string) bool {
	return !Equals(str1, str2)
}

func Contains(str1, str2 string) bool {
	str1 = strings.ToLower(strings.Trim(str1, " "))
	str2 = strings.ToLower(strings.Trim(str2, " "))
	return strings.Contains(str1, str2)
}

func NotContains(str1, str2 string) bool {
	return !Contains(str1, str2)
}

func Empty(str string) bool {
	return Equals("", str)
}

func NotEmpty(str string) bool {
	return !Empty(str)
}
