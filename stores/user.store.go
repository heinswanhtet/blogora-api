package stores

import (
	"context"
	"fmt"

	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

const createUserQuery = `
INSERT INTO users (id, email, password, first_name, last_name)
VALUES (?, ?, ?, ?, ?)
`
const getUserQuery = `
SELECT id, email, password, first_name, last_name, created_at, updated_at
FROM users
WHERE id = ?
`

func (s *Store) CrateUser(ctx context.Context, user *types.User) (*types.User, error) {
	id := utils.GenerateUUID()

	_, err := s.db.ExecContext(ctx, createUserQuery,
		id,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
	)
	if err != nil {
		return nil, err
	}

	var i types.User
	err = s.db.QueryRowContext(ctx, getUserQuery, id).Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getUserIdByEmailQuery = `
SELECT id
FROM users
WHERE email = ?
`

func (s *Store) GetUserIdByEmail(ctx context.Context, email string) (string, error) {
	var i types.User

	if err := s.db.QueryRowContext(ctx, getUserIdByEmailQuery, email).Scan(&i.ID); err != nil {
		return "", err
	}

	if i.ID == "" {
		return "", fmt.Errorf("user not found")
	}

	return i.ID, nil
}

const getLoginUserQuery = `
SELECT id, email, password
FROM users
WHERE email = ?
`

func (s *Store) GetLoginUser(ctx context.Context, email string) (*types.User, error) {
	var i types.User

	if err := s.db.QueryRowContext(ctx, getLoginUserQuery, email).Scan(&i.ID, &i.Email, &i.Password); err != nil {
		return nil, err
	}

	return &i, nil
}
