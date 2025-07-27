package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/heinswanhtet/blogora-api/constants"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

const createStartupQuery = `
INSERT INTO startup (id, title, slug, author_id, description, category, image, pitch, created_at, updated_at, created_by, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (s *Store) CrateStartup(ctx context.Context, startup *types.Startup) (*types.Startup, error) {
	id := utils.GenerateUUID()

	_, err := s.db.ExecContext(ctx, createStartupQuery,
		id,
		startup.Title,
		startup.Slug,
		startup.AuthorId,
		startup.Description,
		startup.Category,
		startup.Image,
		startup.Pitch,
		time.Now().UTC(),
		time.Now().UTC(),
		startup.CreatedBy,
		startup.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return s.GetStartup(ctx, id)
}

const getStartupQuery = `
SELECT *
FROM startup
WHERE id = ? AND deleted is NULL
`

func (s *Store) GetStartup(ctx context.Context, id string) (*types.Startup, error) {
	var i types.Startup
	err := s.db.QueryRowContext(ctx, getStartupQuery, id).Scan(
		&i.ID,
		&i.Title,
		&i.Slug,
		&i.AuthorId,
		&i.Views,
		&i.Description,
		&i.Category,
		&i.Image,
		&i.Pitch,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("startup not found")
	}
	return &i, nil
}

const getStartupBySlugQuery = `
SELECT *
FROM startup
WHERE slug = ?
`

func (s *Store) GetStartupBySlug(ctx context.Context, slug string) (*types.Startup, error) {
	var i types.Startup
	err := s.db.QueryRowContext(ctx, getStartupBySlugQuery, slug).Scan(
		&i.ID,
		&i.Title,
		&i.Slug,
		&i.AuthorId,
		&i.Views,
		&i.Description,
		&i.Category,
		&i.Image,
		&i.Pitch,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return &i, err
}

var validStartupSearchFields []string = append([]string{"author.name", "author.email"}, constants.StartupFields...)

func (s *Store) GetStartups(
	ctx context.Context,
	limit, offset int,
	sort_by, sort_type, search string,
	allowedSearchList *[]string,
	otherQuery *map[string]string,
) (*[]*types.Startup, int, error) {

	getTotalStartupsQuery := `
		SELECT COUNT(*) AS total
		FROM startup
		INNER JOIN author ON startup.author_id = author.id
		WHERE startup.deleted is NULL
	`

	searchQuery, getTotalStartupsQuery := utils.GetSearchQuery(
		validStartupSearchFields,
		search,
		allowedSearchList,
		otherQuery,
		getTotalStartupsQuery,
	)

	getStartupsQuery := fmt.Sprintf(`
		SELECT startup.id, title, slug, author_id, views, description, category, startup.image, pitch, startup.created_at, startup.updated_at, startup.created_by, startup.updated_by, startup.deleted, author.id, author.name, author.username, author.email, author.image
		FROM startup
		INNER JOIN author ON startup.author_id = author.id
		WHERE startup.deleted IS NULL %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, searchQuery, sort_by, sort_type)

	rows, err := s.db.QueryContext(ctx, getStartupsQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	startupList := make([]*types.Startup, 0)
	for rows.Next() {
		i := new(types.Startup)
		err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Slug,
			&i.AuthorId,
			&i.Views,
			&i.Description,
			&i.Category,
			&i.Image,
			&i.Pitch,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.Deleted,
			&i.Author.ID,
			&i.Author.Name,
			&i.Author.Username,
			&i.Author.Email,
			&i.Author.Image,
		)
		if err != nil {
			return nil, 0, err
		}
		startupList = append(startupList, i)
	}

	var i int
	if err := s.db.QueryRowContext(ctx, getTotalStartupsQuery).Scan(&i); err != nil {
		return nil, 0, err
	}

	return &startupList, i, nil
}

func (s *Store) UpdateStartup(ctx context.Context, id string, updateData *types.StartupPayload) (*types.Startup, error) {
	oldData, err := s.GetStartup(ctx, id)
	if err != nil {
		return nil, err
	}

	jwtPayload, err := utils.GetJWTPayload(ctx)
	if err != nil {
		return nil, err
	}

	query, args, update_ind := utils.GetSetQuery(
		id,
		"startup",
		&[]*utils.SetQuery{
			utils.NewSetQuery(updateData.Title, "title"),
			utils.NewSetQuery(updateData.Description, "description"),
			utils.NewSetQuery(updateData.Category, "category"),
			utils.NewSetQuery(updateData.Image, "image"),
			utils.NewSetQuery(updateData.Pitch, "pitch"),
			utils.NewSetQuery(&jwtPayload.UserId, "updated_by"),
		},
	)

	if !update_ind {
		return oldData, nil
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	newData, err := s.GetStartup(ctx, id)
	if err != nil {
		return nil, err
	}

	if *newData.Title != *oldData.Title {
		slug := utils.GenerateUniqueSlug(*newData.Title, func(slug string) bool {
			_, err := s.GetStartupBySlug(ctx, slug)
			return err == nil
		})

		query, args, _ := utils.GetSetQuery(
			id,
			"startup",
			&[]*utils.SetQuery{
				utils.NewSetQuery(&slug, "slug"),
			},
		)

		_, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		return s.GetStartup(ctx, id)
	}

	return newData, nil
}

const deleteStartupQuery = `
UPDATE startup SET deleted = ?, updated_at = ?, updated_by = ?
WHERE id = ?
`

func (s *Store) DeleteStartup(ctx context.Context, id string) error {
	startup, err := s.GetStartup(ctx, id)
	if err != nil {
		return err
	}

	if startup.Deleted == nil {
		jwtPayload, err := utils.GetJWTPayload(ctx)
		if err != nil {
			return err
		}
		_, err = s.db.ExecContext(ctx, deleteStartupQuery, id, time.Now().UTC(), jwtPayload.UserId, id)
		return err
	}

	return nil
}
