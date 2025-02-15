package server

import (
	"merchshop/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) ListenAndServe() error {
	apiServer := api.NewServer()

	r := gin.Default()
	r.Use(api.JSONErrorHandler)

	handler := api.NewStrictHandler(apiServer, nil)
	api.RegisterHandlers(r, handler)

	httpServer := &http.Server{
		Handler: r,
		Addr:    s.addr,
	}

	return httpServer.ListenAndServe()
}
