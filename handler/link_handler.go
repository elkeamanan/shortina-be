package handler

import (
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) handlerStoreLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)
	req := &linkDomain.StoreLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedActorUserId, _ := uuid.Parse(actorUserId)
	var createdBy *uuid.UUID
	if parsedActorUserId != uuid.Nil {
		createdBy = &parsedActorUserId
	}

	err = s.linkService.StoreLink(ctx, req, createdBy)
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
	redirection, err := s.linkService.GetLinkRedirection(ctx, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if redirection == "" {
		http.Error(w, fmt.Sprintf("redirection is not found with identifier %s", key), http.StatusNotFound)
		return
	}

	response := &linkDomain.GetLinkRedirectionResponse{
		Redirection: redirection,
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (s *Server) handlerGetLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)
	statusFilter := r.URL.Query().Get("status")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 || pageSize == 0 {
		http.Error(w, "wrong pagination request, cannot be zero", http.StatusBadRequest)
		return
	}

	if actorUserId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	links, paginationResponse, err := s.linkService.GetPaginatedLinks(
		ctx,
		linkDomain.LinkPredicate{UserID: actorUserId, Status: linkDomain.LinkStatus(statusFilter)},
		util.PaginationParam{
			PageSize:    uint32(page),
			CurrentPage: uint32(pageSize),
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &linkDomain.GetLinksResponse{
		Pagination: paginationResponse,
		Links:      links,
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (s *Server) handlerGetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)
	parsedActorUserId, _ := uuid.Parse(actorUserId)

	if parsedActorUserId == uuid.Nil {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	link, err := s.linkService.GetLink(ctx, linkDomain.LinkPredicate{
		ID: id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if link == nil {
		http.Error(w, "link is not found", http.StatusNotFound)
		return
	}

	var linkCreatedBy string
	if link.CreatedBy != nil {
		linkCreatedBy = link.CreatedBy.String()
	}

	if linkCreatedBy != actorUserId {
		http.Error(w, "cannot access link that are not created by current user", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (s *Server) handlerUpdateLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	req := &linkDomain.UpdateLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)

	link, err := s.linkService.GetLink(ctx, linkDomain.LinkPredicate{
		ID: id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if link == nil {
		http.Error(w, "link is not found", http.StatusNotFound)
		return
	}

	var linkCreatedBy string
	if link.CreatedBy != nil {
		linkCreatedBy = link.CreatedBy.String()
	}

	if linkCreatedBy != actorUserId {
		http.Error(w, "cannot access link that are not created by current user", http.StatusUnauthorized)
		return
	}

	err = s.linkService.UpdateLink(ctx, &linkDomain.Link{
		Key:         req.Key,
		Redirection: req.Redirection,
		Status:      req.Status,
	}, linkDomain.LinkPredicate{ID: id})

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update link with err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
