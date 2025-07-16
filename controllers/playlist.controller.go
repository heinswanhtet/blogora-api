package controllers

import (
	"net/http"

	"github.com/heinswanhtet/blogora-api/services"
	"github.com/heinswanhtet/blogora-api/stores"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type PlaylistController struct {
	playlistService *services.PlaylistService
}

func NewPlaylistController(store *stores.Store) *PlaylistController {
	return &PlaylistController{
		playlistService: services.NewPlaylistService(store),
	}
}

func (c *PlaylistController) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var playlist types.PlaylistPayload

	if err := utils.ParseJSON(r, &playlist); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(playlist); err != nil {
		errors := utils.GetValidationErrors(err)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	result, status, err := c.playlistService.CreatePlaylist(r.Context(), &playlist)

	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "playlist created successfully", nil)
}

func (c *PlaylistController) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, pageSize := utils.GetPageAndPageSize(query)
	sort_by := utils.GetSanitizedQuery(
		query, "sort_by", "created_at",
		"title",
		"slug",
		"created_at",
		"updated_at",
	)
	sort_type := utils.GetSanitizedQuery(
		query, "sort_type", "desc",
		"asc",
		"desc",
	)
	search := query.Get("search")
	otherQuery := utils.GetRestOfQuery(query)

	result, total, status, err := c.playlistService.GetPlaylists(
		r.Context(),
		page,
		pageSize,
		sort_by,
		sort_type,
		search,
		&otherQuery,
	)

	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		result,
		"playlists fetched successfully",
		utils.GenerateMetaPagination(page, pageSize, total),
	)
}

func (c *PlaylistController) HandleGetSinglePlaylist(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, status, err := c.playlistService.GetSinglePlaylist(r.Context(), id)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "playlist fetched successfully", nil)
}

func (c *PlaylistController) HandleUpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updateData types.PlaylistPayload

	if err := utils.ParseJSON(r, &updateData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, status, err := c.playlistService.UpdatePlaylist(r.Context(), id, &updateData)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "playlist updated successfully", nil)
}

func (c *PlaylistController) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	status, err := c.playlistService.DeletePlaylist(r.Context(), id)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, nil, "playlist deleted successfully", nil)
}
