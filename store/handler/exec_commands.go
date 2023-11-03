package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	store "github.com/yyyoichi/Distributed-Task-Management-System/pkg/database"
)

var (
	ErrInvalidBodyProperty = errors.New("invalid body property")
	ErrSyntaxInvalidCmd    = errors.New("syntax error: invalid comand")
	ErrSyntaxInvalidArgs   = errors.New("syntax error: invalid args")
	ErrNoDataFound         = errors.New("sql error: no data found")
)

func Exec(commands []string, s *store.TStore) (string, error) {
	var resp string
	switch commands[0] {
	case "create":
		if len(commands[1:]) != 1 {
			err := fmt.Sprintf("%s: create cmd required 1 arg", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		id := s.Create(commands[1])
		resp = fmt.Sprintf("Created TODO[ID:%d]", id)
	case "list":
		list := []string{"TODO: "}
		cmps := []string{"COMPLETED TODO:"}
		for id, todo := range s.ByID {
			s := fmt.Sprintf("%d: %s", id, todo.Task)
			if todo.Completed {
				cmps = append(cmps, s)
			} else {
				list = append(list, s)
			}
		}
		resp = fmt.Sprintf("%s\n%s", strings.Join(list, "\n"), strings.Join(cmps, "\n"))
	case "update":
		if len(commands[1:]) != 2 {
			err := fmt.Sprintf("%s: update cmd required 2 args", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		id, err := strconv.Atoi(commands[1])
		if err != nil {
			err := fmt.Sprintf("%s: second argument must be a number", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		if commands[2] != "complete" && commands[2] != "open" {
			err := fmt.Sprintf("%s: third argument must be 'complete' or 'open'", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		if err := s.Update(id, commands[2] == "complete"); err != nil {
			return "", fmt.Errorf("%s: %s", ErrNoDataFound, err)
		}
		resp = fmt.Sprintf("Updated TODO[ID:%d]", id)
	case "delete":
		if len(commands[1:]) != 1 {
			err := fmt.Sprintf("%s: update cmd required 1 arg", ErrSyntaxInvalidArgs)
			return "", errors.New(err)
		}
		id, err := strconv.Atoi(commands[1])
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
		err := fmt.Errorf("%s: '%s'", ErrSyntaxInvalidCmd, commands[0])
		return "", err
	}
	return resp, nil
}
