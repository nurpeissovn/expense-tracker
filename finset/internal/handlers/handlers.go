package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/finset/app/internal/db"
	"github.com/finset/app/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler holds shared dependencies.
type Handler struct {
	DB *db.Pool
}

func New(d *db.Pool) *Handler {
	return &Handler{DB: d}
}

// ─── helpers ─────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseBody(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ─── GET /api/transactions ───────────────────────────────

func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	txs, err := h.DB.ListTransactions(r.Context())
	if err != nil {
		log.Printf("ListTransactions: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch transactions")
		return
	}
	writeJSON(w, http.StatusOK, txs)
}

// ─── GET /api/transactions/{id} ──────────────────────────

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tx, err := h.DB.GetTransaction(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "transaction not found")
		return
	}
	writeJSON(w, http.StatusOK, tx)
}

// ─── POST /api/transactions ──────────────────────────────

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTransactionRequest
	if err := parseBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Trim whitespace
	req.Type     = strings.TrimSpace(req.Type)
	req.Category = strings.TrimSpace(req.Category)
	req.Method   = strings.TrimSpace(req.Method)
	req.Note     = strings.TrimSpace(req.Note)
	req.Date     = strings.TrimSpace(req.Date)

	if req.Method == "" {
		req.Method = "Cash"
	}
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	if errMsg := req.Validate(); errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	id := uuid.New().String()
	tx, err := h.DB.CreateTransaction(r.Context(), id, req)
	if err != nil {
		log.Printf("CreateTransaction: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}
	writeJSON(w, http.StatusCreated, tx)
}

// ─── DELETE /api/transactions/{id} ───────────────────────

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	found, err := h.DB.DeleteTransaction(r.Context(), id)
	if err != nil {
		log.Printf("DeleteTransaction: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to delete transaction")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "transaction not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─── GET /api/stats ──────────────────────────────────────

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.DB.GetStats(r.Context())
	if err != nil {
		log.Printf("GetStats: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// ─── GET /api/monthly-flow?months=7 ──────────────────────

func (h *Handler) GetMonthlyFlow(w http.ResponseWriter, r *http.Request) {
	months := 7
	if m := r.URL.Query().Get("months"); m != "" {
		if n := 0; json.Unmarshal([]byte(m), &n) == nil && n > 0 && n <= 24 {
			months = n
		}
	}
	rows, err := h.DB.GetMonthlyFlow(r.Context(), months)
	if err != nil {
		log.Printf("GetMonthlyFlow: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch monthly flow")
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

// ─── GET /api/category-breakdown ─────────────────────────

func (h *Handler) GetCategoryBreakdown(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.GetCategoryBreakdown(r.Context())
	if err != nil {
		log.Printf("GetCategoryBreakdown: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch category breakdown")
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

// ─── POST /api/import ────────────────────────────────────
// Accepts JSON array of transactions or {transactions:[…]}.
// Uses ON CONFLICT DO NOTHING — existing records are never overwritten.

func (h *Handler) ImportTransactions(w http.ResponseWriter, r *http.Request) {
	type importBody struct {
		Transactions []models.Transaction `json:"transactions"`
	}

	var body importBody
	if err := parseBody(r, &body); err != nil {
		// Maybe it's a raw array
		writeError(w, http.StatusBadRequest, "invalid JSON: expected {transactions:[…]}")
		return
	}

	if len(body.Transactions) == 0 {
		writeError(w, http.StatusBadRequest, "no transactions provided")
		return
	}

	// Assign new IDs to any records missing one
	for i := range body.Transactions {
		if body.Transactions[i].ID == "" {
			body.Transactions[i].ID = uuid.New().String()
		}
		if body.Transactions[i].CreatedAt.IsZero() {
			body.Transactions[i].CreatedAt = time.Now()
		}
	}

	inserted, err := h.DB.BulkInsert(r.Context(), body.Transactions)
	if err != nil {
		log.Printf("BulkInsert: %v", err)
		writeError(w, http.StatusInternalServerError, "import failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"inserted": inserted,
		"skipped":  len(body.Transactions) - inserted,
		"total":    len(body.Transactions),
	})
}

// ─── GET /api/health ─────────────────────────────────────

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unreachable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
