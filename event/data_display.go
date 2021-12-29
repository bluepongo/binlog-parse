package event

import "fmt"

// DisplayData determine the Event type and choose the display function
func DisplayData(headerMap map[string]interface{}, dataMap map[string]interface{}) {

	switch EventType(headerMap["type_code"].(int64)) {
	case "QUERY_EVENT":
		DisplayQueryEvent(dataMap)
	case "FORMAT_DESCRIPTION_EVENT":
		DisplayFormatDescriptionEvent(dataMap)
	case "XID_EVENT":
		DisplayXIDEvent(dataMap)
	case "TABLE_MAP_EVENT":
		DisplayTableMapEvent(dataMap)
	}
}

// DisplayQueryEvent display the query_event info
func DisplayQueryEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("thread_id: %v\t", dataMap["slave_proxy_id"])
	fmt.Printf("query_exec_time: %v seconds\t", dataMap["query_exec_time"])
	//fmt.Printf("db_len: %v\t", dataMap["db_len"])
	fmt.Printf("error_code: %v\t", dataMap["error_code"])
	//fmt.Printf("status_vars_len: %v\t", dataMap["status_vars_len"])
	fmt.Printf("db: %v\t", dataMap["db"])
	fmt.Printf("query: %v\n", dataMap["query"])
	fmt.Printf("status variables: %v\n", dataMap["status_variables"])
}

// DisplayFormatDescriptionEvent display the format_description_event info
func DisplayFormatDescriptionEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("binlog_version: %v\t", dataMap["binlog_version"])
	fmt.Printf("server_version: %v\t", dataMap["server_version"])
	fmt.Printf("create_timestamp: %v\t", dataMap["create_timestamp"])
	fmt.Printf("header_length: %v\n", dataMap["header_length"])

	fmt.Printf("array of post-header: %v\n", dataMap["array_of_post-header"])
}

// DisplayXIDEvent display the xid_event info
func DisplayXIDEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("Xid: %v\n", dataMap["xid"])
}

// DisplayTableMapEvent display the table_map_event info
func DisplayTableMapEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("db name: %v\t", dataMap["db_name"])
	fmt.Printf("table name: %v\t", dataMap["table_name"])

	//fmt.Printf("array of col types: %v\t", dataMap["array_of_col_types"])
	//fmt.Printf("metadata block: %v\t", dataMap["metadata_block"])
	fmt.Printf("mapped to number: %v\n", dataMap["m_null_bits"])
	if dataMap["optional_meta_fields"] != nil {
		fmt.Printf("optional meta fields: [type=%v value=%v]\n",
			dataMap["optional_meta_fields"].([]interface{})[0],
			//dataMap["optional_meta_fields"].([]interface{})[1],
			dataMap["optional_meta_fields"].([]interface{})[2],
		)
	}

}
