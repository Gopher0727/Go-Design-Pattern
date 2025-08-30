package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Message   string `json:"message,omitempty"`
}

func RegisterHandlers(mux *http.ServeMux, repo UserRepository) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Name == "" || req.Email == "" {
			http.Error(w, "Bad Request: name and email are required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		u := &User{
			Name:      req.Name,
			Email:     req.Email,
			CreatedAt: time.Now(),
		}
		id, err := repo.Create(ctx, u)
		if err != nil {
			// 唯一键冲突则返回已存在的用户并提示
			var me *mysql.MySQLError
			if errors.As(err, &me) && me.Number == 1062 {
				if exist, e := repo.GetByEmail(ctx, req.Email); e == nil && exist != nil {
					w.Header().Set("Location", "/users/"+strconv.FormatInt(exist.ID, 10))
					writeJSON(w, http.StatusOK, userResponse{
						ID:        exist.ID,
						Name:      exist.Name,
						Email:     exist.Email,
						CreatedAt: exist.CreatedAt.Format(time.RFC3339),
						Message:   "User already exists!",
					})
					return
				}
			}
			http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		created, err := repo.GetByID(ctx, id)
		if err != nil || created == nil {
			http.Error(w, "created but cannot load", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", "/users/"+strconv.FormatInt(id, 10))
		writeJSON(w, http.StatusCreated, userResponse{
			ID:        created.ID,
			Name:      created.Name,
			Email:     created.Email,
			CreatedAt: created.CreatedAt.Format(time.RFC3339),
			Message:   "Created!",
		})
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
