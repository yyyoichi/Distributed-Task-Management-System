package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"yyyoichi/Distributed-Task-Management-System/store/store"

	"github.com/go-playground/validator"
)

type StoreHandlers struct {
	s store.TStore
}

func NewStoreHandler() StoreHandlers {
	return StoreHandlers{store.NewStore()}
}

func (sh *StoreHandlers) commandsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: commandsHandler")
	var data struct {
		Task string `json:"task" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Printf("[KV] Get cmds '%s'\n", data.Task)
	cmds := strings.Split(data.Task, " ")
	resp, err := sh.s.Read(cmds)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}
	log.Println("[KV] Response 200")
	w.Write(bytes.NewBufferString(resp).Bytes())
	w.WriteHeader(http.StatusOK)
}
