package model

// Member data structure
type Member struct {
	ID        string `jsonapi:"primary" json:"id"`
	FirstName string `jsonapi:"attr,firstName" json:"firstName"`
	LastName  string `jsonapi:"attr,lastName" json:"lastName"`
}

// MemberError data structure
type MemberError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// Parameters data structure
type Parameters struct {
	Query    string
	StrPage  string
	Page     int
	StrLimit string
	Limit    int
	Offset   int
	Status   string
	Sort     string
	OrderBy  string
}

// ListMembers data structure
type ListMembers struct {
	ID        string    `jsonapi:"primary,members" json:"members"`
	Name      string    `jsonapi:"attr,name" json:"name"`
	Members   []*Member `jsonapi:"relation,member" json:"member"`
	TotalData int       `json:"totalData"`
}
