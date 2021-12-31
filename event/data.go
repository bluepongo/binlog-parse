package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/bluepongo/binlog-parse/util"
)

var colTypeCodeMap = map[int64]string{
	0:   "DECIMAL",
	1:   "TINY",
	2:   "SHORT",
	3:   "LONG",
	4:   "FLOAT",
	5:   "DOUBLE",
	6:   "NULL",
	7:   "TIMESTAMP",
	8:   "LONGLONG",
	9:   "INT24",
	10:  "DATE",
	11:  "TIME",
	12:  "DATETIME",
	13:  "YEAR",
	14:  "NEWDATE",
	15:  "VARCHAR",
	16:  "BIT",
	17:  "TIMESTAMP2",
	18:  "DATETIME2",
	19:  "TIME2",
	20:  "TYPED_ARRAY",
	243: "INVALID",
	244: "BOOL",
	245: "JSON",
	246: "NEWDECIMAL",
	247: "ENUM",
	248: "SET",
	249: "TINY_BLOB",
	250: "MEDIUM_BLOB",
	251: "LONG_BLOB",
	252: "BLOG",
	253: "VAR_STRING",
	254: "STRING",
	255: "GEOMETRY",
}

// ParseData determine the Event type and choose the parse function
func ParseData(content []string, pos int, headerMap map[string]interface{}) map[string]interface{} {
	eventBody := content[pos+DefaultHeaderLength : headerMap["end_log_p"].(int64)-4]
	switch GetEventType(headerMap["type_code"].(int64)) {
	case "QUERY_EVENT":
		return ParseQueryEvent(eventBody)
	case "FORMAT_DESCRIPTION_EVENT":
		return ParseFormatDescriptionEvent(eventBody)
	case "XID_EVENT":
		return ParseXIDEvent(eventBody)
	case "TABLE_MAP_EVENT":
		return ParseTableMapEvent(eventBody)
	case "WRITE_EVENT":
		return ParseWriteEvent(eventBody)
	case "UPDATE_EVENT":
		return ParseUpdateEvent(eventBody)
	case "DELETE_EVENT":
		return ParseDeleteEvent(eventBody)
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

type eventParser struct {
	bodys  []string
	cursor int64
}

func (ep *eventParser) extractBody(l int64) string {
	if l == 0 {
		l = int64(len(ep.bodys)) - ep.cursor
	}
	result := ""
	for _, t := range util.EndianConversion(ep.bodys[ep.cursor : ep.cursor+l]) {
		result += t
	}
	ep.cursor += l
	return result
}

func (ep *eventParser) extractBodyToChar(l int64) string {
	if l == 0 {
		l = int64(len(ep.bodys)) - ep.cursor
	}
	result := ""
	for _, t := range ep.bodys[ep.cursor : ep.cursor+l] {
		t, _ = util.Base16ToChar(t)
		result += t
	}
	ep.cursor += l
	return result
}

func (ep *eventParser) pushCursor(len int64) {
	ep.cursor += len
}

// ParseQueryEvent parsing the query_event
func ParseQueryEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// slave_proxy_id
	slaveProxyID, _ := util.Base16ToBase10(ep.extractBody(4))
	dataMap["slave_proxy_id"] = slaveProxyID

	// query_exec_time
	queryExecTime, _ := util.Base16ToBase10(ep.extractBody(4))
	dataMap["query_exec_time"] = queryExecTime

	// db_len
	dbLen, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["db_len"] = dbLen

	// error_code
	errorCode, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["error_code"] = errorCode

	// status_vars_len
	statusVarsLen, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["status_vars_len"] = statusVarsLen

	// status variables
	dataMap["status_variables"] = ep.extractBody(statusVarsLen)

	// db
	dataMap["db"] = ep.extractBodyToChar(dbLen)
	ep.pushCursor(1)

	// query
	dataMap["query"] = ep.extractBodyToChar(0)

	return dataMap
}

// ParseFormatDescriptionEvent parsing the format_description_event
func ParseFormatDescriptionEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// binlog_version
	binlogVersion, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["binlog_version"] = binlogVersion

	// server_version
	serverVersion := ep.extractBodyToChar(50)
	serverVersion = strings.Replace(serverVersion, "\u0000", "", -1)
	dataMap["server_version"] = serverVersion

	// create_timestamp
	createTimeUnix, _ := util.Base16ToBase10(ep.extractBody(4))
	if createTimeUnix == 0 {
		dataMap["create_time"] = 0
		dataMap["create_timestamp"] = 0
	} else {
		dataMap["create_time"] = createTimeUnix
		createTimeStamp := time.Unix(createTimeUnix, 0).Format(DefaultTimeStampFormat)
		dataMap["create_timestamp"] = createTimeStamp
	}

	// header_length
	headerLength, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["header_length"] = headerLength

	// array of post-header
	dataMap["array_of_post-header"] = ep.extractBody(0)

	return dataMap
}

