package routes

import (
	"github.com/heinswanhtet/blogora-api/methods"
	"github.com/heinswanhtet/blogora-api/services"
)

func (h *Handler) RegisterAuthRoutes(router *methods.CustomMux) {
	authRouter := methods.NewCustomMux()
	service := services.NewService(h.store)

	authRouter.Attach("POST", "/register", service.HandleRegister)
	authRouter.Attach("POST", "/login", service.HandleLogin)

	router.Use("/auth/", authRouter)
}
