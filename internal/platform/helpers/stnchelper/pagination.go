package stnchelper

type PageLink struct {
	Number int
	Link   string
	Active bool
}
type Pagination struct {
	TotalItems  int64
	TotalPages  int
	CurrentPage int
	NextPage    int
	PrevPage    int
	Pages       []int
	HasNext     bool
	HasPrev     bool
	PageList    []int
}
