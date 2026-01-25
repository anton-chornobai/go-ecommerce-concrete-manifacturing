package handlers

import "net/http"

type ProdctsHandlerInterface interface {
	GetProducts() http.Handler
}

type ProdctsHandler struct {
}

func (r *ProdctsHandler) GetProducts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
