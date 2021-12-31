package event

import (
	"time"

	"github.com/bluepongo/binlog-parse/util"
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
	for _, t := range util.EndianConversion(header[:4]) {
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
	for _, t := range util.EndianConversion(header[5:9]) {
		serverIDBase16 += t
	}
	serverID, _ := util.Base16ToBase10(serverIDBase16)
	headerMap["server_id"] = serverID
	// fmt.Printf("server_id: %d\t", serverID)

	// event length
	var eventLenBase16 string
	for _, t := range util.EndianConversion(header[9:13]) {
		eventLenBase16 += t
	}
	eventLen, _ := util.Base16ToBase10(eventLenBase16)
	headerMap["event_len"] = eventLen
	// fmt.Printf("event_len: %d\t", eventLen)

	// end log position
	var endLogPositionBase16 string
	for _, t := range util.EndianConversion(header[13:17]) {
		endLogPositionBase16 += t
	}
	endLogPosition, _ := util.Base16ToBase10(endLogPositionBase16)
	headerMap["end_log_p"] = endLogPosition
	// fmt.Printf("end_log_p: %d\t", endLogPosition)

	// flags
	var flagsBase16 string
	for _, t := range util.EndianConversion(header[17:19]) {
		flagsBase16 += t
	}
	flags, _ := util.Base16ToBase10(flagsBase16)
	headerMap["flags"] = flags
	// fmt.Printf("flags: %d\n", flags)
	return headerMap
}

var eventTypeMap = map[int64]string{
	2:  "QUERY_EVENT",
	15: "FORMAT_DESCRIPTION_EVENT",
	16: "XID_EVENT",
	19: "TABLE_MAP_EVENT",
	30: "WRITE_EVENT",
	31: "UPDATE_EVENT",
	32: "DELETE_EVENT",
	33: "GTID_EVENT",
	34: "ANONYMOUS_GTID_LOG_EVENT",
	35: "PREVIOUS_GTID_EVENT",
}

// GetEventType mark the event type according to type code
func GetEventType(typeCode int64) string {
	result, ok := eventTypeMap[typeCode]
	if !ok {
		return "unknown"
	}
	return result
}
