package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/store"
)

type StoreHandlers struct {
	s *store.TStore
}

func NewStoreHandler() StoreHandlers {
	return StoreHandlers{store.NewStore()}
}

func (sh *StoreHandlers) commandsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: commandsHandler")
	var data struct {
		Task string `json:"task" validate:"required"`
	}
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}

	log.Printf("[KV] Get cmds '%s'\n", data.Task)
	cmds := strings.Split(data.Task, " ")
	resp, err := Exec(cmds, sh.s)
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

func (sh *StoreHandlers) differencesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: differencesHandler")
	var data struct {
		// 同期バージョン
		Version int `json:"version" validate:"required"`
	}
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}

	log.Printf("[KV] Get version '%d'\n", data.Version)

	resp := []store.TodoDateset{}
	for id, todo := range sh.s.GetLatestVersionTodo(data.Version) {
		resp = append(resp, store.TodoDateset{
			ID:        id,
			Task:      todo.Task,
			Completed: todo.Completed,
			Deleted:   todo.Deleted,
			Version:   todo.Version,
		})
	}

	personJSON, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("[KV] Response 200")
	w.Write(personJSON)
	w.WriteHeader(http.StatusOK)
}

func (sh *StoreHandlers) syncHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: syncHandler")
	var data struct {
		// 同期バージョン
		Version int                 `json:"version" validate:"required"`
		Todos   []store.TodoDateset `json:"todos" validate:"required"`
	}
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}
	sh.s.Sync(r.Context(), data.Version, data.Todos)

	log.Println("[KV] Response 200")
	w.WriteHeader(http.StatusOK)
}
