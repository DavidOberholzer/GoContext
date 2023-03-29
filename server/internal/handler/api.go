package handler

import "errors"

type apiCreateRequest struct {
	Name string `json:"name"`
}

func (a *apiCreateRequest) Validate() error {
	if a.Name == "" {
		return errors.New("'name' is required")
	}
	return nil
}
