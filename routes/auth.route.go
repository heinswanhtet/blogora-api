package routes

import (
	"github.com/heinswanhtet/blogora-api/controllers"
	"github.com/heinswanhtet/blogora-api/methods"
)

// func (h *Handler) RegisterAuthRoutes(router *methods.CustomMux) {
// 	authRouter := methods.NewCustomMux()
// 	service := services.NewService(h.store)

// 	authRouter.Attach("POST", "/register", service.HandleRegister)
// 	authRouter.Attach("POST", "/login", service.HandleLogin)
// 	authRouter.Attach("POST", "/sso", service.HandleOAuthLogin)

// 	router.Use("/auth/", authRouter)
// }

func (h *Handler) RegisterAuthRoutes(router *methods.CustomMux) {
	authRouter := methods.NewCustomMux()
	controller := controllers.NewAuthController(h.store)

	authRouter.Attach("POST", "/sso", controller.HandleOAuthLogin)

	router.Use("/auth/", authRouter)
}
