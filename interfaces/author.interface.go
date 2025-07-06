package interfaces

import (
	"context"

	"github.com/heinswanhtet/blogora-api/types"
)

type AuthorStore interface {
	CrateAuthor(ctx context.Context, author *types.Author) (*types.Author, error)

	GetAuthor(ctx context.Context, id string) (*types.Author, error)

	GetAuthorIdByEmail(ctx context.Context, email string) (string, error)

	GetAuthors(
		ctx context.Context,
		limit, offset int,
		sort_by, sort_type, search string,
		allowedSearchList *[]string,
		otherQuery *map[string]string,
	) (*[]*types.Author, int, error)

	UpdateAuthor(ctx context.Context, id string, updateData *types.AuthorPayload) (*types.Author, error)

	DeleteAuthor(ctx context.Context, id string) error
}
