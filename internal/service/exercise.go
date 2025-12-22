package service

import (
	//"fmt"

	"github.com/section14/train-track/internal/model"
)

type ExerciseManager interface {
    GetExercises() []model.Exercise
    GetExercise(id int) model.Exercise
    AddExercise(e model.Exercise) error
    UpdateExercise(e model.Exercise) error
    DeleteExercise(id int) error
}

type ExerciseService struct {
    service ExerciseManager
}

func NewExerciseService(service ExerciseManager) *ExerciseService {
    return &ExerciseService{service: service}
}

func (es *ExerciseService) GetAll() []model.Exercise {
    exercises := es.service.GetExercises()
    return exercises
}

func (es *ExerciseService) Add(name string) error {
    e := model.Exercise {
        ID: 0,
        Name: name,
    }

    err := es.service.AddExercise(e)
    return err
}

func (es *ExerciseService) Update(id int, name string) error {
    e := model.Exercise {
        ID: id,
        Name: name,
    }

    err := es.service.UpdateExercise(e)
    return err
}

func (es *ExerciseService) Delete(id int) error {
    err := es.service.DeleteExercise(id)
    return err
}
