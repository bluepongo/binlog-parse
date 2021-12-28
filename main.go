package main

import (
	"github.com/bluepongo/binlog-parse/parsing"
)

const (
	DefaultBinlogFilePath = "./binlog_test.txt"
)

func main() {
	//unix := "61c57ff2"
	//t, _ := parsing.Base16ToAscii(unix)
	//fmt.Println(t)
	//fmt.Println(time.Unix(t, 0))
	parsing.ParseBinlog(DefaultBinlogFilePath)
}
