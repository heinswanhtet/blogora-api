package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heinswanhtet/blogora-api/interfaces"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type PlaylistService struct {
	store interfaces.PlaylistStore
}

func NewPlaylistService(store interfaces.PlaylistStore) *PlaylistService {
	return &PlaylistService{
		store: store,
	}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, data *types.PlaylistPayload) (*types.Playlist, int, error) {
	slug := utils.GenerateUniqueSlug(*data.Title, func(slug string) bool {
		_, err := s.store.GetPlaylistBySlug(ctx, slug)
		return err == nil
	})

	jwtPayload, err := utils.GetJWTPayload(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result, err := s.store.CratePlaylist(ctx, &types.Playlist{
		Title:     data.Title,
		Slug:      &slug,
		CreatedBy: &jwtPayload.UserId,
		UpdatedBy: &jwtPayload.UserId,
	})

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return result, http.StatusCreated, nil
}

func (s *PlaylistService) GetPlaylists(
	ctx context.Context,
	page, pageSize int,
	sort_by, sort_type, search string,
	otherQuery *map[string]string,
) (*[]*types.Playlist, int, int, error) {

	allowedSearchList := []string{"title", "slug"}

	playlists, total, err := s.store.GetPlaylists(
		ctx,
		pageSize,
		utils.GetOffsetToPaginate(page, pageSize),
		sort_by,
		sort_type,
		search,
		&allowedSearchList,
		otherQuery,
	)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	return playlists, total, http.StatusOK, nil
}

func (s *PlaylistService) GetSinglePlaylist(ctx context.Context, id string) (*types.Playlist, int, error) {
	playlist, err := s.store.GetPlaylist(ctx, id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return playlist, http.StatusOK, nil
}

func (s *PlaylistService) UpdatePlaylist(
	ctx context.Context,
	id string,
	updateData *types.PlaylistPayload,
) (*types.Playlist, int, error) {
	permission_ind, err := s.checkPermission(ctx, id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if !permission_ind {
		return nil, http.StatusForbidden, fmt.Errorf("permission denied")
	}

	playlist, err := s.store.UpdatePlaylist(ctx, id, updateData)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return playlist, http.StatusOK, nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id string) (int, error) {
	permission_ind, err := s.checkPermission(ctx, id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	if !permission_ind {
		return http.StatusForbidden, fmt.Errorf("permission denied")
	}

	err = s.store.DeletePlaylist(ctx, id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func (s *PlaylistService) checkPermission(ctx context.Context, id string) (bool, error) {
	jwtPayload, err := utils.GetJWTPayload(ctx)
	if err != nil {
		return false, err
	}

	startup, err := s.store.GetPlaylist(ctx, id)
	if err != nil {
		return false, err
	}

	if *startup.CreatedBy != jwtPayload.UserId {
		return false, nil
	}

	return true, nil
}
