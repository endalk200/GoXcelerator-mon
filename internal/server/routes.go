package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/endalk200/GoXcelerator/internal/database"
	"github.com/endalk200/GoXcelerator/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var validate *validator.Validate

func (s *Server) RegisterRoutes() http.Handler {
	validate = validator.New()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Post("/auth/signup", s.SignupHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Message string `json:"message"`
	}

	utils.Response(w, http.StatusOK, Response{
		Message: "Hello world",
	})
}

type CreateUserRequestSchema struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
}

func (s *Server) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequestSchema

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(body); err != nil {
		// Extract validation errors
		validationErrors := err.(validator.ValidationErrors)
		utils.ResponseError(w, http.StatusBadRequest, validationErrors.Error())
		return
	}

	ctx := r.Context()
	user, err := s.db.AddUser(ctx, database.AddUserParams{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	})
	if err != nil {
		log.Printf("%v", err)
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create a new user")
		return
	}

	utils.Response(w, http.StatusCreated, user)
}
