package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator"
)

var (
	ErrInvalidBodyProperty = errors.New("invalid body property")
	ErrSyntaxInvalidCmd    = errors.New("syntax error: invalid comand")
	ErrSyntaxInvalidArgs   = errors.New("syntax error: invalid args")
	ErrNoDataFound         = errors.New("sql error: no data found")
)

func NewStore() TStore {
	return TStore{ByID: make(map[int]*struct {
		task      string
		completed bool
	})}
}

type TStore struct {
	mu   sync.Mutex
	ByID map[int]*struct {
		task      string
		completed bool
	}
}

type ReqJson struct {
	Task string `json:"task" validate:"required"`
}

func (s *TStore) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Body)
	var data ReqJson
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
	resp, err := s.read(cmds)
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

func (s *TStore) read(cmds []string) (string, error) {
	var resp string
	switch cmds[0] {
	case "create":
		if len(cmds[1:]) != 1 {
			err := fmt.Sprintf("%s: create cmd required 1 arg", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		id := s.Create(cmds[1])
		resp = fmt.Sprintf("Created TODO[ID:%d]", id)
	case "list":
		list := []string{"TODO: "}
		cmps := []string{"COMPLETED TODO:"}
		for id, todo := range s.ByID {
			s := fmt.Sprintf("%d: %s", id, todo.task)
			if todo.completed {
				cmps = append(cmps, s)
			} else {
				list = append(list, s)
			}
		}
		resp = fmt.Sprintf("%s\n%s", strings.Join(list, "\n"), strings.Join(cmps, "\n"))
	case "update":
		if len(cmds[1:]) != 2 {
			err := fmt.Sprintf("%s: update cmd required 2 args", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		id, err := strconv.Atoi(cmds[1])
		if err != nil {
			err := fmt.Sprintf("%s: second argument must be a number", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		if cmds[2] != "complete" && cmds[2] != "open" {
			err := fmt.Sprintf("%s: third argument must be 'complete' or 'open'", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		if err := s.Update(id, cmds[2] == "complete"); err != nil {
			return "", fmt.Errorf("%s: %s", ErrNoDataFound, err)
		}
		resp = fmt.Sprintf("Updated TODO[ID:%d]", id)
	case "delete":
		if len(cmds[1:]) != 1 {
			err := fmt.Sprintf("%s: update cmd required 1 arg", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		id, err := strconv.Atoi(cmds[1])
		if err != nil {
			err := fmt.Sprintf("%s: second argument must be a number", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		if err := s.Delete(id); err != nil {
			err := fmt.Errorf("%s: %s", ErrNoDataFound, err)
			return "", err
		}
		resp = fmt.Sprintf("Deleted TODO[ID:%d]", id)
	case "help":
		h := []string{
			"create <task>: Create a new todo with the specified task.",
			"list: List all todos, separated into incomplete and completed todos.",
			"update <id> <status>: Update the status of the todo with the specified id. Status can be 'complete' or 'open'.",
			"delete <id>: Delete the todo with the specified id.",
		}
		resp = strings.Join(h, "\n")
	default:
		err := fmt.Errorf("%s: '%s'", ErrSyntaxInvalidCmd, cmds[0])
		return "", err
	}
	return resp, nil
}

func (s *TStore) Create(task string) int {
	s.mu.Lock()
	id := s.nextID()
	s.ByID[id] = &struct {
		task      string
		completed bool
	}{task, false}
	s.mu.Unlock()
	return id
}

func (s *TStore) Update(id int, completed bool) error {
	todo, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	s.mu.Lock()
	todo.completed = completed
	s.mu.Unlock()
	return nil
}

func (s *TStore) Delete(id int) error {
	_, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	s.mu.Lock()
	delete(s.ByID, id)
	s.mu.Unlock()
	return nil
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
