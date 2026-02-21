package service

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"strconv"
	"time"

	//"strconv"
	//"time"

	"github.com/section14/train-track/internal/model"
)

type WorkoutManager interface {
	GetWorkouts() []model.Workout
	GetWorkout(id int) ([]model.Movement, error)
	AddMovement(workoutId int) ([]model.Movement, error)
	UpdateMovement(m model.Movement) ([]model.Movement, error)
	DeleteMovement(id int, workoutId int) ([]model.Movement, error)
	AddWorkout() error
	UpdateWorkout(e model.Workout) ([]model.Movement, error)
	//DeleteWorkout(id int) ([]model.Workout, error)
	DeleteWorkout(id int) error
}

type WorkoutService struct {
	service WorkoutManager
}

func NewWorkoutService(service WorkoutManager) *WorkoutService {
	return &WorkoutService{service: service}
}

func (ws *WorkoutService) GetAll() []model.Workout {
	return ws.service.GetWorkouts()
}

func (ws *WorkoutService) Get(r *http.Request) ([]model.Movement, error) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return ws.service.GetWorkout(int(id))
}

func (ws *WorkoutService) Add() error {
	err := ws.service.AddWorkout()
	return err
}

func (ws *WorkoutService) Delete(r *http.Request) error {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := ws.service.DeleteWorkout(int(id))
	return err
}

func (ws *WorkoutService) AddMovement(r *http.Request) ([]model.Movement, error) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	return ws.service.AddMovement(int(id))
}

func (ws *WorkoutService) DeleteMovement(r *http.Request) ([]model.Movement, error) {
	workoutIdStr := r.PathValue("workoutId")
	movementIdStr := r.PathValue("movementId")
	movementId, _ := strconv.ParseInt(movementIdStr, 10, 64)
	workoutId, _ := strconv.ParseInt(workoutIdStr, 10, 64)

	return ws.service.DeleteMovement(int(movementId), int(workoutId))
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
