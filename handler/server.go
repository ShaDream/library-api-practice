package handler

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

type Server struct {
	r *mux.Router
	s http.Server
}

func NewServer() *Server {
	router := GetRouter()
	addr := viper.GetString("api.address")
	return &Server{r: router,
		s: http.Server{
			Addr:    addr,
			Handler: router,
		}}
}

func (s *Server) Start() error {
	return s.s.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