// ParseXIDEvent parsing the xid_event
func ParseXIDEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// XID
	xid, _ := util.Base16ToBase10(ep.extractBody(0))
	dataMap["xid"] = xid

	return dataMap
}

// ParseTableMapEvent parsing the table_map_event
func ParseTableMapEvent(eventBody []string) map[string]interface{} {
	//fmt.Println(len(eventBody))
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// table_id
	tableID, _ := util.Base16ToBase10(ep.extractBody(6))
	dataMap["table_id"] = tableID

	// Reserved
	reserved, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["reserved"] = reserved

	// db len
	dbLen, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["db_len"] = dbLen

	// db name
	dataMap["db_name"] = ep.extractBodyToChar(dbLen)
	ep.pushCursor(1)

	// table len
	tableLen, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["table_len"] = tableLen

	// table name
	dataMap["table_name"] = ep.extractBodyToChar(tableLen)
	ep.pushCursor(1)

	// no of cols
	noOfCols, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["no_of_cols"] = noOfCols

	// array of col types
	var arrayOfColTypes []string
	for _, t := range util.EndianConversion(ep.bodys[ep.cursor : ep.cursor+noOfCols]) {
		typeCode, _ := util.Base16ToBase10(t)
		arrayOfColTypes = append(arrayOfColTypes, colTypeCodeMap[typeCode])
	}
	dataMap["array_of_col_types"] = arrayOfColTypes
	ep.pushCursor(noOfCols)

	// metadata len
	metadataLen, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["metadata_len"] = metadataLen

	// metadata block
	dataMap["metadata_block"] = ep.extractBody(metadataLen)

	// m_null_bits
	mNullBitsBase, _ := util.Base16ToBase2(ep.extractBody((noOfCols + 7) / 8))
	dataMap["m_null_bits"] = mNullBitsBase

	// optional_meta_fields
	if util.Int64ToInt(ep.cursor) < len(ep.bodys) {
		var optionalMetaFieldsBase16 []interface{}
		// type
		typeCode, _ := util.Base16ToBase10(ep.extractBody(1))
		optionalMetaFieldsBase16 = append(optionalMetaFieldsBase16, typeCode)

		// length
		length, _ := util.Base16ToBase10(ep.extractBody(1))
		optionalMetaFieldsBase16 = append(optionalMetaFieldsBase16, length)

		// value
		valueBase10, _ := util.Base16ToBase10(ep.extractBody(length))
		optionalMetaFieldsBase16 = append(optionalMetaFieldsBase16, valueBase10)

		dataMap["optional_meta_fields"] = optionalMetaFieldsBase16
	}

	return dataMap
}

// ParseMetadataBlock
// TODO: parsing the meatadata block in table_map_event
func ParseMetadataBlock() {

}

