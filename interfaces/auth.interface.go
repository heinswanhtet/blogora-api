package interfaces

import (
	"context"

	"github.com/heinswanhtet/blogora-api/types"
)

type AuthStore interface {
	CrateAuthor(ctx context.Context, author *types.Author) (*types.Author, error)

	GetAuthor(ctx context.Context, id string) (*types.Author, error)

	GetAuthorIdByEmail(ctx context.Context, email string) (string, error)
}
