package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heinswanhtet/blogora-api/interfaces"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type StartupService struct {
	store         interfaces.StartupStore
	authorService *AuthorService
}

func NewStartupService(store interfaces.StartupStore) *StartupService {
	return &StartupService{
		store:         store,
		authorService: NewAuthorService(store.(interfaces.AuthorStore)),
	}
}

func (s *StartupService) CreateStartup(ctx context.Context, data *types.StartupPayload) (*types.Startup, int, error) {
	_, status, err := s.authorService.GetSingleAuthor(ctx, *data.AuthorId)
	if err != nil {
		return nil, status, fmt.Errorf("author_id '%s' not found", *data.AuthorId)
	}

	slug := utils.GenerateUniqueSlug(*data.Title, func(slug string) bool {
		_, err := s.store.GetStartupBySlug(ctx, slug)
		return err == nil
	})

	result, err := s.store.CrateStartup(ctx, &types.Startup{
		Title:       data.Title,
		Slug:        &slug,
		AuthorId:    data.AuthorId,
		Description: data.Description,
		Category:    data.Category,
		Image:       data.Image,
		Pitch:       data.Pitch,
	})

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return result, http.StatusCreated, nil
}

func (s *StartupService) GetStartups(
	ctx context.Context,
	page, pageSize int,
	sort_by, sort_type, search string,
	otherQuery *map[string]string,
) (*[]*types.Startup, int, int, error) {

	allowedSearchList := []string{"title", "author.name"}

	startups, total, err := s.store.GetStartups(
		ctx,
		pageSize,
		utils.GetOffsetToPaginate(page, pageSize),
		sort_by,
		sort_type,
		search,
		&allowedSearchList,
		otherQuery,
	)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	return startups, total, http.StatusOK, nil
}

func (s *StartupService) GetSingleStartup(ctx context.Context, id string) (*types.Startup, int, error) {
	startup, err := s.store.GetStartup(ctx, id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	author, err := s.authorService.store.GetAuthor(ctx, *startup.AuthorId)
	if err != nil || author == nil {
		return nil, http.StatusNotFound, err
	}
	startup.Author = *author

	return startup, http.StatusOK, nil
}

func (s *StartupService) UpdateStartup(
	ctx context.Context,
	id string,
	updateData *types.StartupPayload,
) (*types.Startup, int, error) {

	startup, err := s.store.UpdateStartup(ctx, id, updateData)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return startup, http.StatusOK, nil
}

func (s *StartupService) DeleteStartup(ctx context.Context, id string) (int, error) {
	err := s.store.DeleteStartup(ctx, id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
