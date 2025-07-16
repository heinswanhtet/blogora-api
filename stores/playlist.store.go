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
INSERT INTO playlist (id, title, slug, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
`

func (s *Store) CratePlaylist(ctx context.Context, playlist *types.Playlist) (*types.Playlist, error) {
	id := utils.GenerateUUID()

	_, err := s.db.ExecContext(ctx, createPlaylistQuery,
		id,
		playlist.Title,
		playlist.Slug,
		time.Now().UTC(),
		time.Now().UTC(),
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

	query, args, update_ind := utils.GetSetQuery(
		id,
		"playlist",
		&[]*utils.SetQuery{
			utils.NewSetQuery(updateData.Title, "title"),
		},
	)

	if !update_ind {
		return oldData, nil
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetPlaylist(ctx, id)
}

const deletePlaylistQuery = `
UPDATE playlist SET deleted = ?, updated_at = ?
WHERE id = ?
`

func (s *Store) DeletePlaylist(ctx context.Context, id string) error {
	playlist, err := s.GetPlaylist(ctx, id)
	if err != nil {
		return nil
	}

	if playlist.Deleted == nil {
		_, err = s.db.ExecContext(ctx, deletePlaylistQuery, id, time.Now().UTC(), id)
	}

	return err
}
