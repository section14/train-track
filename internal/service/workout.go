package service

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"time"

	//"strconv"
	//"time"

	"github.com/section14/train-track/internal/model"
)

type WorkoutManager interface {
	GetWorkouts() []model.Workout
	GetWorkout(id int) []model.Movement
	GetLastWorkout() []model.Movement
	UpdateMovement(m model.Movement) ([]model.Movement, error)
	AddWorkout() error
	UpdateWorkout(e model.Workout) ([]model.Movement, error)
	DeleteWorkout(id int) error
}

type WorkoutService struct {
	service WorkoutManager
}

func NewWorkoutService(service WorkoutManager) *WorkoutService {
	return &WorkoutService{service: service}
}

func (ws *WorkoutService) GetAll() []model.Workout {
	workouts := ws.service.GetWorkouts()
	return workouts
}

func (ws *WorkoutService) GetLast() []model.Movement {
	last := ws.service.GetLastWorkout()
	return last
}

func (ws *WorkoutService) Add() error {
	err := ws.service.AddWorkout()
	return err
}

func (ws *WorkoutService) Update(r *http.Request) ([]model.Movement, error) {
    var m = model.Movement{}

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&m)
    if err != nil {
        return []model.Movement{}, err
    }

    m.Date = time.Now()
	movements, err := ws.service.UpdateMovement(m)
	return movements, err
}
