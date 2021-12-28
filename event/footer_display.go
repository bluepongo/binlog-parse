package event

import "fmt"

// DisplayFooter display the crc code
func DisplayFooter(crc string) {
	fmt.Printf("crc: 0x%s\n", crc)
}
