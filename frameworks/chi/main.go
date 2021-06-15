package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gxravel/leetcode/frameworks/chi/internal/api/handler"
	"github.com/gxravel/leetcode/frameworks/chi/internal/db"
	"github.com/gxravel/leetcode/frameworks/chi/internal/model/book"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connString := "user:pass@tcp(localhost:3306)/default"
	myDB, err := db.Init("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}
	env := handler.Environment{Books: book.New(myDB)}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/book", func(r chi.Router) {
		r.Get("/", env.AllBooks)

		r.Route("/book/{bookID}", func(r chi.Router) {
			r.Use(env.BookCtx)
			r.Get("/", env.Book)
		})
	})

	r.Mount("/admin", handler.AdminRouter())

	srv := http.Server{
		Addr:    ":8090",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
