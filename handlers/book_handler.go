package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"biblioteca-api/models"
)

// BookStore abstrai o acesso a dados — assim conseguimos mockar nos testes.
type BookStore interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (*models.Book, error)
	Create(b *models.Book) error
	Update(id int, b *models.Book) (bool, error)
	Delete(id int) (bool, error)
}

type BookHandler struct {
	store BookStore
}

func NewBookHandler(store BookStore) *BookHandler {
	return &BookHandler{store: store}
}

// BooksHandler trata GET /books e POST /books
func (h *BookHandler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// BookHandler trata GET /books/{id}, PUT /books/{id} e DELETE /books/{id}
func (h *BookHandler) BookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"id inválido"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) list(w http.ResponseWriter, _ *http.Request) {
	books, err := h.store.GetAll()
	if err != nil {
		http.Error(w, `{"error":"erro ao buscar livros"}`, http.StatusInternalServerError)
		return
	}
	if books == nil {
		books = []models.Book{}
	}
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) create(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, `{"error":"corpo inválido"}`, http.StatusBadRequest)
		return
	}
	if err := h.store.Create(&b); err != nil {
		http.Error(w, `{"error":"erro ao criar livro"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func (h *BookHandler) get(w http.ResponseWriter, _ *http.Request, id int) {
	b, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, `{"error":"erro ao buscar livro"}`, http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, `{"error":"livro não encontrado"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(b)
}

func (h *BookHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, `{"error":"corpo inválido"}`, http.StatusBadRequest)
		return
	}
	found, err := h.store.Update(id, &b)
	if err != nil {
		http.Error(w, `{"error":"erro ao atualizar livro"}`, http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, `{"error":"livro não encontrado"}`, http.StatusNotFound)
		return
	}
	b.ID = id
	json.NewEncoder(w).Encode(b)
}

func (h *BookHandler) delete(w http.ResponseWriter, _ *http.Request, id int) {
	found, err := h.store.Delete(id)
	if err != nil {
		http.Error(w, `{"error":"erro ao deletar livro"}`, http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, `{"error":"livro não encontrado"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
