package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"yyyoichi/Distributed-Task-Management-System/store/store"
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
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
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

type todoDataset struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
	Deleted   bool   `json:"deleted"`
	Version   int    `json:"version"`
}

func (sh *StoreHandlers) differencesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[KV] Call: differencesHandler")
	var data struct {
		Version int `json:"version" validate:"required"`
	}
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}

	log.Printf("[KV] Get version '%d'\n", data.Version)

	resp := []todoDataset{}
	for id, todo := range sh.s.GetLatestVersionTodo(data.Version) {
		resp = append(resp, todoDataset{
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
		// 最大更新バージョン+1
		NextVersion int           `json:"nextVersion" validate:"required"`
		Todos       []todoDataset `json:"todos" validate:"required"`
	}
	if err := parseBody(r, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes.NewBufferString(err.Error()).Bytes())
		return
	}
	var wg sync.WaitGroup

	syncTodo := func(todos []todoDataset) {
		defer wg.Done()
		var todoWg sync.WaitGroup
		todoWg.Add(len(todos))
		for _, todo := range todos {
			go func(t todoDataset) {
				defer todoWg.Done()
				err := sh.s.SyncTodoAt(t.ID, store.Todo{
					Task:      t.Task,
					Completed: t.Completed,
					Deleted:   t.Deleted,
					Version:   t.Version,
				})
				if err != nil {
					log.Println(err.Error())
				}
			}(todo)
		}
		todoWg.Wait()
	}

	syncVersion := func(nextVersion int) {
		defer wg.Done()
		sh.s.SyncNextVersion(nextVersion)
	}

	wg.Add(2)
	go syncTodo(data.Todos)
	go syncVersion(data.NextVersion)
	wg.Wait()
	log.Println("[KV] Response 200")
	w.WriteHeader(http.StatusOK)
}
