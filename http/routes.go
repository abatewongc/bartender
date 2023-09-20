package http

import "github.com/go-chi/chi"

func (svc *service) Routes() chi.Router {
	r := chi.NewRouter()

	// todo: search (by champ name, skin name), sort, filter
	r.Route("/skins", func(r chi.Router) {
		r.Get("/", svc.handleGetSkins())
		r.Get("/{skindId}", svc.handleGetSkin())
	})

	r.Route("/blacklist", func(r chi.Router) {
		r.Put("/{skinId}", svc.handleBlacklistPut())
		r.Delete("/{skinId}", svc.handleBlacklistDelete())
	})

	return r
}
