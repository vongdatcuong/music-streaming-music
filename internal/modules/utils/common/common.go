package common_utils

func GetUInt32Pointer(value uint32) *uint32 {
	temp := value
	return &temp
}

func GetUInt64Pointer(value uint64) *uint64 {
	temp := value
	return &temp
}

func GetStringPointer(value string) *string {
	temp := value
	return &temp
}

func Contains(elems []uint64, v uint64) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
