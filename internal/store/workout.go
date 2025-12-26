package store

import (
	"fmt"
	"time"

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
	var workouts []model.Workout

	q := `
        SELECT id, date
        FROM workout
        ORDER BY date DESC
    `

	stmt, err := ws.env.Db.Prepare(q)
	if err != nil {
		fmt.Println("error preparing GetWorkouts: ", err)
		return workouts
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("error querying GetWorkouts: ", err)
		return workouts
	}
	defer rows.Close()

	for rows.Next() {
		var w model.Workout
		var unixInt int64
		err = rows.Scan(&w.ID, &unixInt)

		w.Date = time.Unix(unixInt, 0)

		if err != nil {
			if err == rows.Err() {
				fmt.Println("error scanning GetWorks: ", err)
			}
		}

		workouts = append(workouts, w)
	}

	return workouts
}

func (ws *WorkoutStore) GetWorkout(id int) []model.Movement {
	var movements []model.Movement

	q := `
        SELECT id, sets, reps, date
        FROM movement
        WHERE workout_id = ?
    `

	stmt, err := ws.env.Db.Prepare(q)
	if err != nil {
		fmt.Println("error preparing GetWorkout: ", err)
		return movements
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		fmt.Println("error querying GetWorkout: ", err)
		return movements
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Movement
		var unixInt int64

		err := rows.Scan(&m.ID, &m.Sets, &m.Reps, &unixInt)
		m.Date = time.Unix(unixInt, 0)

		if err != nil {
			if err == rows.Err() {
				fmt.Println("error scanning GetWorkout: ", err)
			}
		}

		movements = append(movements, m)
	}

	return movements
}

func (ws *WorkoutStore) GetLastWorkout() []model.Movement {
	var movements []model.Movement

	q := `
        SELECT id, workout_id, exercise_id, sets, reps, date
        FROM movement
        WHERE workout_id = (SELECT id FROM workout ORDER BY id DESC LIMIT 1)
    `

	stmt, err := ws.env.Db.Prepare(q)
	if err != nil {
		fmt.Println("error preparing GetLastWorkout: ", err)
		return movements
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("error querying GetLastWorkout: ", err)
		return movements
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Movement
		var unixInt int64

		err := rows.Scan(&m.ID, &m.WorkoutID, &m.ExerciseID, &m.Sets, &m.Reps, &unixInt)
		m.Date = time.Unix(unixInt, 0)

		if err != nil {
			if err == rows.Err() {
				fmt.Println("error scanning GetLastWorkout: ", err)
			}
		}

		movements = append(movements, m)
	}

	return movements
}

func (ws *WorkoutStore) AddWorkout() error {
	stmt, err := ws.env.Db.Prepare("INSERT INTO workout(date) VALUES(?)")
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = stmt.Exec(now.Unix())
	if err != nil {
		return err
	}

	return nil
}

func (ws *WorkoutStore) UpdateWorkout(e model.Workout) error {
	return nil
}

func (ws *WorkoutStore) DeleteWorkout(id int) error {
	return nil
}
