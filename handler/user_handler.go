package handler

import (
	"context"
	userDomain "elkeamanan/shortina/internal/user/domain"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &userDomain.RegisterUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   fmt.Sprintf("failed to decode request with err: %s", err.Error()),
		})
		return
	}

	err = s.userService.RegisterUser(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to register user with err: %s", err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success":true}`))
}

func (s *Server) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &userDomain.LoginUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   fmt.Sprintf("failed to decode request with err: %s", err.Error()),
		})
		return
	}

	tokenPair, err := s.userService.LoginUser(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to login user with err: %s", err.Error()),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(tokenPair)
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

func (s *Server) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actorUserId := ctx.Value(XAuthorizedUserIDHeaderKey).(string)

	vars := mux.Vars(r)
	id := vars["id"]

	if id != actorUserId {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 403,
			Message:   "cannot update user information of another user",
		})
		return
	}

	req := &userDomain.UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   fmt.Sprintf("failed to decode request with err: %s", err.Error()),
		})
		return
	}

	err = s.userService.UpdateUser(ctx, uuid.MustParse(id), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to update user with err: %s", err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

func (s *Server) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 400,
			Message:   "refresh token is required",
		})
		return
	}

	tokenPair, err := s.userService.RefreshUserToken(ctx, refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode: 500,
			Message:   fmt.Sprintf("failed to refresh user token with err: %s", err.Error()),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(tokenPair)
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
