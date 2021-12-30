package event

import (
	"fmt"
	"github.com/bluepongo/binlog-parse/util"
	"strings"
	"time"
)

// ParseData determine the Event type and choose the parse function
func ParseData(content []string, pos int, headerMap map[string]interface{}) map[string]interface{} {
	eventBody := content[pos+DefaultHeaderLength : headerMap["end_log_p"].(int64)-4]
	switch EventType(headerMap["type_code"].(int64)) {
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
func ParseFormatDescriptionEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// binlog_version
	var binlogVersionBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		binlogVersionBase16 += t
	}
	pos += 2
	binlogVersion, _ := util.Base16ToBase10(binlogVersionBase16)
	dataMap["binlog_version"] = binlogVersion

	// server_version
	var serverVersion string
	for _, t := range eventBody[pos : pos+50] {
		t, _ = util.Base16ToChar(t)
		serverVersion += t
	}
	pos += 50
	serverVersion = strings.Replace(serverVersion, "\u0000", "", -1)
	dataMap["server_version"] = serverVersion

	// create_timestamp
	var createTimeBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+4]) {
		createTimeBase16 += t
	}
	createTimeUnix, _ := util.Base16ToBase10(createTimeBase16)
	if createTimeUnix == 0 {
		dataMap["create_time"] = 0
		dataMap["create_timestamp"] = 0
	} else {
		dataMap["create_time"] = createTimeUnix
		createTimeStamp := time.Unix(createTimeUnix, 0).Format(DefaultTimeStampFormat)
		dataMap["create_timestamp"] = createTimeStamp
	}

	// header_length
	var headerLengthBase16 string
	headerLengthBase16 = eventBody[pos]
	pos += 1
	headerLength, _ := util.Base16ToBase10(headerLengthBase16)
	dataMap["header_length"] = headerLength

	// array of post-header
	var arrayOfPostHeaderBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos:]) {
		arrayOfPostHeaderBase16 += t
	}
	dataMap["array_of_post-header"] = arrayOfPostHeaderBase16

	return dataMap
}

// ParseXIDEvent parsing the xid_event
func ParseXIDEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// XID
	var xidBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos:]) {
		xidBase16 += t
	}
	pos += 2
	xid, _ := util.Base16ToBase10(xidBase16)
	dataMap["xid"] = xid

	return dataMap
}

