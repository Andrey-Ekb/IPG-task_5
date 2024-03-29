package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"main_test.go/internal/repository/map_repository"
)

func main() {
	log.Println("start program")

	// создаем репозиторий для хранения данных (map)
	repository := map_repository.New()

	// создаем сервис для обработки запросов клиента
	friendService := friend_service.New(repository)

	// создаем новый роутер
	r := chi.NewRouter()
	// показывает информацию по запросам
	r.Use(middleware.Logger)

	// создать друга
	r.Post("/create", friendService.CreateFriend)
	// подружить двух друзей
	r.Post("/make_friends", friendService.MakeFriend)
	// показать всех друзей пользователя
	r.Get("/friends/{user_id}", friendService.GetFriends)
	// удалить пользователя
	r.Delete("/{user}", friendService.DeleteFriend)
	// Обновить возраст пользователя
	r.Put("/{user_id}", friendService.UpdateAge)

	// запускаем работу сервиса на 8080 порту
	http.ListenAndServe(":8080", r)
}
