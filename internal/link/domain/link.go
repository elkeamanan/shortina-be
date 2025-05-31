package domain

import (
	"database/sql"
	"elkeamanan/shortina/util"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func GetUpdateLinkMap(link *Link) map[string]any {
	result := map[string]any{}

	if link.Key != "" {
		result[ColumnLinkKey] = link.Key
	}

	if link.Redirection != "" {
		result[ColumnLinkRedirection] = link.Redirection
	}

	if link.Status != "" && link.Status != LinkStatusUnspecified {
		result[ColumnLinkStatus] = link.Status
	}

	return result
}

type LinkPredicate struct {
	ID     string
	Key    string
	UserID string
	Status LinkStatus
}

func (p LinkPredicate) ToWherePredicate() sq.Sqlizer {
	result := sq.And{}

	if p.ID != "" {
		result = append(result, sq.Eq{ColumnLinkID: p.ID})
	}

	if p.Key != "" {
		result = append(result, sq.Eq{ColumnLinkKey: p.Key})
	}

	if p.UserID != "" {
		result = append(result, sq.Eq{ColumnLinkCreatedBy: p.UserID})
	}

	if p.Status != "" && p.Status != LinkStatusUnspecified {
		result = append(result, sq.Eq{ColumnLinkStatus: p.Status})
	}

	return result
}
