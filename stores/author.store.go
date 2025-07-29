package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/heinswanhtet/blogora-api/constants"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

const createAuthorQuery = `
INSERT INTO author (id, name, username, email, image, bio, created_at, updated_at, created_by, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (s *Store) CrateAuthor(ctx context.Context, author *types.Author) (*types.Author, error) {
	id := utils.GenerateUUID()

	_, err := s.db.ExecContext(ctx, createAuthorQuery,
		id,
		author.Name,
		author.Username,
		author.Email,
		author.Image,
		author.Bio,
		time.Now().UTC(),
		time.Now().UTC(),
		id,
		id,
	)
	if err != nil {
		return nil, err
	}

	return s.GetAuthor(ctx, id)
}

const getAuthorQuery = `
SELECT id, name, username, email, image, bio, created_at, updated_at, created_by, updated_by, deleted
FROM author
WHERE id = ? AND deleted is NULL
`

func (s *Store) GetAuthor(ctx context.Context, id string) (*types.Author, error) {
	var i types.Author
	err := s.db.QueryRowContext(ctx, getAuthorQuery, id).Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Email,
		&i.Image,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.Deleted,
	)
	if err != nil {
		return nil, fmt.Errorf("author not found")
	}
	return &i, nil
}

const getAuthorIdByEmailQuery = `
SELECT id
FROM author
WHERE email = ?
`

func (s *Store) GetAuthorIdByEmail(ctx context.Context, email string) (string, error) {
	var i types.Author

	if err := s.db.QueryRowContext(ctx, getAuthorIdByEmailQuery, email).Scan(&i.ID); err != nil {
		return "", err
	}

	if *i.ID == "" {
		return "", fmt.Errorf("author not found")
	}

	return *i.ID, nil
}

// const getAuthorsQuery = `
// SELECT *
// FROM author
// WHERE deleted is NULL
// ORDER BY ? DESC
// LIMIT ? OFFSET ?
// `

// const getTotalAuthorsQuery = `
// SELECT COUNT(*) AS total
// FROM author
// WHERE deleted is NULL
// `

func (s *Store) GetAuthors(
	ctx context.Context,
	limit, offset int,
	sort_by, sort_type, search string,
	allowedSearchList *[]string,
	otherQuery *map[string]string,
) (*[]*types.Author, int, error) {

	getTotalAuthorsQuery := `
		SELECT COUNT(*) AS total
		FROM author
		WHERE deleted is NULL
	`

	searchQuery, getTotalAuthorsQuery := utils.GetSearchQuery(
		constants.AuthorFields,
		search,
		allowedSearchList,
		otherQuery,
		getTotalAuthorsQuery,
	)

	getAuthorsQuery := fmt.Sprintf(`
		SELECT * 
		FROM author 
		WHERE deleted IS NULL %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, searchQuery, sort_by, sort_type)

	rows, err := s.db.QueryContext(ctx, getAuthorsQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	authorList := make([]*types.Author, 0)
	for rows.Next() {
		i := new(types.Author)
		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.Email,
			&i.Image,
			&i.Bio,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Deleted,
			&i.CreatedBy,
			&i.UpdatedBy,
		)
		if err != nil {
			return nil, 0, err
		}
		authorList = append(authorList, i)
	}

	var i int
	if err := s.db.QueryRowContext(ctx, getTotalAuthorsQuery).Scan(&i); err != nil {
		return nil, 0, err
	}

	return &authorList, i, nil
}

func (s *Store) UpdateAuthor(ctx context.Context, id string, updateData *types.AuthorPayload) (*types.Author, error) {
	oldData, err := s.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}

	// setClauses := []string{}
	// args := []any{}
	// if updateData.Name != nil {
	// 	setClauses = append(setClauses, "name = ?")
	// 	args = append(args, *updateData.Name)
	// }
	// if updateData.Username != nil {
	// 	setClauses = append(setClauses, "username = ?")
	// 	args = append(args, *updateData.Username)
	// }
	// if updateData.Image != nil {
	// 	setClauses = append(setClauses, "image = ?")
	// 	args = append(args, *updateData.Image)
	// }
	// if updateData.Bio != nil {
	// 	setClauses = append(setClauses, "bio = ?")
	// 	args = append(args, *updateData.Bio)
	// }
	// if len(setClauses) == 0 {
	// 	return oldData, nil
	// }
	// setClauses = append(setClauses, "updated_at = ?")
	// args = append(args, time.Now().UTC())
	// query := fmt.Sprintf("UPDATE author SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	// args = append(args, id)

	// query, args, update_ind := utils.GetSetQueryMap(
	// 	id,
	// 	"author",
	// 	[]map[string]any{
	// 		{
	// 			"field": updateData.Name,
	// 			"col":   "name",
	// 		},
	// 		{
	// 			"field": updateData.Username,
	// 			"col":   "username",
	// 		},
	// 		{
	// 			"field": updateData.Image,
	// 			"col":   "image",
	// 		},
	// 		{
	// 			"field": updateData.Bio,
	// 			"col":   "bio",
	// 		},
	// 	},
	// )

	// need to handle for updated_by
	query, args, update_ind := utils.GetSetQuery(
		id,
		"author",
		&[]*utils.SetQuery{
			utils.NewSetQuery(updateData.Name, "name"),
			utils.NewSetQuery(updateData.Username, "username"),
			utils.NewSetQuery(updateData.Image, "image"),
			utils.NewSetQuery(updateData.Bio, "bio"),
		},
	)

	if !update_ind {
		return oldData, nil
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetAuthor(ctx, id)
}

const deleteAuthorQuery = `
UPDATE author SET deleted = ?, updated_at = ?
WHERE id = ?
`

func (s *Store) DeleteAuthor(ctx context.Context, id string) error {
	author, err := s.GetAuthor(ctx, id)
	if err != nil {
		return nil // no need to return user not found
	}

	// need to handle for updated_by
	if author.Deleted == nil {
		_, err = s.db.ExecContext(ctx, deleteAuthorQuery, id, time.Now().UTC(), id)
	}

	return err
}
