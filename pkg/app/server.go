package app

import (
	"crabigateur-api/pkg/api"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router *gin.Engine
	userService api.UserService
	cardService api.CardService
}

func NewServer(router *gin.Engine, userService api.UserService, cardService api.CardService) *Server{
	return &Server{
		router: router,
		userService: userService,
		cardService: cardService,
	}
}

func (s *Server) Run() error {
	// run function that initializes the routes
	r := s.Routes()
	s.RegisterValidators()

	// run the server through the router
	err := r.Run()
	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}

func (s *Server) RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("sortable", api.ValidSortOrders)
		v.RegisterStructValidation(api.CardStructValidation, api.Card{})
	}
}