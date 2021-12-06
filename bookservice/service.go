package bookservices

import (
	"context"
	"errors"
	"sync"
)



type Service interface{
	PostBook(ctx context.Context, b Book) error
	GetBook(ctx context.Context, id string)(Book, error)
	PutBook(ctx context.Context, id string, b Book) error
	PatchBook(ctx context.Context, id string, b Book) error
}

type Book struct{
	ID string `json:"id"`
	Title string `json:"titolo,omitempty"`
	Authors string `json:"authors,omitempty"`
}

/*Puo essere inutile ai fini dell'esercizio
type Author struct{
	ID       string `json:"id"`
	FullName string `json:"fullname,omitempty"`
}*/


var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

//In sostanza: funzione per evitare che molteplici threads accedino e modifichno contemporaneamente la stessa risorsa
type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Book
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Book{},
	}
}


func (s *inmemService) PostBook(ctx context.Context, b Book) error {
	s.mtx.Lock()// chiude l'accesso alla risorsa
	defer s.mtx.Unlock()// finita l'esecuzione di QUESTA funzione riapre
	if _, ok := s.m[b.ID]; ok {
		return ErrAlreadyExists // POST = create, don't overwrite
	}
	s.m[b.ID] = b
	return nil
}

func (s *inmemService) GetBook(ctx context.Context, id string) (Book, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	b, ok := s.m[id]
	if !ok {
		return Book{}, ErrNotFound
	}
	return b, nil
}

func (s *inmemService) PutBook(ctx context.Context, id string, b Book) error {
	if id != b.ID {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = b // PUT = create or update
	return nil
}

func (s *inmemService) PatchBook(ctx context.Context, id string, b Book) error {
	if b.ID != "" && id != b.ID {
		return ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[id]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}
	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Profile definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if b.Title != "" {
		existing.Title = b.Title
	}
	if len(b.Authors) > 0 {
		existing.Authors = b.Authors
	}
	s.m[id] = existing
	return nil
}

