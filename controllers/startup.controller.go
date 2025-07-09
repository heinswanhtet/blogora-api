package controllers

import (
	"net/http"

	"github.com/heinswanhtet/blogora-api/services"
	"github.com/heinswanhtet/blogora-api/stores"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type StartupController struct {
	startupService *services.StartupService
}

func NewStartupController(store *stores.Store) *StartupController {
	return &StartupController{
		startupService: services.NewStartupService(store),
	}
}

func (c *StartupController) HandleCreateStartup(w http.ResponseWriter, r *http.Request) {
	var startup types.StartupPayload

	if err := utils.ParseJSON(r, &startup); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(startup); err != nil {
		errors := utils.GetValidationErrors(err)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	result, status, err := c.startupService.CreateStartup(r.Context(), &startup)

	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "startup created successfully", nil)
}

func (c *StartupController) HandleGetStartups(w http.ResponseWriter, r *http.Request) {
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

	result, total, status, err := c.startupService.GetStartups(
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
		"startups fetched successfully",
		utils.GenerateMetaPagination(page, pageSize, total),
	)
}

func (c *StartupController) HandleGetSingleStartup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, status, err := c.startupService.GetSingleStartup(r.Context(), id)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "startup fetched successfully", nil)
}

func (c *StartupController) HandleUpdateStartup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updateData types.StartupPayload

	if err := utils.ParseJSON(r, &updateData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, status, err := c.startupService.UpdateStartup(r.Context(), id, &updateData)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, result, "startup updated successfully", nil)
}

func (c *StartupController) HandleDeleteStartup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	status, err := c.startupService.DeleteStartup(r.Context(), id)
	if err != nil {
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, status, nil, "startup deleted successfully", nil)
}
