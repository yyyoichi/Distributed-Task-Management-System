package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"
)

type StoreHandlers struct {
	dc *document.TDocument
}

func NewStoreHandler() StoreHandlers {
	return StoreHandlers{document.NewTDocument()}
}

func (sh *StoreHandlers) CommandsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: CommandsHandler")
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
	resp, err := Exec(cmds, sh.dc)
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

func (sh *StoreHandlers) DifferencesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: DifferencesHandler")
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

	resp := []document.TodoDataset{}
	for id, todo := range sh.dc.GetLatestVersionTodo(data.Version) {
		resp = append(resp, document.TodoDataset{
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

func (sh *StoreHandlers) SynchronizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: SynchronizeHandler")
	var data struct {
		// 同期バージョン
		Version int                    `json:"version" validate:"required"`
		Todos   []document.TodoDataset `json:"todos" validate:"required"`
	}
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}
	sh.dc.Synchronize(r.Context(), data.Version, data.Todos)

	log.Println("[KV] Response 200")
	w.WriteHeader(http.StatusOK)
}
