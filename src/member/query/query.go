package query

import "github.com/Rifannurmuhammad/go-api-echo/src/member/model"

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// MemberQuery interface abstraction
type MemberQuery interface {
	GetListMembers(params *model.Parameters) <-chan ResultQuery
	GetTotalMembers(params *model.Parameters) <-chan ResultQuery
}
