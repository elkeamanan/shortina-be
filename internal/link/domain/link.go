package domain

import (
	"database/sql"
	"elkeamanan/shortina/util"
	"time"

	"github.com/google/uuid"
)

type LinkStatus string

const (
	LinkStatusUnspecified LinkStatus = "unspecified"
	LinkStatusActive      LinkStatus = "active"
	LinkStatusInactive    LinkStatus = "inactive"
)

type Link struct {
	ID          uuid.UUID  `json:"id"`
	Key         string     `json:"key"`
	Redirection string     `json:"redirection"`
	Status      LinkStatus `json:"status"`
	CreatedBy   *uuid.UUID `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

var (
	TableLink = "links"

	ColumnLinkID          = "id"
	ColumnLinkKey         = "key"
	ColumnLinkRedirection = "redirection"
	ColumnLinkStatus      = "status"
	ColumnLinkCreatedBy   = "created_by"
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
		ColumnLinkStatus,
		ColumnLinkCreatedBy,
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

func GetInsertLinkValues(link *Link) []any {
	if link == nil {
		return nil
	}

	var createdBy sql.NullString
	if link.CreatedBy != nil {
		createdBy = sql.NullString{
			String: link.CreatedBy.String(),
			Valid:  true,
		}
	}

	return []any{
		link.ID.String(),
		link.Key,
		link.Redirection,
		link.Status,
		createdBy,
	}
}
