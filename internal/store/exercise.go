package store

import (
	"fmt"

	"github.com/section14/train-track/internal/config"
	"github.com/section14/train-track/internal/model"
)

type ExerciseStore struct {
	env *config.Env
}

func NewExerciseStore(env *config.Env) *ExerciseStore {
	return &ExerciseStore{env: env}
}

func (es *ExerciseStore) GetExercises() []model.Exercise {
	var exercises []model.Exercise

	stmt, err := es.env.Db.Prepare("SELECT id, name FROM exercise ORDER BY name")
	if err != nil {
		fmt.Println("error preparing GetExercises: ", err)
		return exercises
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("error querying GetExercises: ", err)
		return exercises
	}
	defer rows.Close()

	for rows.Next() {
		var e model.Exercise
		err = rows.Scan(&e.ID, &e.Name)
		if err != nil {
			if err == rows.Err() {
				fmt.Println("error scanning GetExercises: ", err)
			}
		}

		exercises = append(exercises, e)
	}

	return exercises
}

func (es *ExerciseStore) GetExercise(id int) model.Exercise {
	var e model.Exercise

	stmt, err := es.env.Db.Prepare("SELECT id, name FROM exercise WHERE id=?")
	if err != nil {
		fmt.Println("error preparing GetExercise: ", err)
		return e
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&e.ID, &e.Name)
	if err != nil {
		fmt.Println("error querying GetExercise: ", err)
		return e
	}

	return e
}

func (es *ExerciseStore) AddExercise(e model.Exercise) error {
	stmt, err := es.env.Db.Prepare("INSERT INTO exercise(name) VALUES(?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(e.Name)
	if err != nil {
		return err
	}

	return nil
}

func (es *ExerciseStore) UpdateExercise(e model.Exercise) error {
	stmt, err := es.env.Db.Prepare("UPDATE exercise SET name=? WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(e.Name, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (es *ExerciseStore) DeleteExercise(id int) error {
	stmt, err := es.env.Db.Prepare("DELETE FROM exercise WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
