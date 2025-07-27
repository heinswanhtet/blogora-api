package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heinswanhtet/blogora-api/interfaces"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

type AuthorService struct {
	store interfaces.AuthorStore
}

func NewAuthorService(store interfaces.AuthorStore) *AuthorService {
	return &AuthorService{
		store: store,
	}
}

func (s *AuthorService) CreateAuthor(ctx context.Context, data *types.AuthorPayload) (*types.Author, int, error) {
	_, err := s.store.GetAuthorIdByEmail(ctx, *data.Email)
	if err == nil {
		return nil, http.StatusForbidden, fmt.Errorf("author with email:%s already exists", *data.Email)
	}

	result, err := s.store.CrateAuthor(ctx, &types.Author{
		// Name:     *author.Name,
		// Username: *author.Username,
		// Email:    *author.Email,
		// Image:    utils.Maybe(author.Image).Else(""),
		// Bio:      utils.Maybe(author.Bio).Else(""),
		Name:     data.Name,
		Username: data.Username,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
	})

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return result, http.StatusCreated, nil
}

func (s *AuthorService) GetAuthors(
	ctx context.Context,
	page, pageSize int,
	sort_by, sort_type, search string,
	otherQuery *map[string]string,
) (*[]*types.Author, int, int, error) {

	allowedSearchList := []string{"name", "username", "email"}

	authors, total, err := s.store.GetAuthors(
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

	return authors, total, http.StatusOK, nil
}

func (s *AuthorService) GetSingleAuthor(ctx context.Context, id string) (*types.Author, int, error) {
	author, err := s.store.GetAuthor(ctx, id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return author, http.StatusOK, nil
}

func (s *AuthorService) UpdateAuthor(
	ctx context.Context,
	id string,
	updateData *types.AuthorPayload,
) (*types.Author, int, error) {

	author, err := s.store.UpdateAuthor(ctx, id, updateData)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return author, http.StatusOK, nil
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id string) (int, error) {
	err := s.store.DeleteAuthor(ctx, id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
