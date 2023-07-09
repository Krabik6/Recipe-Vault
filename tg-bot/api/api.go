package api

import "net/http"

type Api struct {
	Client http.Client
}

func NewApi() *Api {
	return &Api{
		Client: http.Client{},
	}
}
