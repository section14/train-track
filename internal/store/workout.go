package store

import (
	"github.com/section14/train-track/internal/config"
	"github.com/section14/train-track/internal/model"
)

type WorkoutStore struct {
	env *config.Env
}

func NewWorkoutStore(env *config.Env) *WorkoutStore {
	return &WorkoutStore{env: env}
}

func (ws *WorkoutStore) GetWorkouts() []model.Workout {
    return []model.Workout{}
}

func (ws *WorkoutStore) GetWorkout(id int) model.Workout {
    return model.Workout{}
}

func (ws *WorkoutStore) AddWorkout(e model.Workout) error {
    return nil
}

func (ws *WorkoutStore) UpdateWorkout(e model.Workout) error {
    return nil
}

func (ws *WorkoutStore) DeleteWorkout(id int) error {
    return nil
}
