package viewmodels

type UserForm struct {
	ID    uint
	Name  string
	Email string
}

type PageLink struct {
	Number int
	Link   string
	Active bool
}

type UserListPage struct {
	Title        string
	Users        []UserForm
	Pages        []PageLink
	HasPages     bool
	HasPrev      bool
	HasNext      bool
	PrevLink     string
	NextLink     string
	FlashSuccess string
	FlashError   string
	CsrfToken    string
}
