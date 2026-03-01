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

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Load .env in development (ignored if file doesn't exist — Railway injects env vars directly)
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ── Database ──────────────────────────────────────────
	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}
	defer pool.Close()

	// Run safe migration (never drops or alters existing data)
	if err := pool.Migrate(); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// ── Router ────────────────────────────────────────────
	h := handlers.New(pool)
	r := chi.NewRouter()

	// Middleware
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

	// ── API routes ────────────────────────────────────────
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", h.Health)
		r.Get("/stats", h.GetStats)
		r.Get("/monthly-flow", h.GetMonthlyFlow)
		r.Get("/category-breakdown", h.GetCategoryBreakdown)

		r.Get("/transactions", h.ListTransactions)
		r.Post("/transactions", h.CreateTransaction)
		r.Get("/transactions/{id}", h.GetTransaction)
		r.Delete("/transactions/{id}", h.DeleteTransaction)

		r.Post("/import", h.ImportTransactions)
	})

	// ── Static frontend (SPA) ─────────────────────────────
	// Serve everything in the embedded static/ folder.
	// Any route that isn't /api/* falls through to index.html.
	staticFS, err := fs.Sub(staticFiles, "../static")
	if err != nil {
		log.Fatalf("embed static: %v", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		// Try to serve the exact file first; fall back to index.html for SPA routing
		_, statErr := fs.Stat(staticFS, req.URL.Path[1:])
		if statErr != nil {
			// Serve index.html for client-side navigation
			f, err := staticFS.Open("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			defer f.Close()
			stat, _ := f.Stat()
			http.ServeContent(w, req, "index.html", stat.ModTime(), f.(interface {
				Read([]byte) (int, error)
				Seek(int64, int) (int64, error)
			}))
			return
		}
		fileServer.ServeHTTP(w, req)
	})

	// ── Start ─────────────────────────────────────────────
	log.Printf("🚀 FinSet server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
