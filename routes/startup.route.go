package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
)

func (h *Handler) RegisterStartupRoutes(router *methods.CustomMux) {
	controller := controllers.NewStartupController(h.store)

	router.Attach("POST", "/startup", controller.HandleCreateStartup)
	router.Attach("GET", "/startup", controller.HandleGetStartups)
	router.Attach("GET", "/startup/{id}", controller.HandleGetSingleStartup)
	router.Attach("PUT", "/startup/{id}", controller.HandleUpdateStartup)
	router.Attach("DELETE", "/startup/{id}", controller.HandleDeleteStartup)
}
