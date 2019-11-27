package repository

import (
	"github.com/Rifannurmuhammad/go-api-echo/src/member/model"
)

//ResultRepository data struct
type ResultRepository struct {
	Result interface{}
	Error  error
}

//MemberRepository interface abstraction
type MemberRepository interface {
	Save(member model.Member) <-chan ResultRepository
	Load(uid string) <-chan ResultRepository
	LoadMember(uid string) ResultRepository
}
