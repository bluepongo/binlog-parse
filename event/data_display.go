package event

import "fmt"

// DisplayData determine the Event type and choose the display function
func DisplayData(headerMap map[string]interface{}, dataMap map[string]interface{}) {

	switch EventType(headerMap["type_code"].(int64)) {
	case "QUERY_EVENT":
		DisplayQueryEvent(dataMap)
	}
}

// DisplayQueryEvent display the query_event
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
