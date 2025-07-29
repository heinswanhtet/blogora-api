package controllers

import (
	"net/http"

	"github.com/heinswanhtet/blogora-api/services"
	"github.com/heinswanhtet/blogora-api/stores"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(store *stores.Store) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(store),
	}
}

func (c *AuthController) HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.SsoPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := utils.GetValidationErrors(err)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	user, status, err := c.authService.OAuthLogin(r.Context(), payload)

	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	jwt, err := utils.CreateJWT(*user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	response := struct {
		types.Author
		Token string `json:"token"`
	}{
		Author: *user,
		Token:  jwt,
	}

	utils.WriteJSON(w, status, response, "user logged in successfully", nil)
}
