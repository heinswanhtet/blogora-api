package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
	middleware "github.com/heinswanhtet/blogora-api/middlewares"
)

func (h *Handler) RegisterPlaylistRoutes(router *methods.CustomMux) {
	controller := controllers.NewPlaylistController(h.store)

	router.Attach("POST", "/playlist", controller.HandleCreatePlaylist, middleware.AuthenticateToken(h.store))
	router.Attach("GET", "/playlist", controller.HandleGetPlaylists)
	router.Attach("GET", "/playlist/{id}", controller.HandleGetSinglePlaylist)
	router.Attach("PUT", "/playlist/{id}", controller.HandleUpdatePlaylist, middleware.AuthenticateToken(h.store))
	router.Attach("DELETE", "/playlist/{id}", controller.HandleDeletePlaylist, middleware.AuthenticateToken(h.store))
}
