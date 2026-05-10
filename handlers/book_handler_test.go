package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"biblioteca-api/models"
)

// ─── Mock do BookStore ───────────────────────────────────────────────────────

type mockStore struct {
	books   map[int]models.Book
	nextID  int
	failOps bool
}

func newMockStore() *mockStore {
	return &mockStore{
		books: map[int]models.Book{
			1: {ID: 1, Title: "Livro 1", Author: "Autor 1", Year: 2020, Available: true},
			2: {ID: 2, Title: "Livro 2", Author: "Autor 2", Year: 2021, Available: false},
		},
		nextID: 3,
	}
}

func (m *mockStore) GetAll() ([]models.Book, error) {
	if m.failOps {
		return nil, errors.New("erro simulado")
	}
	out := make([]models.Book, 0, len(m.books))
	for _, b := range m.books {
		out = append(out, b)
	}
	return out, nil
}

func (m *mockStore) GetByID(id int) (*models.Book, error) {
	if m.failOps {
		return nil, errors.New("erro simulado")
	}
	b, ok := m.books[id]
	if !ok {
		return nil, nil
	}
	return &b, nil
}

func (m *mockStore) Create(b *models.Book) error {
	if m.failOps {
		return errors.New("erro simulado")
	}
	b.ID = m.nextID
	m.nextID++
	m.books[b.ID] = *b
	return nil
}

func (m *mockStore) Update(id int, b *models.Book) (bool, error) {
	if m.failOps {
		return false, errors.New("erro simulado")
	}
	if _, ok := m.books[id]; !ok {
		return false, nil
	}
	b.ID = id
	m.books[id] = *b
	return true, nil
}

func (m *mockStore) Delete(id int) (bool, error) {
	if m.failOps {
		return false, errors.New("erro simulado")
	}
	if _, ok := m.books[id]; !ok {
		return false, nil
	}
	delete(m.books, id)
	return true, nil
}

// ─── Testes ──────────────────────────────────────────────────────────────────

// 1) Listar livros com sucesso
func TestListBooks_OK(t *testing.T) {
	h := NewBookHandler(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	h.BooksHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado %d, recebi %d", http.StatusOK, rec.Code)
	}

	var got []models.Book
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("esperado 2 livros, recebi %d", len(got))
	}
}

// 2) Buscar livro por ID que existe
func TestGetBook_Found(t *testing.T) {
	h := NewBookHandler(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	rec := httptest.NewRecorder()
	h.BookHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado %d, recebi %d", http.StatusOK, rec.Code)
	}

	var got models.Book
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if got.ID != 1 || got.Title != "Livro 1" {
		t.Fatalf("livro errado: %+v", got)
	}
}

// 3) Buscar livro que não existe → 404
func TestGetBook_NotFound(t *testing.T) {
	h := NewBookHandler(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	rec := httptest.NewRecorder()
	h.BookHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperado %d, recebi %d", http.StatusNotFound, rec.Code)
	}
}

// 4) Criar livro novo
func TestCreateBook_OK(t *testing.T) {
	store := newMockStore()
	h := NewBookHandler(store)

	body := `{"title":"Novo","author":"Autor X","year":2024,"available":true}`
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.BooksHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperado %d, recebi %d", http.StatusCreated, rec.Code)
	}

	var got models.Book
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if got.ID == 0 || got.Title != "Novo" {
		t.Fatalf("livro criado errado: %+v", got)
	}
	if len(store.books) != 3 {
		t.Fatalf("esperado 3 livros no store, tem %d", len(store.books))
	}
}

// 5) Criar livro com JSON inválido → 400
func TestCreateBook_InvalidJSON(t *testing.T) {
	h := NewBookHandler(newMockStore())

	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader("{ json quebrado"))
	rec := httptest.NewRecorder()
	h.BooksHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado %d, recebi %d", http.StatusBadRequest, rec.Code)
	}
}

// 6) Atualizar livro existente
func TestUpdateBook_OK(t *testing.T) {
	store := newMockStore()
	h := NewBookHandler(store)

	body := `{"title":"Atualizado","author":"Novo Autor","year":2025,"available":false}`
	req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.BookHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado %d, recebi %d", http.StatusOK, rec.Code)
	}
	if store.books[1].Title != "Atualizado" {
		t.Fatalf("livro não foi atualizado: %+v", store.books[1])
	}
}

// 7) Deletar livro
func TestDeleteBook_OK(t *testing.T) {
	store := newMockStore()
	h := NewBookHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	rec := httptest.NewRecorder()
	h.BookHandler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("esperado %d, recebi %d", http.StatusNoContent, rec.Code)
	}
	if _, ok := store.books[1]; ok {
		t.Fatalf("livro 1 ainda existe no store após delete")
	}
}

// 8) ID inválido na URL → 400
func TestBookHandler_InvalidID(t *testing.T) {
	h := NewBookHandler(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	rec := httptest.NewRecorder()
	h.BookHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado %d, recebi %d", http.StatusBadRequest, rec.Code)
	}
}
