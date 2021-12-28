package event

import (
	"fmt"
	"github.com/bluepongo/binlog-parse/util"
)

// ParseData determine the Event type and choose the parse function
func ParseData(content []string, pos int, headerMap map[string]interface{}) map[string]interface{} {
	eventBody := content[pos+DefaultHeaderLength : headerMap["end_log_p"].(int64)-4]
	switch EventType(headerMap["type_code"].(int64)) {
	case "QUERY_EVENT":
		return ParseQueryEvent(eventBody)
	case "FORMAT_DESCRIPTION_EVENT":
		ParseFormatDescriptionEvent(eventBody)
	case "XID_EVENT":
		ParseXIDEvent(eventBody)
	case "TABLE_MAP_EVENT":
		ParseTableMapEvent(eventBody)
	case "WRITE_EVENT":
		ParseWriteEvent(eventBody)
	case "UPDATE_EVENT":
		ParseUpdateEvent(eventBody)
	case "DELETE_EVENT":
		ParseDeleteEvent(eventBody)
	case "GTID_EVENT":
		ParseGtidEvent(eventBody)
	case "ANONYMOUS_GTID_LOG_EVENT":
		ParseAnonymousGtidLogEvent(eventBody)
	case "PREVIOUS_GTID_EVENT":
		ParsePreviousGtidEvent(eventBody)
	default:
		return nil
	}
	return nil
}

// ParseQueryEvent parsing the query_event
func ParseQueryEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// slave_proxy_id
	var slaveProxyIDBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+4]) {
		slaveProxyIDBase16 += t
	}
	pos += 4
	slaveProxyID, _ := util.Base16ToBase10(slaveProxyIDBase16)
	dataMap["slave_proxy_id"] = slaveProxyID

	// query_exec_time
	var queryExecTimeBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+4]) {
		queryExecTimeBase16 += t
	}
	pos += 4
	queryExecTime, _ := util.Base16ToBase10(queryExecTimeBase16)
	dataMap["query_exec_time"] = queryExecTime

	// db_len
	var dbLenBase16 string
	dbLenBase16 = eventBody[pos]
	pos += 1
	dbLen, _ := util.Base16ToBase10(dbLenBase16)
	dataMap["db_len"] = dbLen

	// error_code
	var errorCodeBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		errorCodeBase16 += t
	}
	pos += 2
	errorCode, _ := util.Base16ToBase10(errorCodeBase16)
	dataMap["error_code"] = errorCode

	// status_vars_len
	var statusVarsLenBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		statusVarsLenBase16 += t
	}
	pos += 2
	statusVarsLen, _ := util.Base16ToBase10(statusVarsLenBase16)
	dataMap["status_vars_len"] = statusVarsLen

	// status variables
	var statusVariablesBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+dataMap["status_vars_len"].(int64)]) {
		statusVariablesBase16 += t
	}
	pos += dataMap["status_vars_len"].(int64)
	dataMap["status_variables"] = statusVariablesBase16

	// db
	var db string
	for _, t := range eventBody[pos : pos+dataMap["db_len"].(int64)] {
		t, _ = util.Base16ToChar(t)
		db += t
	}
	pos = pos + dataMap["db_len"].(int64) + 1
	dataMap["db"] = db

	// query
	var query string
	for _, t := range eventBody[pos:] {
		t, _ = util.Base16ToChar(t)
		query += t
	}
	dataMap["query"] = query

	return dataMap
}

// ParseFormatDescriptionEvent parsing the format_description_event
func ParseFormatDescriptionEvent(eventBody []string) {
	// fmt.Println(eventBody)

}

// ParseXIDEvent parsing the xid_event
func ParseXIDEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParseTableMapEvent parsing the table_map_event
func ParseTableMapEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParseWriteEvent parsing the write_event
func ParseWriteEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParseUpdateEvent parsing the update_event
func ParseUpdateEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParseDeleteEvent parsing the delete_event
func ParseDeleteEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParseGtidEvent parsing the gtid_event
func ParseGtidEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParseAnonymousGtidLogEvent parsing the anonymous_gtid_log_event
func ParseAnonymousGtidLogEvent(eventBody []string) {
	fmt.Println(eventBody)
}

// ParsePreviousGtidEvent parsing the previous_gtid_event
func ParsePreviousGtidEvent(eventBody []string) {
	fmt.Println(eventBody)
}
