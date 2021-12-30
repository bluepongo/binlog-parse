package event

import "fmt"

// DisplayHeader display the event header info
func DisplayHeader(headerMap map[string]interface{}) {
	eventType := GetEventType(headerMap["type_code"].(int64))
	fmt.Printf("%v\n", headerMap["timestamp"])
	fmt.Printf("[%s]\n", eventType)
	fmt.Printf("server_id: %v\t", headerMap["server_id"])
	fmt.Printf("event_len: %v\t", headerMap["event_len"])
	fmt.Printf("next_event_position: %v\t", headerMap["end_log_p"])
	fmt.Printf("flags: %v\n", headerMap["flags"])
}
