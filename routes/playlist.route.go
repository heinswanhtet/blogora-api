package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
)

func (h *Handler) RegisterPlaylistRoutes(router *methods.CustomMux) {
	controller := controllers.NewPlaylistController(h.store)

	router.Attach("POST", "/playlist", controller.HandleCreatePlaylist)
	router.Attach("GET", "/playlist", controller.HandleGetPlaylists)
	router.Attach("GET", "/playlist/{id}", controller.HandleGetSinglePlaylist)
	router.Attach("PUT", "/playlist/{id}", controller.HandleUpdatePlaylist)
	router.Attach("DELETE", "/playlist/{id}", controller.HandleDeletePlaylist)
}
