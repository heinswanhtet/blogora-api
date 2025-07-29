package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/heinswanhtet/blogora-api/interfaces"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type AuthService struct {
	store interfaces.AuthStore
}

func NewAuthService(store interfaces.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *Service) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := utils.GetValidationErrors(err)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	_, err := s.store.GetUserIdByEmail(r.Context(), user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("user with email:%s already exists", user.Email))
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	result, err := s.store.CrateUser(r.Context(), &types.User{
		Email:     user.Email,
		Password:  hashedPassword,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, result, "User registered successfully", nil)
}

func (s *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var userPayload types.LoginUserPayload

	if err := utils.ParseJSON(r, &userPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(userPayload); err != nil {
		errors := utils.GetValidationErrors(err)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	user, err := s.store.GetLoginUser(r.Context(), userPayload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "email or password incorrect")
		return
	}

	if isPwdValid := utils.ComparePassword(userPayload.Password, user.Password); !isPwdValid {
		utils.WriteError(w, http.StatusBadRequest, "email or password incorrect")
		return
	}

	utils.WriteJSON(w, http.StatusOK, user, "User logged in successfully", nil)
}

func (s *AuthService) OAuthLogin(ctx context.Context, data types.SsoPayload) (*types.Author, int, error) {
	url := "https://www.googleapis.com/oauth2/v3/userinfo"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not create request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *data.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error making http request: %s", err)
	}
	defer res.Body.Close()

	var userData types.GoogleOAuthResp
	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not response body: %s", err)
	}

	if userData.Error != "" {
		return nil, http.StatusBadRequest, fmt.Errorf("%s", userData.Error)
	}

	id, err := s.store.GetAuthorIdByEmail(ctx, userData.Email)
	if err != nil {
		username := utils.FormatUserName(userData.Name)
		user, err := s.store.CrateAuthor(ctx, &types.Author{
			Name:     &userData.Name,
			Username: &username,
			Email:    &userData.Email,
			Image:    &userData.Picture,
		})
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return user, http.StatusOK, nil
	}

	user, err := s.store.GetAuthor(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}
