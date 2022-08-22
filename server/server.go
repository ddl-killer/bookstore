package server

import (
	"encoding/json"
	"net/http"
	"time"

	"bookstore/server/middleware"

	"github.com/gorilla/mux"

	"bookstore/store"
)

type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		s: s,
		srv: &http.Server{
			Addr: addr,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/book", srv.createBookHandler).Methods(http.MethodPost)
	router.HandleFunc("/book/{id}", srv.updateBookHandler).Methods(http.MethodPost)
	router.HandleFunc("/book/{id}", srv.getBookHandler).Methods(http.MethodGet)
	router.HandleFunc("/book", srv.getAllBooksHandler).Methods(http.MethodGet)
	router.HandleFunc("/book/{id}", srv.delBookHandler).Methods(http.MethodDelete)

	srv.srv.Handler = middleware.Logging(middleware.Validating(router))
	return srv
}

func (bs *BookStoreServer) createBookHandler(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bs.s.Create(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

func (bs *BookStoreServer) updateBookHandler(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (bs *BookStoreServer) getBookHandler(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "no id found in request", http.StatusBadRequest)
		return
	}

	book, err := bs.s.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response(w, book)
}

func response(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (bs *BookStoreServer) getAllBooksHandler(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (bs *BookStoreServer) delBookHandler(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		bs.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}
