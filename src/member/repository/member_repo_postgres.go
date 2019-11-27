package repository

import (
	"database/sql"

	"github.com/Rifannurmuhammad/go-api-echo/src/member/model"
)

//MemberRepoPostgres data structure
type MemberRepoPostgres struct {
	db *sql.DB
}

//NewMemberRepoPostgres function for initializing member repository
func NewMemberRepoPostgres(db *sql.DB) *MemberRepoPostgres {
	return &MemberRepoPostgres{db: db}
}

// Load function for loading member data based on user id
func (mr *MemberRepoPostgres) Load(uid string) <-chan ResultRepository {
	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		q := `SELECT email, "firstName", "lastName"
			FROM member WHERE id = $1`

		stmt, err := mr.db.Prepare(q)

		if err != nil {
			output <- ResultRepository{Error: err}
			defer stmt.Close()
			return
		}
		defer stmt.Close()

		// initialize needed variables
		var (
			member   model.Member
			lastName sql.NullString
		)

		err = stmt.QueryRow(uid).Scan(
			&member.FirstName, &lastName,
		)

		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}

		// assign the nullable field to object
		if lastName.Valid {
			member.LastName = lastName.String
		}
		member.ID = uid

		output <- ResultRepository{Result: member}
	}()

	return output
}

// LoadMember function without goroutine
func (mr *MemberRepoPostgres) LoadMember(uid string) ResultRepository {
	q := `SELECT "firstName", "lastName"
			FROM member WHERE id = $1`

	stmt, err := mr.db.Prepare(q)
	defer stmt.Close()

	if err != nil {
		return ResultRepository{Error: err}
	}

	// initialize needed variables
	var (
		member   model.Member
		lastName sql.NullString
	)

	err = stmt.QueryRow(uid).Scan(
		&member.FirstName, &lastName,
	)

	if err != nil {
		return ResultRepository{Error: err}
	}

	// assign the nullable field to object
	if lastName.Valid {
		member.LastName = lastName.String
	}

	return ResultRepository{Result: member}

}

// Save function for saving member data
func (mr *MemberRepoPostgres) Save(member model.Member) <-chan ResultRepository {

	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		tx, err := mr.db.Begin()
		readStmt, err := tx.Prepare(`SELECT "version" FROM member WHERE id = $1`)

		if err != nil {
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}
		defer readStmt.Close()

		var version int
		err = readStmt.QueryRow(member.ID).Scan(&version)

		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}

		q := `INSERT INTO member
				(
					id, "firstName", "lastName"
				)
			VALUES
				(
					$1, $2, $3
				)
			ON CONFLICT(id)
			DO UPDATE SET
				"firstName" = $2, "lastName" = $3`

		stmt, err := tx.Prepare(q)
		defer stmt.Close()

		if err != nil {
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}

		// set null-able variables
		var (
			lastName sql.NullString
		)

		lastName.Valid = false
		if len(member.LastName) > 0 {
			lastName.Valid = true
			lastName.String = member.LastName
		}

		_, err = stmt.Exec(
			member.ID, member.FirstName, lastName,
		)
		if err != nil {
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultRepository{Error: nil}
	}()

	return output
}
