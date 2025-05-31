package repository

import (
	"database/sql"
	"elkeamanan/shortina/internal/link/domain"
	"errors"
)

func scanLinkRows(rows *sql.Rows) ([]*domain.Link, error) {
	var links []*domain.Link

	for rows.Next() {
		var link domain.Link
		err := rows.Scan(
			&link.ID,
			&link.Key,
			&link.Redirection,
			&link.Status,
			&link.CreatedBy,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}

	return links, nil
}
