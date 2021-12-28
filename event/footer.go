package event

import (
	"github.com/bluepongo/binlog-parse/util"
)

func ParseFooter(content []string, headerMap map[string]interface{}) string {
	var crc string
	for _, t := range util.ReverseSlice(content[headerMap["end_log_p"].(int64)-4 : headerMap["end_log_p"].(int64)]) {
		crc += t
	}
	return crc
}
