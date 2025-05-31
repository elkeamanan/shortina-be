package domain

import (
	"elkeamanan/shortina/util"
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
	ColumnsLink map[util.SQLOperation][]string
)

func init() {
	TableColumns := []string{
		ColumnLinkID,
		ColumnLinkKey,
		ColumnLinkRedirection,
		ColumnLinkCreatedAt,
		ColumnLinkUpdatedAt,
	}

	ColumnsLink = make(map[util.SQLOperation][]string)

	for _, op := range []util.SQLOperation{util.InsertOperation, util.SelectOperation, util.JoinOperation, util.UpdateOperation} {
		switch op {
		case util.InsertOperation:
			ColumnsLink[op] = TableColumns[:len(TableColumns)-2]
		case util.SelectOperation:
			ColumnsLink[op] = TableColumns
		case util.JoinOperation:
			for _, column := range TableColumns {
				ColumnsLink[op] = append(ColumnsLink[op], TableLink+"."+column)
			}
		default:
			ColumnsLink[op] = TableColumns
		}
	}
}

func GetLinkColumns(op util.SQLOperation) []string {
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