// ParseTableMapEvent parsing the table_map_event
func ParseTableMapEvent(eventBody []string) map[string]interface{} {
	//fmt.Println(len(eventBody))
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// table_id
	var tableIDBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+6]) {
		tableIDBase16 += t
	}
	pos += 6
	tableID, _ := util.Base16ToBase10(tableIDBase16)
	dataMap["table_id"] = tableID

	// Reserved
	var reservedBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		reservedBase16 += t
	}
	pos += 2
	reserved, _ := util.Base16ToBase10(reservedBase16)
	dataMap["reserved"] = reserved

	// db len
	var dbLenBase16 string
	dbLenBase16 = eventBody[pos]
	pos += 1
	dbLen, _ := util.Base16ToBase10(dbLenBase16)
	dataMap["db_len"] = dbLen

	// db name
	var dbName string
	for _, t := range eventBody[pos : pos+dataMap["db_len"].(int64)] {
		t, _ = util.Base16ToChar(t)
		dbName += t
	}
	pos = pos + dataMap["db_len"].(int64) + 1
	dataMap["db_name"] = dbName

	// table len
	var tableLenBase16 string
	tableLenBase16 = eventBody[pos]
	pos += 1
	tableLen, _ := util.Base16ToBase10(tableLenBase16)
	dataMap["table_len"] = tableLen

	// table name
	var tableName string
	for _, t := range eventBody[pos : pos+dataMap["table_len"].(int64)] {
		t, _ = util.Base16ToChar(t)
		tableName += t
	}
	pos = pos + dataMap["table_len"].(int64) + 1
	dataMap["table_name"] = tableName

	// no of cols
	var noOfColsBase16 string
	noOfColsBase16 = eventBody[pos]
	pos += 1
	noOfCols, _ := util.Base16ToBase10(noOfColsBase16)
	dataMap["no_of_cols"] = noOfCols

	// array of col types
	var arrayOfColTypes []string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+dataMap["no_of_cols"].(int64)]) {
		typeCode, _ := util.Base16ToBase10(t)
		switch typeCode {
		case 0:
			t = "DECIMAL"
		case 1:
			t = "TINY"
		case 2:
			t = "SHORT"
		case 3:
			t = "LONG"
		case 4:
			t = "FLOAT"
		case 5:
			t = "DOUBLE"
		case 6:
			t = "NULL"
		case 7:
			t = "TIMESTAMP"
		case 8:
			t = "LONGLONG"
		case 9:
			t = "INT24"
		case 10:
			t = "DATE"
		case 11:
			t = "TIME"
		case 12:
			t = "DATETIME"
		case 13:
			t = "YEAR"
		case 14:
			t = "NEWDATE"
		case 15:
			t = "VARCHAR"
		case 16:
			t = "BIT"
		case 17:
			t = "TIMESTAMP2"
		case 18:
			t = "DATETIME2"
		case 19:
			t = "TIME2"
		case 20:
			t = "TYPED_ARRAY"
		case 243:
			t = "INVALID"
		case 244:
			t = "BOOL"
		case 245:
			t = "JSON"
		case 246:
			t = "NEWDECIMAL"
		case 247:
			t = "ENUM"
		case 248:
			t = "SET"
		case 249:
			t = "TINY_BLOB"
		case 250:
			t = "MEDIUM_BLOB"
		case 251:
			t = "LONG_BLOB"
		case 252:
			t = "BLOG"
		case 253:
			t = "VAR_STRING"
		case 254:
			t = "STRING"
		case 255:
			t = "GEOMETRY"
		}
		arrayOfColTypes = append(arrayOfColTypes, t)
	}
	pos = pos + dataMap["no_of_cols"].(int64)
	dataMap["array_of_col_types"] = arrayOfColTypes

	// metadata len
	var metadataLenBase16 string
	metadataLenBase16 = eventBody[pos]
	pos += 1
	metadataLen, _ := util.Base16ToBase10(metadataLenBase16)
	dataMap["metadata_len"] = metadataLen

	// metadata block
	var metadataBlock []string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+dataMap["metadata_len"].(int64)]) {
		metadataBlock = append(metadataBlock, t)
	}
	pos += dataMap["metadata_len"].(int64)
	dataMap["metadata_block"] = metadataBlock

	// m_null_bits
	var mNullBitsBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["no_of_cols"].(int64)+7)/8)]) {
		mNullBitsBase16 += t
	}
	pos += (dataMap["no_of_cols"].(int64) + 7) / 8
	mNullBitsBase, _ := util.Base16ToBase2(mNullBitsBase16)
	dataMap["m_null_bits"] = mNullBitsBase

	// optional_meta_fields
	if util.Int64ToInt(pos) < len(eventBody) {
		var optionalMetaFieldsBase16 []interface{}
		// type
		var typeCode int64
		typeCode, _ = util.Base16ToBase10(eventBody[pos])
		pos += 1
		optionalMetaFieldsBase16 = append(optionalMetaFieldsBase16, typeCode)

		// length
		var length int64
		length, _ = util.Base16ToBase10(eventBody[pos])
		pos += 1
		optionalMetaFieldsBase16 = append(optionalMetaFieldsBase16, length)

		// value
		var value string
		for _, t := range util.ReverseSlice(eventBody[pos : pos+length]) {
			value += t
		}
		valueBase10, _ := util.Base16ToBase10(value)
		optionalMetaFieldsBase16 = append(optionalMetaFieldsBase16, valueBase10)

		dataMap["optional_meta_fields"] = optionalMetaFieldsBase16
	}

	return dataMap
}

// ParseMetadataBlock TODO parsing the meatadata block in table_map_event
func ParseMetadataBlock() {

}

// ParseWriteEvent parsing the write_event
func ParseWriteEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// table_id
	var tableIDBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+6]) {
		tableIDBase16 += t
	}
	pos += 6
	tableID, _ := util.Base16ToBase10(tableIDBase16)
	dataMap["table_id"] = tableID

	// Reserved
	var reservedBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		reservedBase16 += t
	}
	pos += 2
	reserved, _ := util.Base16ToBase10(reservedBase16)
	dataMap["reserved"] = reserved

	// var_header_len
	var varHeaderLenBase16 string
	varHeaderLenBase16 = "0x"
	for _, t := range eventBody[pos : pos+2] {
		varHeaderLenBase16 += t
	}
	pos += 2
	varHeaderLen, _ := util.Base16ToBase10(varHeaderLenBase16)
	dataMap["var_header_len"] = varHeaderLen

	// columns_width
	var colWidth int64
	colWidth, _ = util.Base16ToBase10(eventBody[pos])
	pos += 1
	dataMap["columns_width"] = colWidth

	// columns_after_image (don't consider binlog_row_image set to FULL)
	var columnsAfterImageBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		columnsAfterImageBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	columnsAfterImage, _ := util.Base16ToBase2(columnsAfterImageBase16)
	dataMap["columns_after_image"] = columnsAfterImage

	// row Bit-field
	var rowBitFieldBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		rowBitFieldBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	rowBitField, _ := util.Base16ToBase2(rowBitFieldBase16)
	var tmp string
	zeroNo := util.Int64ToInt(dataMap["columns_width"].(int64)) - len(rowBitField)
	for i := zeroNo; i > 0; i-- {
		tmp += "0"
	}
	rowBitField = tmp + rowBitField
	dataMap["row_Bit_field"] = rowBitField

	// row real data
	var rowRealData []string
	rowRealData = eventBody[pos:]
	dataMap["row_real_data"] = rowRealData

	return dataMap
}

