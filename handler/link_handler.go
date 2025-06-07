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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   fmt.Sprintf("failed to decode request with err: %s", err.Error()),
		})
		return
	}

	parsedActorUserId, _ := uuid.Parse(actorUserId)
	var createdBy *uuid.UUID
	if parsedActorUserId != uuid.Nil {
		createdBy = &parsedActorUserId
	}

	err = s.linkService.StoreLink(ctx, req, createdBy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to store link with err: %s", err.Error()),
		})
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

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   "key param is required",
		})
		return
	}

	redirection, err := s.linkService.GetLinkRedirection(ctx, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to get link redirection with err: %s", err.Error()),
		})
		return
	}

	if redirection == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 404,
			Message:   fmt.Sprintf("redirection is not found with identifier %s", key),
		})
		return
	}

	response := &linkDomain.GetLinkRedirectionResponse{
		Redirection: redirection,
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to marshal response with err: %s", err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *Server) handlerGetLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)

	if actorUserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   "invalid user id from token",
		})
		return
	}

	statusFilter := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 || pageSize == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   "wrong pagination request, cannot be zero",
		})
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to get links with err: %s", err.Error()),
		})
		return
	}

	response := &linkDomain.GetLinksResponse{
		Pagination: paginationResponse,
		Links:      links,
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to marshal response with err: %s", err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *Server) handlerGetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   "id param is required",
		})
		return
	}

	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)
	parsedActorUserId, _ := uuid.Parse(actorUserId)

	if parsedActorUserId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   "invalid user id from token",
		})
		return
	}

	link, err := s.linkService.GetLink(ctx, linkDomain.LinkPredicate{
		ID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to get link with err: %s", err.Error()),
		})
		return
	}

	if link == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 404,
			Message:   fmt.Sprintf("link with id %s is not found", id),
		})
		return
	}

	var linkCreatedBy string
	if link.CreatedBy != nil {
		linkCreatedBy = link.CreatedBy.String()
	}

	if linkCreatedBy != actorUserId {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 403,
			Message:   "cannot access link that are not created by current user",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to marshal response with err: %s", err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *Server) handlerUpdateLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	req := &linkDomain.UpdateLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   fmt.Sprintf("failed to decode request with err: %s", err.Error()),
		})
		return
	}

	actorUserId, _ := ctx.Value(XAuthorizedUserIDHeaderKey).(string)

	link, err := s.linkService.GetLink(ctx, linkDomain.LinkPredicate{
		ID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to get link with err: %s", err.Error()),
		})
		return
	}

	if link == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("link with id %s is not found", id),
		})
		return
	}

	var linkCreatedBy string
	if link.CreatedBy != nil {
		linkCreatedBy = link.CreatedBy.String()
	}

	if linkCreatedBy != actorUserId {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 403,
			Message:   "cannot access link that are not created by current user",
		})
		return
	}

	err = s.linkService.UpdateLink(ctx, &linkDomain.Link{
		Key:         req.Key,
		Redirection: req.Redirection,
		Status:      req.Status,
	}, linkDomain.LinkPredicate{ID: id})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to update link with err: %s", err.Error()),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
