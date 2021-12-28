package event

import (
	"github.com/bluepongo/binlog-parse/util"
	"time"
)

const (
	DefaultHeaderLength    = 19
	DefaultTimeStampFormat = "2006-01-02 15:04:05"
)

// ParseHeader paring the event header
func ParseHeader(content []string, pos int) map[string]interface{} {
	header := content[pos : pos+DefaultHeaderLength]
	var headerMap map[string]interface{}
	headerMap = make(map[string]interface{})

	// timestamp
	var timeBase16 string
	for _, t := range util.ReverseSlice(header[:4]) {
		timeBase16 += t
	}
	timeUnix, _ := util.Base16ToBase10(timeBase16)
	headerMap["time"] = timeUnix
	timeStamp := time.Unix(timeUnix, 0).Format(DefaultTimeStampFormat)
	headerMap["timestamp"] = timeStamp
	// fmt.Printf("timestamp: %s\t", timeStamp)

	// type code
	typeCodeBase16 := header[4]
	typeCode, _ := util.Base16ToBase10(typeCodeBase16)
	headerMap["type_code"] = typeCode
	// fmt.Printf("type code: %d\t", typeCode)

	// server id
	var serverIDBase16 string
	for _, t := range util.ReverseSlice(header[5:9]) {
		serverIDBase16 += t
	}
	serverID, _ := util.Base16ToBase10(serverIDBase16)
	headerMap["server_id"] = serverID
	// fmt.Printf("server_id: %d\t", serverID)

	// event length
	var eventLenBase16 string
	for _, t := range util.ReverseSlice(header[9:13]) {
		eventLenBase16 += t
	}
	eventLen, _ := util.Base16ToBase10(eventLenBase16)
	headerMap["event_len"] = eventLen
	// fmt.Printf("event_len: %d\t", eventLen)

	// end log position
	var endLogPositionBase16 string
	for _, t := range util.ReverseSlice(header[13:17]) {
		endLogPositionBase16 += t
	}
	endLogPosition, _ := util.Base16ToBase10(endLogPositionBase16)
	headerMap["end_log_p"] = endLogPosition
	// fmt.Printf("end_log_p: %d\t", endLogPosition)

	// flags
	var flagsBase16 string
	for _, t := range util.ReverseSlice(header[17:19]) {
		flagsBase16 += t
	}
	flags, _ := util.Base16ToBase10(flagsBase16)
	headerMap["flags"] = flags
	// fmt.Printf("flags: %d\n", flags)
	return headerMap
}

// EventType mark the event type according to type code
func EventType(typeCode int64) string {
	switch typeCode {
	case 2:
		return "QUERY_EVENT"
	case 15:
		return "FORMAT_DESCRIPTION_EVENT"
	case 16:
		return "XID_EVENT"
	case 19:
		return "TABLE_MAP_EVENT"
	case 30:
		return "WRITE_EVENT"
	case 31:
		return "UPDATE_EVENT"
	case 32:
		return "DELETE_EVENT"
	case 33:
		return "GTID_EVENT"
	case 34:
		return "ANONYMOUS_GTID_LOG_EVENT"
	case 35:
		return "PREVIOUS_GTID_EVENT"
	default:
		return "unknown"
	}
}