// ParseUpdateEvent parsing the update_event TODO search the columns after image
func ParseUpdateEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// table_id
	var tableIDBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+6]) {
		tableIDBase16 += t
	}
	pos += 6
	tableID, _ := util.Base16ToBase10(tableIDBase16)
	dataMap["table_id"] = tableID

	// Reserved
	var reservedBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		reservedBase16 += t
	}
	pos += 2
	reserved, _ := util.Base16ToBase10(reservedBase16)
	dataMap["reserved"] = reserved

	// var_header_len
	var varHeaderLenBase16 string
	varHeaderLenBase16 = "0x"
	for _, t := range eventBody[pos : pos+2] {
		varHeaderLenBase16 += t
	}
	pos += 2
	varHeaderLen, _ := util.Base16ToBase10(varHeaderLenBase16)
	dataMap["var_header_len"] = varHeaderLen

	// columns_width
	var colWidth int64
	colWidth, _ = util.Base16ToBase10(eventBody[pos])
	pos += 1
	dataMap["columns_width"] = colWidth

	// columns_before_image (don't consider binlog_row_image set to FULL)
	var columnsBeforeImageBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		columnsBeforeImageBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	columnsBeforeImage, _ := util.Base16ToBase2(columnsBeforeImageBase16)
	dataMap["columns_before_image"] = columnsBeforeImage

	// columns_after_image (don't consider binlog_row_image set to FULL)
	var columnsAfterImageBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		columnsAfterImageBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	columnsAfterImage, _ := util.Base16ToBase2(columnsAfterImageBase16)
	dataMap["columns_after_image"] = columnsAfterImage

	// row Bit-field
	var rowBitFieldBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		rowBitFieldBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	rowBitField, _ := util.Base16ToBase2(rowBitFieldBase16)
	var tmp string
	zeroNo := util.Int64ToInt(dataMap["columns_width"].(int64)) - len(rowBitField)
	for i := zeroNo; i > 0; i-- {
		tmp += "0"
	}
	rowBitField = tmp + rowBitField
	dataMap["row_Bit_field"] = rowBitField

	// row real data
	var rowRealData []string
	rowRealData = eventBody[pos:]
	dataMap["row_real_data"] = rowRealData

	return dataMap
}

// ParseDeleteEvent parsing the delete_event
func ParseDeleteEvent(eventBody []string) map[string]interface{} {
	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	var pos int64
	pos = 0

	// table_id
	var tableIDBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+6]) {
		tableIDBase16 += t
	}
	pos += 6
	tableID, _ := util.Base16ToBase10(tableIDBase16)
	dataMap["table_id"] = tableID

	// Reserved
	var reservedBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+2]) {
		reservedBase16 += t
	}
	pos += 2
	reserved, _ := util.Base16ToBase10(reservedBase16)
	dataMap["reserved"] = reserved

	// var_header_len
	var varHeaderLenBase16 string
	varHeaderLenBase16 = "0x"
	for _, t := range eventBody[pos : pos+2] {
		varHeaderLenBase16 += t
	}
	pos += 2
	varHeaderLen, _ := util.Base16ToBase10(varHeaderLenBase16)
	dataMap["var_header_len"] = varHeaderLen

	// columns_width
	var colWidth int64
	colWidth, _ = util.Base16ToBase10(eventBody[pos])
	pos += 1
	dataMap["columns_width"] = colWidth

	// columns_before_image (don't consider binlog_row_image set to FULL)
	var columnsBeforeImageBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		columnsBeforeImageBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	columnsBeforeImage, _ := util.Base16ToBase2(columnsBeforeImageBase16)
	dataMap["columns_before_image"] = columnsBeforeImage

	// row Bit-field
	var rowBitFieldBase16 string
	for _, t := range util.ReverseSlice(eventBody[pos : pos+((dataMap["columns_width"].(int64)+7)/8)]) {
		rowBitFieldBase16 += t
	}
	pos += (dataMap["columns_width"].(int64) + 7) / 8
	rowBitField, _ := util.Base16ToBase2(rowBitFieldBase16)
	var tmp string
	zeroNo := util.Int64ToInt(dataMap["columns_width"].(int64)) - len(rowBitField)
	for i := zeroNo; i > 0; i-- {
		tmp += "0"
	}
	rowBitField = tmp + rowBitField
	dataMap["row_Bit_field"] = rowBitField

	// row real data
	var rowRealData []string
	rowRealData = eventBody[pos:]
	dataMap["row_real_data"] = rowRealData

	return dataMap
}

// ParseRowRealData TODO parsing the row real data in row_event
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
