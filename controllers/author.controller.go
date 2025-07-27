package controllers

import (
	"fmt"
	"net/http"

	"github.com/heinswanhtet/blogora-api/services"
	"github.com/heinswanhtet/blogora-api/stores"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type AuthorController struct {
	authorService *services.AuthorService
}

func NewAuthorController(store *stores.Store) *AuthorController {
	return &AuthorController{
		authorService: services.NewAuthorService(store),
	}
}

func (c *AuthorController) HandleCreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author types.AuthorPayload

	if err := utils.ParseJSON(r, &author); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(author); err != nil {
		errors := utils.GetValidationErrors(err)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	result, status, err := c.authorService.CreateAuthor(r.Context(), &author)

	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "author created successfully", nil)
}

func (c *AuthorController) HandleGetAuthors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, pageSize := utils.GetPageAndPageSize(query)
	sort_by := utils.GetSanitizedQuery(
		query, "sort_by", "created_at",
		"name",
		"username",
		"email",
		"created_at",
		"updated_at",
	)
	sort_type := utils.GetSanitizedQuery(
		query, "sort_type", "desc",
		"asc",
		"desc",
	)
	search := query.Get("search")
	otherQuery := utils.GetRestOfQuery(query)

	result, total, status, err := c.authorService.GetAuthors(
		r.Context(),
		page,
		pageSize,
		sort_by,
		sort_type,
		search,
		&otherQuery,
	)

	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		result,
		"authors fetched successfully",
		utils.GenerateMetaPagination(page, pageSize, total),
	)
}

func (c *AuthorController) HandleGetSingleAuthor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	jwt, err := utils.CreateJWT(id)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(jwt)

	result, status, err := c.authorService.GetSingleAuthor(r.Context(), id)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "author fetched successfully", nil)
}

func (c *AuthorController) HandleUpdateAuthor(w http.ResponseWriter, r *http.Request) {
	jwtPayload, err := utils.GetJWTPayload(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var updateData types.AuthorPayload // won't update email

	if err := utils.ParseJSON(r, &updateData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, status, err := c.authorService.UpdateAuthor(r.Context(), jwtPayload.UserId, &updateData)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "author updated successfully", nil)
}

func (c *AuthorController) HandleDeleteAuthor(w http.ResponseWriter, r *http.Request) {
	jwtPayload, err := utils.GetJWTPayload(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	status, err := c.authorService.DeleteAuthor(r.Context(), jwtPayload.UserId)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, nil, "author deleted successfully", nil)
}
