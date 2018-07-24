package db

type stat struct {
	PageSize          int
	Depth             int
	BranchPageCount   int
	LeafPageCount     int
	OverflowPageCount int
	EntryCount        int
}
