package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
	middleware "github.com/heinswanhtet/blogora-api/middlewares"
)

func (h *Handler) RegisterAuthorRoutes(router *methods.CustomMux) {
	controller := controllers.NewAuthorController(h.store)

	router.Attach("POST", "/author", controller.HandleCreateAuthor)
	router.Attach("GET", "/author", controller.HandleGetAuthors)
	router.Attach("GET", "/author/{id}", controller.HandleGetSingleAuthor)
	router.Attach("PUT", "/author", controller.HandleUpdateAuthor, middleware.AuthenticateToken(h.store))
	router.Attach("DELETE", "/author", controller.HandleDeleteAuthor, middleware.AuthenticateToken(h.store))
}
