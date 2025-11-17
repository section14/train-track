package service

import (
	"fmt"

	"github.com/section14/train-track/internal/model"
)

type WorkoutManager interface {
	GetWorkouts() []model.Workout
	GetWorkout(id int) model.Workout
	AddWorkout(e model.Workout) error
	UpdateWorkout(e model.Workout) error
	DeleteWorkout(id int) error
}

type WorkoutService struct {
	service WorkoutManager
}

func NewWorkoutService(service WorkoutManager) *WorkoutService {
	return &WorkoutService{service: service}
}

func (ws *WorkoutService) GetAll() {
    workouts := ws.service.GetWorkouts()
    fmt.Println("workouts", workouts)
}
