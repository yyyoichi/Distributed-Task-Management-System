package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrInvalidBodyProperty = errors.New("invalid body property")
	ErrSyntaxInvalidCmd    = errors.New("syntax error: invalid comand")
	ErrSyntaxInvalidArgs   = errors.New("syntax error: invalid args")
	ErrNoDataFound         = errors.New("sql error: no data found")
)

func NewStore() TStore {
	return TStore{
		mu:          sync.Mutex{},
		ByID:        make(map[int]*todo),
		nextVersion: 1,
	}
}

type TStore struct {
	mu          sync.Mutex
	nextVersion int
	ByID        map[int]*todo
}

type todo struct {
	task      string
	completed bool
	version   int
	deleted   bool
}

func (s *TStore) Read(cmds []string) (string, error) {
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
	s.ByID[id] = &todo{
		task:      task,
		completed: false,
		deleted:   false,
		version:   s.nextVersion,
	}
	s.nextVersion++
	s.mu.Unlock()
	return id
}

func (s *TStore) Update(id int, completed bool) error {
	s.mu.Lock()
	todo, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	todo.completed = completed
	todo.version = s.nextVersion
	s.nextVersion++
	s.mu.Unlock()
	return nil
}

func (s *TStore) Delete(id int) error {
	s.mu.Lock()
	todo, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	todo.deleted = true
	todo.version = s.nextVersion
	s.nextVersion++
	s.mu.Unlock()
	return nil
}

// [version]よりも大きいバージョンを持つTODOを返す
func (s *TStore) GetLatestVersionTodo(version int) map[int]*todo {
	todos := map[int]*todo{}
	s.mu.Lock()
	for id, todo := range s.ByID {
		if version < todo.version {
			todos[id] = todo
		}
	}
	s.mu.Unlock()
	return todos
}

// 更新バージョンを同期する
func (s *TStore) SyncNextVersion(newVersion int) {
	s.nextVersion = newVersion
}

// データを上書きする。
func (s *TStore) SyncTodoAt(id int, newTodo todo) error {
	s.mu.Lock()
	_, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	s.ByID[id] = &newTodo
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
