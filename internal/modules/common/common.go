package common

type NameValueInt32 struct {
	Name  string
	Value int32
}

type PaginationInfo struct {
	Offset uint64
	Limit  uint64
}
