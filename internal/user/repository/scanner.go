package repository

import (
	"database/sql"
	"elkeamanan/shortina/internal/user/domain"
	"errors"
)

func scanUserRows(rows *sql.Rows) ([]*domain.User, error) {
	users := []*domain.User{}
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Fullname,
			&user.Provider,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
