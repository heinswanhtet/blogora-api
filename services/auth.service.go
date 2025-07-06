package services

import (
	"fmt"
	"net/http"

	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

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
