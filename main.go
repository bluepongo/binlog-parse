package main

import "github.com/bluepongo/binlog-parse/parsing"

const (
	DefaultBinlogFilePath = "./binlog_test.txt"
)

func main() {

	parsing.ParseBinlog(DefaultBinlogFilePath)

}
