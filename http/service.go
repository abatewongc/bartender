package http

import "net/http"

type service struct{}

func (svc *service) handleGetSkins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
