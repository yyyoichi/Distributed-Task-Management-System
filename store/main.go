package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	ErrInvalidBodyProperty = errors.New("invalid body property")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("port is not found")
	}
	addr := fmt.Sprintf(":%s", port)

	var store = NewStore()

	http.HandleFunc("/", store.handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
	log.Println("start database")
}

func NewStore() TStore { return TStore{ByID: make(map[int]string)} }

type TStore struct {
	mu   sync.Mutex
	ByID map[int]string
}

type ReqJson struct {
	Task string `validate:"required"`
}

func (s *TStore) handler(w http.ResponseWriter, r *http.Request) {
	var data *ReqJson
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Save(data.Task)
	w.WriteHeader(http.StatusOK)
}

func (s *TStore) Save(task string) {
	s.mu.Lock()
	id := s.nextID()
	s.ByID[id] = task
	s.mu.Unlock()
}

func (s *TStore) nextID() int {
	max := 0
	for id := range s.ByID {
		if max < id {
			max = id
		}
	}
	return max + 1
}
