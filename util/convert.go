package util

import (
	"fmt"
	"strconv"
)

// Base16ToBase10  convert hexadecimal to ASCII
func Base16ToBase10(base16 string) (asc int64, err error) {
	if asc, err = strconv.ParseInt(base16, 16, 0); err != nil {
		return 0, err
	}
	return asc, nil
}

// Base16ToChar convert hexadecimal to char
func Base16ToChar(base16 string) (c string, err error) {
	var asc int64
	if asc, err = strconv.ParseInt(base16, 16, 0); err != nil {
		return "", err
	}
	c = fmt.Sprintf("%c", asc)
	return c, nil
}

func Int64ToInt(i64 int64) int {
	strInt64 := strconv.FormatInt(i64, 10)
	i16, _ := strconv.Atoi(strInt64)
	return i16
}
