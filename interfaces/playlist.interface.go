package interfaces

import (
	"context"

	"github.com/heinswanhtet/blogora-api/types"
)

type PlaylistStore interface {
	CratePlaylist(ctx context.Context, author *types.Playlist) (*types.Playlist, error)

	GetPlaylist(ctx context.Context, id string) (*types.Playlist, error)

	GetPlaylistBySlug(ctx context.Context, slug string) (*types.Playlist, error)

	GetPlaylists(
		ctx context.Context,
		limit, offset int,
		sort_by, sort_type, search string,
		allowedSearchList *[]string,
		otherQuery *map[string]string,
	) (*[]*types.Playlist, int, error)

	UpdatePlaylist(ctx context.Context, id string, updateData *types.PlaylistPayload) (*types.Playlist, error)

	DeletePlaylist(ctx context.Context, id string) error
}
