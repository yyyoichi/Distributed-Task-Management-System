package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

func parseBody(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		return err
	}
	return nil
}
