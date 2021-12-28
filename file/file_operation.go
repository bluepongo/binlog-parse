package file

import (
	"io/ioutil"
	"strings"
)

// ReadBinlog read the binlog file into the program
func ReadBinlog(filePath string) (content string, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SplitBinlog delimited hexadecimal numbers in strings
func SplitBinlog(content string) []string {
	var hex []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		groups := strings.Split(line, " ")
		for _, group := range groups {
			if len(group) == 4 {
				hex = append(hex, group[:2])
				hex = append(hex, group[2:])
			} else {
				hex = append(hex, group[:2])
			}
		}
	}
	return hex
}
