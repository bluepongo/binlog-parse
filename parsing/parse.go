package parsing

import (
	"fmt"
	"github.com/bluepongo/binlog-parse/event"
	"github.com/bluepongo/binlog-parse/file"
	"github.com/bluepongo/binlog-parse/util"
)

// ParseBinlog parsing the binlog
func ParseBinlog(filePath string) {
	fmt.Println("====================================binlog-parse====================================")
	content, _ := file.ReadBinlog("./binlog_test.txt")
	hex := file.SplitBinlog(content)
	binlogLen := len(hex)
	// Marks the starting point for parsing binlog
	var pos int64
	pos = 4
	for util.Int64ToInt(pos) < binlogLen {
		// parse header
		headerMap := event.ParseHeader(hex, util.Int64ToInt(pos))
		event.DisplayHeader(headerMap)

		// parse data
		dataMap := event.ParseData(hex, util.Int64ToInt(pos), headerMap)
		event.DisplayData(headerMap, dataMap)

		// parse footer
		crc := event.ParseFooter(hex, headerMap)
		event.DisplayFooter(crc)

		fmt.Println()
		// update the postion
		pos = headerMap["end_log_p"].(int64)
	}

}
