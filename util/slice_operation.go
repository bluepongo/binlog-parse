package util

// EndianConversion convert BigEndian/LittleEndian to another
func EndianConversion(slc []string) []string {
	for i, j := 0, len(slc)-1; i < j; i, j = i+1, j-1 {
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}
