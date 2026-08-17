package httpg

import (
	"clodeTask/db"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler хранит пул соединений — создаётся один раз в main и
// передаётся во все методы автоматически (через получатель h).
type Handler struct {
	Pool *pgxpool.Pool
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var note db.Notes
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest) // 400
		return
	}

	if err := db.Create_notes(h.Pool, ctx, note); err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	array, err := db.CheclAll(h.Pool, ctx)
	if err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(array); err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}
}

func (h *Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest) // 400
		return
	}

	note, err := db.GetByID(h.Pool, ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "note not found", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) PutByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest) // 400
		return
	}

	var note db.Notes
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest) // 400
		return
	}
	note.Id = id

	// проверяем, что заметка вообще существует, прежде чем обновлять
	if _, err := db.GetByID(h.Pool, ctx, id); errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "note not found", http.StatusNotFound) // 404
		return
	} else if err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	if err := db.UpdateByID(h.Pool, ctx, note); err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	updated, err := db.GetByID(h.Pool, ctx, id)
	if err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *Handler) DeleteByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest) // 400
		return
	}

	// проверяем существование до удаления, чтобы вернуть корректный 404
	if _, err := db.GetByID(h.Pool, ctx, id); errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "note not found", http.StatusNotFound) // 404
		return
	} else if err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	if err := db.DeleteById(h.Pool, ctx, id); err != nil {
		http.Error(w, "error on server", http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204, без тела
}
