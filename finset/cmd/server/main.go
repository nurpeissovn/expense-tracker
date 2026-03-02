package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/finset/app/internal/db"
	"github.com/finset/app/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

//go:embed static/index.html
var staticFiles embed.FS

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	h := handlers.New(pool)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health",               h.Health)
		r.Get("/stats",                h.GetStats)
		r.Get("/monthly-flow",         h.GetMonthlyFlow)
		r.Get("/category-breakdown",   h.GetCategoryBreakdown)
		r.Get("/transactions",         h.ListTransactions)
		r.Post("/transactions",        h.CreateTransaction)
		r.Get("/transactions/{id}",    h.GetTransaction)
		r.Delete("/transactions/{id}", h.DeleteTransaction)
		r.Post("/import",              h.ImportTransactions)
	})

	indexHTML, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		log.Fatalf("read embedded index.html: %v", err)
	}

	staticFS, _ := fs.Sub(staticFiles, "static")

	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path[1:]
		if path != "" {
			if _, statErr := fs.Stat(staticFS, path); statErr == nil {
				http.FileServer(http.FS(staticFS)).ServeHTTP(w, req)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(indexHTML)
	})

	log.Printf("FinSet listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
