package main

import (
	"forum/internal/app/database"
	serviceHandlers "forum/internal/app/service/handlers"
	serviceRepo "forum/internal/app/service/repository"
	serviceUsecase "forum/internal/app/service/usecase"

	userHandlers "forum/internal/app/user/handlers"
	userRepo "forum/internal/app/user/repository"
	userUsecase "forum/internal/app/user/usecase"

	threadHandlers "forum/internal/app/thread/handlers"
	threadRepo "forum/internal/app/thread/repository"
	threadUCase "forum/internal/app/thread/usecase"

	forumHandlers "forum/internal/app/forum/handlers"
	forumRepo "forum/internal/app/forum/repository"
	forumUsecase "forum/internal/app/forum/usecase"

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

	threadRepository := threadRepo.NewRepo(postgres.GetPostgres())
	if err := threadRepository.Prepare(); err != nil {
		log.Fatalln(err)
	}

	forumRepository := forumRepo.NewRepo(postgres.GetPostgres())
	if err := forumRepository.Prepare(); err != nil {
		log.Fatalln(err)
	}

	userUseCase := userUsecase.NewUseCase(*userRepository)
	serviceUseCase := serviceUsecase.NewUseCase(*serviceRepository)
	threadUseCase := threadUCase.NewUseCase(*threadRepository)
	forumUseCase := forumUsecase.NewUseCase(*forumRepository, *userRepository, *threadRepository)

	forumHandler := forumHandlers.NewHandler(*forumUseCase)
	userHandler := userHandlers.NewHandler(*userUseCase)
	serviceHandler := serviceHandlers.NewHandler(*serviceUseCase)
	threadHandler := threadHandlers.NewHandler(*threadUseCase)

	router := router2.New()

	router.POST("/api/user/{nickname}/create", userHandler.CreateUser)

	router.GET("/api/user/{nickname}/profile", userHandler.GetUserInfo)

	router.POST("/api/user/{nickname}/profile", userHandler.ChangeUser)

	router.POST("/api/service/clear", serviceHandler.ClearDB)

	router.GET("/api/service/status", serviceHandler.Status)

	router.GET("/api/thread/{slug_or_id}/details", threadHandler.ThreadInfo)

	router.POST("/api/thread/{slug_or_id}/details", threadHandler.ChangeThread)

	router.POST("/api/thread/{slug_or_id}/vote", threadHandler.VoteThread)

	router.POST("/api/forum/create", forumHandler.Create)

	//done
	router.GET("/api/forum/{slug}/details", forumHandler.Details)

	//done
	router.POST("/api/forum/{slug}/create", forumHandler.CreateThread)

	//done
	router.GET("/api/forum/{slug}/users", forumHandler.GetUsers)

	//done
	router.GET("/api/forum/{slug}/threads", forumHandler.GetThreads)

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
