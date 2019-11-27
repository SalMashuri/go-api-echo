package main

import (
	"database/sql"

	memberQuery "github.com/Rifannurmuhammad/go-api-echo/src/member/query"
	memberRepo "github.com/Rifannurmuhammad/go-api-echo/src/member/repository"
	memberUseCase "github.com/Rifannurmuhammad/go-api-echo/src/member/usecase"
)

//Service data structure
type Service struct {
	MemberUseCase memberUseCase.MemberUseCase
}

//MakeHandler function, Service's Constructor
func MakeHandler(readDB, writeDB *sql.DB) *Service {
	mRepoRead := memberRepo.NewMemberRepoPostgres(writeDB)
	mRepoWrite := memberRepo.NewMemberRepoPostgres(writeDB)
	mQueryRead := memberQuery.NewMemberQueryPostgres(writeDB)
	mUseCase := memberUseCase.NewMemberUsecase(mRepoRead, mRepoWrite, mQueryRead)

	return &Service{
		MemberUseCase: mUseCase,
	}
}
