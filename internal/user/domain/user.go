package domain

import (
	"elkeamanan/shortina/util"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider string

const (
	LocalProvider  AuthProvider = "local"
	GoogleProvider AuthProvider = "google"
)

type User struct {
	ID        uuid.UUID    `json:"id"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	Fullname  string       `json:"fullname"`
	Provider  AuthProvider `json:"provider"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

var (
	TableUsers = "users"

	ColumnUsersID        = "id"
	ColumnUsersEmail     = "email"
	ColumnUsersPassword  = "password"
	ColumnUsersFullname  = "fullname"
	ColumnUsersProvider  = "provider"
	ColumnUsersCreatedAt = "created_at"
	ColumnUsersUpdatedAt = "updated_at"
)

var (
	ColumnsUsers map[util.SQLOperation][]string
)

func init() {
	TableColumns := []string{
		ColumnUsersID,
		ColumnUsersEmail,
		ColumnUsersPassword,
		ColumnUsersFullname,
		ColumnUsersProvider,
		ColumnUsersCreatedAt,
		ColumnUsersUpdatedAt,
	}

	ColumnsUsers = make(map[util.SQLOperation][]string)

	for _, op := range []util.SQLOperation{util.InsertOperation, util.SelectOperation, util.JoinOperation, util.UpdateOperation} {
		switch op {
		case util.InsertOperation:
			ColumnsUsers[op] = TableColumns[:len(TableColumns)-2]
		case util.SelectOperation:
			ColumnsUsers[op] = TableColumns
		case util.JoinOperation:
			for _, column := range TableColumns {
				ColumnsUsers[op] = append(ColumnsUsers[op], TableUsers+"."+column)
			}
		default:
			ColumnsUsers[op] = TableColumns
		}
	}
}

func GetUsersColumns(op util.SQLOperation) []string {
	return ColumnsUsers[op]
}

func GetInsertUsersValues(user *User) []interface{} {
	return []interface{}{
		user.ID,
		user.Email,
		user.Password,
		user.Fullname,
		user.Provider,
	}
}

func GetUpdateUsersMap(user *User) map[string]interface{} {
	result := make(map[string]interface{})

	if user.Fullname != "" {
		result[ColumnUsersFullname] = user.Fullname
	}

	return result
}

type UserPredicate struct {
	ID       uuid.UUID
	Email    string
	Provider AuthProvider
}

func (p UserPredicate) ToWherePredicate() sq.Sqlizer {
	result := sq.And{}

	if p.ID != uuid.Nil {
		result = append(result, sq.Eq{ColumnUsersID: p.ID.String()})
	}

	if p.Email != "" {
		result = append(result, sq.Eq{ColumnUsersEmail: p.Email})
	}

	if p.Provider != "" {
		result = append(result, sq.Eq{ColumnUsersProvider: p.Provider})
	}

	return result
}

func ValidatePassword(existingPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(existingPassword, password)
	return err == nil
}
