package interfaces

import (
	"context"

	"github.com/heinswanhtet/blogora-api/types"
)

type StartupStore interface {
	CrateStartup(ctx context.Context, author *types.Startup) (*types.Startup, error)

	GetStartup(ctx context.Context, id string) (*types.Startup, error)

	GetStartupBySlug(ctx context.Context, slug string) (*types.Startup, error)

	GetStartups(
		ctx context.Context,
		limit, offset int,
		sort_by, sort_type, search string,
		allowedSearchList *[]string,
		otherQuery *map[string]string,
	) (*[]*types.Startup, int, error)

	UpdateStartup(ctx context.Context, id string, updateData *types.StartupPayload) (*types.Startup, error)

	DeleteStartup(ctx context.Context, id string) error
}
