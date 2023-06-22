package main

import (
	"forum/internal/app/database"
	serviceHandlers "forum/internal/app/service/handlers"
	serviceRepo "forum/internal/app/service/repository"
	serviceUsecase "forum/internal/app/service/usecase"
	userHandlers "forum/internal/app/user/handlers"
	userRepo "forum/internal/app/user/repository"
	userUsecase "forum/internal/app/user/usecase"
	router2 "github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

func main() {

	postgres, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := userRepo.NewRepo(postgres.GetPostgres())
	if err := userRepository.Prepare(); err != nil {
		log.Fatalln(err)
	}

	serviceRepository := serviceRepo.NewRepo(postgres.GetPostgres())
	if err := serviceRepository.Prepare(); err != nil {
		log.Fatalln(err)
	}

	userUseCase := userUsecase.NewUseCase(*userRepository)
	serviceUseCase := serviceUsecase.NewUseCase(*serviceRepository)

	userHandler := userHandlers.NewHandler(*userUseCase)
	serviceHandler := serviceHandlers.NewHandler(*serviceUseCase)

	router := router2.New()

	router.POST("/api/user/{nickname}/create", userHandler.CreateUser)

	router.GET("/api/user/{nickname}/profile", userHandler.GetUserInfo)

	router.POST("/api/user/{nickname}/profile", userHandler.ChangeUser)

	router.POST("/api/service/clear", serviceHandler.ClearDB)

	router.GET("/api/service/status", serviceHandler.Status)

	if err := fasthttp.ListenAndServe(":5000", router.Handler); err != nil {
		log.Fatal(err)
	}
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"url":    r.URL,
			"method": r.Method,
			"body":   r.Body,
		}).Info()
		next.ServeHTTP(w, r)
	})
}
