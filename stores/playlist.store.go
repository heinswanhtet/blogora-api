package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/heinswanhtet/blogora-api/constants"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

const createPlaylistQuery = `
INSERT INTO playlist (id, title, slug, created_at, updated_at, created_by, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

func (s *Store) CratePlaylist(ctx context.Context, playlist *types.Playlist) (*types.Playlist, error) {
	id := utils.GenerateUUID()

	_, err := s.db.ExecContext(ctx, createPlaylistQuery,
		id,
		playlist.Title,
		playlist.Slug,
		time.Now().UTC(),
		time.Now().UTC(),
		playlist.CreatedBy,
		playlist.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return s.GetPlaylist(ctx, id)
}

const getPlaylistQuery = `
SELECT *
FROM playlist
WHERE id = ? AND deleted is NULL
`

func (s *Store) GetPlaylist(ctx context.Context, id string) (*types.Playlist, error) {
	var i types.Playlist
	err := s.db.QueryRowContext(ctx, getPlaylistQuery, id).Scan(
		&i.ID,
		&i.Title,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("playlist not found")
	}
	return &i, nil
}

const getPlaylistBySlugQuery = `
SELECT *
FROM playlist
WHERE slug = ?
`

func (s *Store) GetPlaylistBySlug(ctx context.Context, slug string) (*types.Playlist, error) {
	var i types.Playlist
	err := s.db.QueryRowContext(ctx, getPlaylistBySlugQuery, slug).Scan(
		&i.ID,
		&i.Title,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return &i, err
}

func (s *Store) GetPlaylists(
	ctx context.Context,
	limit, offset int,
	sort_by, sort_type, search string,
	allowedSearchList *[]string,
	otherQuery *map[string]string,
) (*[]*types.Playlist, int, error) {

	getTotalPlaylistsQuery := `
		SELECT COUNT(*) AS total
		FROM playlist
		WHERE deleted is NULL
	`

	searchQuery, getTotalPlaylistsQuery := utils.GetSearchQuery(
		constants.PlaylistFields,
		search,
		allowedSearchList,
		otherQuery,
		getTotalPlaylistsQuery,
	)

	getPlaylistsQuery := fmt.Sprintf(`
		SELECT *
		FROM playlist
		WHERE deleted IS NULL %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, searchQuery, sort_by, sort_type)

	rows, err := s.db.QueryContext(ctx, getPlaylistsQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	playlistList := make([]*types.Playlist, 0)
	for rows.Next() {
		i := new(types.Playlist)
		err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Slug,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Deleted,
			&i.CreatedBy,
			&i.UpdatedBy,
		)
		if err != nil {
			return nil, 0, err
		}
		playlistList = append(playlistList, i)
	}

	var i int
	if err := s.db.QueryRowContext(ctx, getTotalPlaylistsQuery).Scan(&i); err != nil {
		return nil, 0, err
	}

	return &playlistList, i, nil
}

