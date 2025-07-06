package services

import store "github.com/heinswanhtet/blogora-api/stores"

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}