// ParseWriteEvent parsing the write_event
func ParseWriteEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// table_id
	tableID, _ := util.Base16ToBase10(ep.extractBody(6))
	dataMap["table_id"] = tableID

	// Reserved
	reserved, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["reserved"] = reserved

	// var_header_len
	varHeaderLen, _ := util.Base16ToBase10("0x" + ep.extractBody(2))
	dataMap["var_header_len"] = varHeaderLen

	// columns_width
	colWidth, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["columns_width"] = colWidth

	// columns_after_image (don't consider binlog_row_image set to FULL)
	columnsAfterImage, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))
	dataMap["columns_after_image"] = columnsAfterImage

	// row Bit-field
	rowBitField, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))

	var tmp string
	zeroNo := util.Int64ToInt(dataMap["columns_width"].(int64)) - len(rowBitField)
	for i := zeroNo; i > 0; i-- {
		tmp += "0"
	}
	rowBitField = tmp + rowBitField
	dataMap["row_Bit_field"] = rowBitField

	// row real data
	var rowRealData []string
	rowRealData = ep.bodys[ep.cursor:]
	dataMap["row_real_data"] = rowRealData

	return dataMap
}

// ParseUpdateEvent parsing the update_event
// TODO: search the columns after image
func ParseUpdateEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// table_id
	tableID, _ := util.Base16ToBase10(ep.extractBody(6))
	dataMap["table_id"] = tableID

	// Reserved
	reserved, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["reserved"] = reserved

	// var_header_len
	varHeaderLen, _ := util.Base16ToBase10("0x" + ep.extractBody(2))
	dataMap["var_header_len"] = varHeaderLen

	// columns_width
	colWidth, _ := util.Base16ToBase10(ep.extractBody(1))
	dataMap["columns_width"] = colWidth

	// columns_before_image (don't consider binlog_row_image set to FULL)
	columnsBeforeImage, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))
	dataMap["columns_before_image"] = columnsBeforeImage

	// columns_after_image (don't consider binlog_row_image set to FULL)
	columnsAfterImage, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))
	dataMap["columns_after_image"] = columnsAfterImage

	// row Bit-field
	rowBitField, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))
	var tmp string
	zeroNo := util.Int64ToInt(colWidth) - len(rowBitField)
	for i := zeroNo; i > 0; i-- {
		tmp += "0"
	}
	rowBitField = tmp + rowBitField
	dataMap["row_Bit_field"] = rowBitField

	// row real data
	var rowRealData []string
	rowRealData = ep.bodys[ep.cursor:]
	dataMap["row_real_data"] = rowRealData

	return dataMap
}

// ParseDeleteEvent parsing the delete_event
func ParseDeleteEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	ep := eventParser{eventBody, 0}

	// table_id
	tableID, _ := util.Base16ToBase10(ep.extractBody(6))
	dataMap["table_id"] = tableID

	// Reserved
	reserved, _ := util.Base16ToBase10(ep.extractBody(2))
	dataMap["reserved"] = reserved

	// var_header_len
	varHeaderLen, _ := util.Base16ToBase10("0x" + ep.extractBody(2))
	dataMap["var_header_len"] = varHeaderLen

	// columns_width
	colWidth, _ := util.Base16ToBase10(ep.extractBody(0))
	dataMap["columns_width"] = colWidth

	// columns_before_image (don't consider binlog_row_image set to FULL)
	columnsBeforeImage, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))
	dataMap["columns_before_image"] = columnsBeforeImage

	// row Bit-field
	rowBitField, _ := util.Base16ToBase2(ep.extractBody((colWidth + 7) / 8))
	var tmp string
	zeroNo := util.Int64ToInt(colWidth) - len(rowBitField)
	for i := zeroNo; i > 0; i-- {
		tmp += "0"
	}
	rowBitField = tmp + rowBitField
	dataMap["row_Bit_field"] = rowBitField

	// row real data
	var rowRealData []string
	rowRealData = ep.bodys[ep.cursor:]
	dataMap["row_real_data"] = rowRealData

	return dataMap
}

// ParseRowRealData
// TODO: parsing the row real data in row_event
func ParseRowRealData() {

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
