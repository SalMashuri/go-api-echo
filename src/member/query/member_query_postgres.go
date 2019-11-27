package query

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Rifannurmuhammad/go-api-echo/src/member/model"
	"github.com/getsentry/raven-go"
)

// MemberQueryPostgres data structure
type MemberQueryPostgres struct {
	db *sql.DB
}

// NewMemberQueryPostgres function for initializing member query
func NewMemberQueryPostgres(db *sql.DB) *MemberQueryPostgres {
	return &MemberQueryPostgres{db: db}
}

// GetListMembers function for getting list of members
func (mq *MemberQueryPostgres) GetListMembers(params *model.Parameters) <-chan ResultQuery {

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		strQuery, queryValues := mq.generateQuery(params)

		if len(params.OrderBy) > 0 {
			params.OrderBy = fmt.Sprintf(`"%s"`, params.OrderBy)
		}

		sq := fmt.Sprintf(`SELECT id, "firstName", "lastName"
			FROM member
			%s
			ORDER BY %s %s
			LIMIT %d OFFSET %d`, strQuery, params.OrderBy, params.Sort, params.Limit, params.Offset)

		rows, err := mq.db.Query(sq, queryValues...)
		if err != nil {
			output <- ResultQuery{Error: nil}
			return
		}
		defer rows.Close()

		var members model.ListMembers
		for rows.Next() {
			var (
				member   model.Member
				lastName sql.NullString
			)

			err = rows.Scan(
				&member.ID, &member.FirstName, &lastName,
			)

			if err != nil {
				output <- ResultQuery{Error: err}
				return
			}

			// assign the nullable field to object
			if lastName.Valid {
				member.LastName = lastName.String
			}

			members.Members = append(members.Members, &member)
		}

		output <- ResultQuery{Result: members}
	}()

	return output
}

// generateQuery function for generating query
func (mq *MemberQueryPostgres) generateQuery(params *model.Parameters) (string, []interface{}) {
	var (
		strQuery, idx string
		queryStrOR    []string
		queryStrAND   []string
		queryValues   []interface{}
		lq            int
	)

	if len(params.Query) > 0 {
		queries := strings.Split(params.Query, " ")
		lq = len(queries)
		if lq > 1 {
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `("firstName" || ' ' || "lastName"  || ' ' || "id"  ilike $`+strconv.Itoa(len(queryStrOR)+1)+`)`)

		} else {
			queryStrOR = append(queryStrOR, `"firstName" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"lastName" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"id" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
		}

		idx = strconv.Itoa(len(queryStrOR) + 1)
	}

	if len(params.Status) > 0 {
		intLentOR := 0
		if len(queryStrOR) > 0 {
			intLentOR = len(queryStrOR)
		}

		idx = strconv.Itoa(intLentOR + 1)

		queryStrAND = append(queryStrAND, `status = $`+idx)
		queryValues = append(queryValues, params.Status)
	}

	if len(queryStrOR) > 0 || len(queryStrAND) > 0 {
		if len(queryStrOR) > 0 {
			strQuery = fmt.Sprintf(`(%s)`, strings.Join(queryStrOR, " OR "))
			queryStrAND = append(queryStrAND, strQuery)
		}

		if len(queryStrAND) > 0 {
			strQuery = strings.Join(queryStrAND, " AND ")
		}
	}

	if len(strQuery) > 0 {
		strQuery = fmt.Sprintf(" WHERE %s", strQuery)
	}

	return strQuery, queryValues
}

// GetTotalMembers function for getting total of members
func (mq *MemberQueryPostgres) GetTotalMembers(params *model.Parameters) <-chan ResultQuery {
	ctx := "MemberQuery-GetTotalMembers"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		var totalData int

		strQuery, queryValues := mq.generateQuery(params)

		sq := fmt.Sprintf(`SELECT count(id) FROM member %s`, strQuery)
		err := mq.db.QueryRow(sq, queryValues...).Scan(&totalData)
		if err != nil {
			raven.CaptureError(err, map[string]string{ctx: err.Error()})
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: totalData}
	}()

	return output
}
