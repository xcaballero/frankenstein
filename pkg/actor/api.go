package actor

import (
	"fmt"
	"net/http"
)

type API struct {
	action            chan func()
	pending           map[string]source
	committedSegments int64
	committedBytes    int64
	quitc             chan struct{}
}

func NewAPI() http.Handler {
	return &API{
		action:  make(chan func()),
		pending: make(map[string]source),
		quitc:   make(chan struct{}),
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handleCommit(w, r)
}

type source struct {
	segment Segment
	reading bool
}
type segment struct {
	size int64
}
type Segment interface {
	Size() int64
	Commit() error
}

func (s *segment) Size() int64 {
	return s.size
}
func (s *segment) Commit() error {
	return nil
}

func (a *API) handleCommit(w http.ResponseWriter, r *http.Request) {
	var (
		notFound  = make(chan struct{})
		notRead   = make(chan struct{})
		commitErr = make(chan error)
		commitOK  = make(chan int64)
	)

	a.action <- func() {
		id := r.URL.Query().Get("id")
		s, ok := a.pending[id]
		if !ok {
			close(notFound)
			return
		}
		if !s.reading {
			close(notRead)
			return
		}
		sz := s.segment.Size()
		if err := s.segment.Commit(); err != nil {
			commitErr <- err
			return
		}
		delete(a.pending, id)
		commitOK <- sz
	}

	select {
	case <-notFound:
		http.NotFound(w, r)
	case <-notRead:
		http.Error(w, "segment hasn't been read yet; can't commit", http.StatusBadRequest)
	case err := <-commitErr:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case n := <-commitOK:
		a.committedSegments++
		a.committedBytes += n
		fmt.Fprint(w, "Commited segment!\n")
	}
}