func (s *Store) UpdatePlaylist(ctx context.Context, id string, updateData *types.PlaylistPayload) (*types.Playlist, error) {
	oldData, err := s.GetPlaylist(ctx, id)
	if err != nil {
		return nil, err
	}

	jwtPayload, err := utils.GetJWTPayload(ctx)
	if err != nil {
		return nil, err
	}

	query, args, update_ind := utils.GetSetQuery(
		id,
		"playlist",
		&[]*utils.SetQuery{
			utils.NewSetQuery(updateData.Title, "title"),
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

	newData, err := s.GetPlaylist(ctx, id)
	if err != nil {
		return nil, err
	}

	if *newData.Title != *oldData.Title {
		slug := utils.GenerateUniqueSlug(*newData.Title, func(slug string) bool {
			_, err := s.GetPlaylistBySlug(ctx, slug)
			return err == nil
		})

		query, args, _ := utils.GetSetQuery(
			id,
			"playlist",
			&[]*utils.SetQuery{
				utils.NewSetQuery(&slug, "slug"),
			},
		)

		_, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		return s.GetPlaylist(ctx, id)
	}

	return newData, nil
}

const deletePlaylistQuery = `
UPDATE playlist SET deleted = ?, updated_at = ?, updated_by = ?
WHERE id = ?
`

func (s *Store) DeletePlaylist(ctx context.Context, id string) error {
	playlist, err := s.GetPlaylist(ctx, id)
	if err != nil {
		return err
	}

	if playlist.Deleted == nil {
		jwtPayload, err := utils.GetJWTPayload(ctx)
		if err != nil {
			return err
		}
		_, err = s.db.ExecContext(ctx, deletePlaylistQuery, id, time.Now().UTC(), jwtPayload.UserId, id)
		return err
	}

	return nil
}

func (s *Store) CreateStartupPlaylist(ctx context.Context, playlistId string, startupId_list []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	jwtPayload, err := utils.GetJWTPayload(ctx)
	if err != nil {
		return err
	}

	for _, startupId := range startupId_list {
		var i *string
		err = tx.QueryRowContext(ctx, "SELECT deleted FROM startup_playlist WHERE playlist_id = ? AND startup_id = ?", playlistId, startupId).Scan(&i)

		if err == nil {
			if i != nil {
				_, err := tx.ExecContext(ctx, "UPDATE startup_playlist SET deleted = NULL, updated_at = ?, updated_by = ? WHERE playlist_id = ? AND startup_id = ?", time.Now().UTC(), jwtPayload.UserId, playlistId, startupId)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
			continue
		}

		query := `
			INSERT INTO startup_playlist (startup_id, playlist_id, created_at, updated_at, created_by,  updated_by) 
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err := tx.ExecContext(ctx, query, startupId, playlistId, time.Now().UTC(), time.Now().UTC(), jwtPayload.UserId, jwtPayload.UserId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (s *Store) RemoveStartupPlaylist(ctx context.Context, playlistId string, startupId_list []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	jwtPayload, err := utils.GetJWTPayload(ctx)
	if err != nil {
		return err
	}

	for _, startupId := range startupId_list {
		var i *string
		err = tx.QueryRowContext(ctx, "SELECT deleted FROM startup_playlist WHERE playlist_id = ? AND startup_id = ?", playlistId, startupId).Scan(&i)
		if err != nil || i != nil {
			continue
		}

		query := `
			UPDATE startup_playlist SET deleted = ?, updated_at = ?, updated_by = ? WHERE playlist_id = ? AND startup_id = ?
		`
		_, err := tx.ExecContext(ctx, query, playlistId, time.Now().UTC(), jwtPayload.UserId, playlistId, startupId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (s *Store) GetStartupsPlaylist(
	ctx context.Context,
	playlistId string,
	limit, offset int,
	sort_by, sort_type, search string,
	allowedSearchList *[]string,
	otherQuery *map[string]string,
) (*[]*types.Startup, int, error) {

	getTotalStartupsQuery := fmt.Sprintf(`
		SELECT COUNT(*) AS total
		FROM startup_playlist
		INNER JOIN startup ON startup_playlist.startup_id = startup.id
		INNER JOIN author ON startup.author_id = author.id
		WHERE startup.deleted is NULL AND startup_playlist.deleted is NULL AND startup_playlist.playlist_id = "%s"
	`, playlistId)

	searchQuery, getTotalStartupsQuery := utils.GetSearchQuery(
		validStartupSearchFields,
		search,
		allowedSearchList,
		otherQuery,
		getTotalStartupsQuery,
	)

	getStartupsQuery := fmt.Sprintf(`
		SELECT startup.id, title, slug, author_id, views, description, category, startup.image, pitch, startup.created_at, startup.updated_at, startup.created_by, startup.updated_by, startup.deleted, author.id, author.name, author.username, author.email, author.image
		FROM startup_playlist
		INNER JOIN startup ON startup_playlist.startup_id = startup.id
		INNER JOIN author ON startup.author_id = author.id
		WHERE startup.deleted IS NULL AND startup_playlist.deleted is NULL AND startup_playlist.playlist_id = "%s" %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, playlistId, searchQuery, sort_by, sort_type)

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
