package routes

import store "github.com/heinswanhtet/blogora-api/stores"

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{
		store: store,
	}
}
