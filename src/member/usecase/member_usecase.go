package usecase

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Rifannurmuhammad/go-api-echo/helper"
	"github.com/Rifannurmuhammad/go-api-echo/src/member/model"
	"github.com/Rifannurmuhammad/go-api-echo/src/member/query"
	"github.com/Rifannurmuhammad/go-api-echo/src/member/repository"
)

// MemberUsecase data structure
type MemberUsecaseImpl struct {
	MemberRepoRead  repository.MemberRepository
	MemberRepoWrite repository.MemberRepository
	MemberQueryRead query.MemberQuery
}

// NewMemberUsecase will create new an MemberUsecase object representation of memberUsecase interface
func NewMemberUsecase(memberRepoRead, memberRepoWrite repository.MemberRepository, memberQueryRead query.MemberQuery) MemberUseCase {
	return &MemberUsecaseImpl{
		MemberRepoRead:  memberRepoRead,
		MemberRepoWrite: memberRepoWrite,
		MemberQueryRead: memberQueryRead,
	}
}

// GetListMembers function for getting all list of members
func (me *MemberUsecaseImpl) GetListMembers(params *model.Parameters) <-chan ResultUsecase {
	output := make(chan ResultUsecase)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				output <- ResultUsecase{HTTPStatus: http.StatusInternalServerError, Error: err}
			}
			close(output)
		}()

		memberResult := <-me.MemberQueryRead.GetListMembers(params)
		if memberResult.Error != nil {
			httpStatus := http.StatusInternalServerError
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "member")
			}

			output <- ResultUsecase{Error: memberResult.Error, HTTPStatus: httpStatus}
		}

		member := memberResult.Result.(model.ListMembers)

		totalResult := <-me.MemberQueryRead.GetTotalMembers(params)
		if totalResult.Error != nil {
			output <- ResultUsecase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		total := totalResult.Result.(int)
		member.TotalData = total
		output <- ResultUsecase{Result: member}
	}()
	return output
}
