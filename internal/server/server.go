package server

import (
	"merchshop/internal/api"
	"merchshop/internal/middleware"
	"merchshop/internal/repository"
	"merchshop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	addr         string
	merchService *service.MerchService
}

func NewServer(addr string, repo repository.MerchRepository) *Server {
	return &Server{
		addr:         addr,
		merchService: service.NewMerchService(repo),
	}
}

func (s *Server) ListenAndServe() error {
	apiServer := api.NewAPIServer(s.merchService)

	r := gin.Default()
	r.Use(api.JSONErrorHandler)
	r.Use(middleware.JWTMiddleware())

	handler := api.NewStrictHandler(apiServer, nil)
	api.RegisterHandlers(r, handler)

	httpServer := &http.Server{
		Handler: r,
		Addr:    s.addr,
	}

	return httpServer.ListenAndServe()
}
