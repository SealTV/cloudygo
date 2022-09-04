package kvstore

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(addr string) error {
	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", s.keyValuePutHandler).Methods(http.MethodPut)
	r.HandleFunc("/v1/{key}", s.keyValueGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/{key}", s.keyValueDeleteHandler).Methods(http.MethodDelete)

	return http.ListenAndServeTLS(addr, "./certs/server.crt", "./certs/server.key", r)
}

func (s *Server) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, value)
}

func (s *Server) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.WritePut(key, string(value))

	if err := Put(key, string(value)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	logger.WriteDelete(key)

	if err := Delete(key); err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	InitializeTransactionLog()
	server := NewServer()
	log.Fatal(server.Run(":8080"))
}
