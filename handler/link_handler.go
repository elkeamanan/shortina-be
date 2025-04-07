package handler

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) handlerStoreLink(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req := &linkDomain.StoreLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.linkService.StoreLink(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to store link with err: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handlerGetLinkRedirection(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	messages, err := s.linkService.GetLinkRedirection(ctx, &linkDomain.GetLinkRedirectionRequest{
		Key: vars["key"],
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
