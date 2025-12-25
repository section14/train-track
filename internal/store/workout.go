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
        ORDER BY date
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
        err = rows.Scan(&w.ID, &w.Date)

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
        SELECT m.id, m.sets, m.reps, m.date
        FROM workout w
        WHERE id=?
        ORDER BY m.date
        INNER JOIN
        movement m ON w.id = o.workout_id
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
