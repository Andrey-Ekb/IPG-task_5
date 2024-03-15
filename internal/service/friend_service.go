package friend_service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/login/test_project/internal/entity"
	"github.com/login/test_project/internal/repository"
)

type friendService struct {
	repo repository.Repository
}

// конструктор для создания сервиса
func New(repo repository.Repository) *friendService {
	return &friendService{
		repo: repo,
	}
}

// Создать пользователя
func (s *friendService) CreateFriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	friend := &entity.Friend{}
	if err := json.NewDecoder(r.Body).Decode(&friend); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(friend.Name) == 0 || friend.Age < 1 {
		writeError(w, "incorrect values", http.StatusBadRequest)
		return
	}

	id, err := s.repo.CreateFriend(*friend)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body := entity.ResponseCreate{
		Status: http.StatusOK,
		Id:     id,
	}
	resp, err := jsoniter.Marshal(body)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

// Подружить двух пользователей
func (s *friendService) MakeFriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	makeFriend := &entity.MakeFriend{}
	if err := json.NewDecoder(r.Body).Decode(&makeFriend); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(makeFriend.SourceId); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(makeFriend.TargetId); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	nameSource, nameTarget, err := s.repo.MakeFriend(makeFriend.SourceId, makeFriend.TargetId)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := fmt.Sprintf("%s и %s теперь друзья", nameSource, nameTarget)
	body := entity.ResponseMakeFriend{
		Status: http.StatusOK,
		Result: result,
	}
	resp, err := jsoniter.Marshal(body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(resp)
}

// Удаляем пользователя
func (s *friendService) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "user")
	if _, err := uuid.Parse(id); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	name, err := s.repo.DeleteFriend(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	body := entity.ResponseMakeFriend{
		Status: http.StatusOK,
		Result: name,
	}
	resp, err := jsoniter.Marshal(body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(resp)
}

// Показать всех друзей пользователя
func (s *friendService) GetFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "user_id")
	if _, err := uuid.Parse(id); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	friends, err := s.repo.GetFriends(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := jsoniter.Marshal(friends)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(resp)
}

// Обновить возраст
func (s *friendService) UpdateAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newAge := &entity.NewAge{}
	if err := json.NewDecoder(r.Body).Decode(&newAge); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := chi.URLParam(r, "user_id")
	if _, err := uuid.Parse(id); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := s.repo.UpdateAge(id, newAge.Age)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
