package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
	middleware "github.com/heinswanhtet/blogora-api/middlewares"
)

func (h *Handler) RegisterAuthorRoutes(router *methods.CustomMux) {
	controller := controllers.NewAuthorController(h.store)

	router.Attach("POST", "/author", controller.HandleCreateAuthor)
	router.Attach("GET", "/author", controller.HandleGetAuthors, middleware.AuthenticateToken(h.store))
	router.Attach("GET", "/author/{id}", controller.HandleGetSingleAuthor)
	router.Attach("PUT", "/author/{id}", controller.HandleUpdateAuthor)
	router.Attach("DELETE", "/author/{id}", controller.HandleDeleteAuthor)
}
