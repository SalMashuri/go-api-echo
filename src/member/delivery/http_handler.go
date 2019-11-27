package delivery

import (
	"errors"
	"math"
	"net/http"

	"github.com/Rifannurmuhammad/go-api-echo/src/member/model"
	"github.com/Rifannurmuhammad/go-api-echo/src/member/usecase"
	"github.com/Rifannurmuhammad/go-api-echo/src/shared"
	"github.com/Rifannurmuhammad/go-api-echo/helper"
	"github.com/labstack/echo"
)

// HTTPMemberHandler  represent the httphandler for article
type HTTPMemberHandler struct {
	MemberUseCase usecase.MemberUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(memberUseCase usecase.MemberUseCase) *HTTPMemberHandler {
	return &HTTPMemberHandler{MemberUseCase: memberUseCase}
}

// MountMe function for mounting routes
func (handler *HTTPMemberHandler) MountMe(group *echo.Group) {
	group.GET("", handler.FetchMember)
}

// FetchMember function for getting member
func (h *HTTPMemberHandler) FetchMember(c echo.Context) error {

	params := model.Parameters{
		Query:    c.QueryParam("query"),
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		Sort:     c.QueryParam("sort"),
		OrderBy:  c.QueryParam("orderBy"),
		Status:   c.QueryParam("status"),
	}

	memberResult := <-h.MemberUseCase.GetListMembers(&params)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	member, ok := memberResult.Result.(model.ListMembers)
	if !ok {
		err := errors.New("result is not list of members")
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	member.ID = helper.RandomString(8)
	member.Name = "list of members"

	totalPage := math.Ceil(float64(member.TotalData) / float64(params.Limit))

	if len(member.Members) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, "Get Members Response", make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: member.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, "Get Members Response", member.Members, meta).JSON(c)
}
