package api

import (
	"net/http"
)

type Router struct {
	mux     *http.ServeMux
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{
		mux:     http.NewServeMux(),
		handler: handler,
	}
}

func (router *Router) Setup() *http.ServeMux {
	router.mux.HandleFunc("GET /questions/", router.handler.GetQuestions)
	router.mux.HandleFunc("POST /questions/", router.handler.CreateQuestion)
	router.mux.HandleFunc("GET /questions/{id}", router.handler.GetQuestion)
	router.mux.HandleFunc("DELETE /questions/{id}", router.handler.DeleteQuestion)

	router.mux.HandleFunc("POST /questions/{id}/answers/", router.handler.CreateAnswer)
	router.mux.HandleFunc("GET /answers/{id}", router.handler.GetAnswer)
	router.mux.HandleFunc("DELETE /answers/{id}", router.handler.DeleteAnswer)

	return router.mux
}
