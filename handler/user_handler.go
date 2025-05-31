package handler

import (
	"context"
	userDomain "elkeamanan/shortina/internal/user/domain"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &userDomain.RegisterUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.userService.RegisterUser(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenPair, err := s.userService.LoginUser(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(tokenPair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "cannot update user information of another user", http.StatusUnauthorized)
		return
	}

	req := &userDomain.UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.userService.UpdateUser(ctx, uuid.MustParse(id), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

func (s *Server) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	tokenPair, err := s.userService.RefreshUserToken(ctx, refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(tokenPair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
