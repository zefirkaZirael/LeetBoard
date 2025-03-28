package main

import (
	"1337bo4rd/internal/UI/handler"
	"1337bo4rd/internal/UI/middleware"
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/external"
	"1337bo4rd/internal/infrastructure/repository"
	"1337bo4rd/internal/infrastructure/s3"
	"1337bo4rd/internal/service"
	"1337bo4rd/utils"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	// Проверяем cmd
	utils.CheckFlags()

	// Подключаемся к БД
	db, err := repository.CheckDB()
	if err != nil {
		slog.Error("Failed to start program", "CheckDB err:", err)
		log.Fatal(err)
	}
	defer db.Close()
	slog.Info("Server successfully connected to DB")

	// Создаём репозитории
	postRepo := repository.NewPostRepository(db)
	sessionRepo := repository.NewSessionRepo(db)
	commentRepo := repository.NewCommentRepository(db)
	storageRepo := s3.NewS3Repo()
	externalRepo := external.DefaultExternalAPI()

	// Создаем хранилище в Triple-s
	err = storageRepo.InitBuckets()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err)
	}

	// Создаём сервисы (передаём ВСЕ аргументы)
	postServ := service.NewPostService(postRepo, sessionRepo, storageRepo, db)
	commentService := service.NewCommentService(commentRepo, postRepo, sessionRepo, storageRepo)
	sessionServ := service.NewUserSessionService(sessionRepo, externalRepo)
	// Создаём обработчики
	postHandler := handler.NewPostHandler(postServ, sessionServ, db)
	commentHandler := handler.NewCommentHandler(commentService, sessionServ, db) // ✅ Теперь правильно!
	sessionHandler := middleware.NewSessionHandler(sessionServ, db)

	// Создаём маршрутизатор
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			utils.ErrorPage(w, errors.New("method is undefined"), http.StatusMethodNotAllowed)
			return
		}
		utils.ErrorPage(w, errors.New("undefined URL address"), http.StatusBadRequest)
	})
	mux.HandleFunc("GET /catalog", sessionHandler.MainPage)
	mux.HandleFunc("POST /submit-post", postHandler.CreatePost)
	mux.HandleFunc("GET /catalog/post/{id}", postHandler.ServePost)
	mux.HandleFunc("GET /archive", postHandler.ServeArchivePage)
	mux.HandleFunc("GET /archive/post/{id}", postHandler.ServeArchivePost)
	mux.HandleFunc("POST /submit-comment", commentHandler.SubmitCommentHandler)
	mux.HandleFunc("GET /catalog/create-post", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/create-post.html")
	})

	// Запускаем сервер
	slog.Info(fmt.Sprintf("Server started at %d port", *domain.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *domain.Port), mux))
}
