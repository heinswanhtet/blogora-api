package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
	middleware "github.com/heinswanhtet/blogora-api/middlewares"
)

func (h *Handler) RegisterStartupRoutes(router *methods.CustomMux) {
	controller := controllers.NewStartupController(h.store)

	router.Attach("POST", "/startup", controller.HandleCreateStartup, middleware.AuthenticateToken(h.store))
	router.Attach("GET", "/startup", controller.HandleGetStartups)
	router.Attach("GET", "/startup/{id}", controller.HandleGetSingleStartup)
	router.Attach("PUT", "/startup/{id}", controller.HandleUpdateStartup, middleware.AuthenticateToken(h.store))
	router.Attach("DELETE", "/startup/{id}", controller.HandleDeleteStartup, middleware.AuthenticateToken(h.store))
}
