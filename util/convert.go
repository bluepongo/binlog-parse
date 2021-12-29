package util

import (
	"fmt"
	"strconv"
)

// Base16ToBase10  convert hexadecimal to ASCII
func Base16ToBase10(base16 string) (base10 int64, err error) {
	if base10, err = strconv.ParseInt(base16, 16, 0); err != nil {
		return 0, err
	}
	return base10, nil
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

// Int64ToInt convert int64 to int
func Int64ToInt(i64 int64) int {
	strInt64 := strconv.FormatInt(i64, 10)
	i16, _ := strconv.Atoi(strInt64)
	return i16
}

// Base16ToBase2 convert base16 to base2
func Base16ToBase2(base16 string) (strBase2 string, err error) {
	var base2 int64
	if base2, err = strconv.ParseInt(base16, 16, 10); err != nil {
		return "", err
	}
	strBase2 = strconv.FormatInt(base2, 2)
	return strBase2, nil
}
