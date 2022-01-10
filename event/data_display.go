package event

import (
	"fmt"
	"github.com/bluepongo/binlog-parse/util"
)

// DisplayData determine the Event type and choose the display function
func DisplayData(headerMap map[string]interface{}, dataMap map[string]interface{}) {

	switch GetEventType(headerMap["type_code"].(int64)) {
	case "QUERY_EVENT":
		DisplayQueryEvent(dataMap)
	case "FORMAT_DESCRIPTION_EVENT":
		DisplayFormatDescriptionEvent(dataMap)
	case "XID_EVENT":
		DisplayXIDEvent(dataMap)
	case "TABLE_MAP_EVENT":
		DisplayTableMapEvent(dataMap)
	case "WRITE_EVENT":
		DisplayWriteEvent(dataMap)
	case "UPDATE_EVENT":
		DisplayUpdateEvent(dataMap)
	case "DELETE_EVENT":
		DisplayDeleteEvent(dataMap)
	case "GTID_EVENT":
		DisplayGtidEvent(dataMap)
	case "ANONYMOUS_GTID_LOG_EVENT":
		DisplayAnonymousGtidEvent(dataMap)
	case "PREVIOUS_GTID_LOG_EVENT":
		DisplayPreviousGtidsEvent(dataMap)
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
	fmt.Printf("table id: %v\t", dataMap["table_id"])
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

// DisplayWriteEvent  display the write_event info
func DisplayWriteEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("table_id: %v\t", dataMap["table_id"])
	fmt.Printf("row Bit-field: %v\n", dataMap["table_id"])
	fmt.Printf("row real data: %v\n", dataMap["row_real_data"])
}

// DisplayUpdateEvent  display the update_event info
func DisplayUpdateEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("table_id: %v\t", dataMap["table_id"])
	fmt.Printf("row Bit-field: %v\n", dataMap["table_id"])
	fmt.Printf("row real data: %v\n", dataMap["row_real_data"])
}

// DisplayDeleteEvent  display the delete_event info
func DisplayDeleteEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("table_id: %v\n", dataMap["table_id"])
	fmt.Printf("row Bit-field: %v\n", dataMap["table_id"])
	fmt.Printf("row real data: %v\n", dataMap["row_real_data"])
}

// DisplayGtidEvent display the gtid_event info
func DisplayGtidEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("bin-log type: %v\t", dataMap["flags"])
	fmt.Printf("server_uuid: %v:%v\t", dataMap["server_uuid"], dataMap["gno"])
	fmt.Printf("last committed: %v\t", dataMap["last_commit"])
	fmt.Printf("seq number: %v\n", dataMap["seq_number"])
}

// DisplayAnonymousGtidEvent display the gtid_event info
func DisplayAnonymousGtidEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	fmt.Printf("bin-log type: %v\t", dataMap["flags"])
	fmt.Printf("server_uuid: %v:%v\t", dataMap["server_uuid"], dataMap["gno"])
	fmt.Printf("last committed: %v\t", dataMap["last_commit"])
	fmt.Printf("seq number: %v\n", dataMap["seq_number"])
}

// DisplayPreviousGtidsEvent display the previous gtids event info
func DisplayPreviousGtidsEvent(dataMap map[string]interface{}) {
	fmt.Printf("[data body]\n")
	if dataMap["num_of_sids"].(int64) == 0 {
		fmt.Println("empty")
	}
	for m := 0; m < util.Int64ToInt(dataMap["num_of_sids"].(int64)); m++ {
		for n := 0; n < util.Int64ToInt(dataMap["sids"].([]map[string]interface{})[m]["n_intervals"].(int64)); n++ {
			interStart := dataMap["sids"].([]map[string]interface{})[m]["inter_start_next"].([][2]int64)[n][0]
			interEnd := dataMap["sids"].([]map[string]interface{})[m]["inter_start_next"].([][2]int64)[n][1] - 1
			fmt.Printf("%v: %v-%v\n",
				dataMap["sids"].([]map[string]interface{})[m]["server_uuid"], interStart, interEnd)
		}
	}
}
