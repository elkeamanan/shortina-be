package handler

import (
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) handlerStoreLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &linkDomain.StoreLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.linkService.StoreLink(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to store link with err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handlerGetLinkRedirection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	key := vars["key"]
	messages, err := s.linkService.GetLinkRedirection(ctx, &linkDomain.GetLinkRedirectionRequest{
		Key: key,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if messages == nil {
		http.Error(w, fmt.Sprintf("message is not found with identifier %s", key), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
