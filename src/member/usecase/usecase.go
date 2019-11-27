package usecase

import (
	"github.com/Rifannurmuhammad/go-api-echo/src/member/model"
)

// ResultUsecase data structure
type ResultUsecase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.MemberError
}

// MemberUseCase interface abstraction
type MemberUseCase interface {
	GetListMembers(data *model.Parameters) <-chan ResultUsecase
}
