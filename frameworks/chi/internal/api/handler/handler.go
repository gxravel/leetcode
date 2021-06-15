package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gxravel/leetcode/frameworks/chi/internal/model/book"
)

type Permission string

func (p Permission) IsAdmin() bool {
	return p == "my-secret"
}

type Environment struct {
	Books book.Wrapper
}

func (e *Environment) AllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := e.Books.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (e *Environment) Book(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	book, ok := ctx.Value("book").(*book.Book)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		perm, ok := ctx.Value("acl.permission").(Permission)
		if !ok || !perm.IsAdmin() {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(map[string]string{"msg": "here's my fellow"})
	})
	return r
}

func (e *Environment) BookCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		book, err := e.Books.ByID(bookID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "book", book)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
