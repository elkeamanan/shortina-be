package domain

import (
	"elkeamanan/shortina/util/sql"
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Redirection string    `json:"redirection"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	TableLink = "links"

	ColumnLinkID          = "id"
	ColumnLinkKey         = "key"
	ColumnLinkRedirection = "redirection"
	ColumnLinkCreatedAt   = "created_at"
	ColumnLinkUpdatedAt   = "updated_at"
)

var (
	ColumnsLink map[sql.Operation][]string
)

func init() {
	TableColumns := []string{
		ColumnLinkID,
		ColumnLinkKey,
		ColumnLinkRedirection,
		ColumnLinkCreatedAt,
		ColumnLinkUpdatedAt,
	}

	ColumnsLink = make(map[sql.Operation][]string)

	for _, op := range []sql.Operation{sql.InsertOperation, sql.SelectOperation, sql.JoinOperation, sql.UpdateOperation} {
		switch op {
		case sql.InsertOperation:
			ColumnsLink[op] = TableColumns[:len(TableColumns)-2]
		case sql.SelectOperation:
			ColumnsLink[op] = TableColumns
		case sql.JoinOperation:
			for _, column := range TableColumns {
				ColumnsLink[op] = append(ColumnsLink[op], TableLink+"."+column)
			}
		default:
			ColumnsLink[op] = TableColumns
		}
	}
}

func GetLinkColumns(op sql.Operation) []string {
	return ColumnsLink[op]
}

func GetInsertLinkValues(link *Link) []interface{} {
	if link == nil {
		return nil
	}

	return []interface{}{
		link.ID.String(),
		link.Key,
		link.Redirection,
	}
}
