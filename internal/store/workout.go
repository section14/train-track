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
				fmt.Println("error scanning GetWorkouts: ", err)
			}
		}

		workouts = append(workouts, w)
	}

	return workouts
}

func (ws *WorkoutStore) GetWorkout(id int) ([]model.Movement, error) {
	var movements []model.Movement

	q := `
        SELECT id, workout_id, exercise_id, sets, reps, weight, date
        FROM movement
        WHERE workout_id = ?
    `

	stmt, err := ws.env.Db.Prepare(q)
	if err != nil {
		return movements, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return movements, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Movement
		var unixInt int64

		err := rows.Scan(&m.ID, &m.WorkoutID, &m.ExerciseID, &m.Sets, &m.Reps, &m.Weight, &unixInt)
		m.Date = time.Unix(unixInt, 0)

		if err != nil {
			if err == rows.Err() {
				fmt.Println("error scanning GetWorkout: ", err)
			}
		}

		movements = append(movements, m)
	}

	return movements, nil
}

func (ws *WorkoutStore) GetLastWorkout() []model.Movement {
	var movements []model.Movement

	q := `
        SELECT id, workout_id, exercise_id, sets, reps, weight, date
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

		err := rows.Scan(&m.ID, &m.WorkoutID, &m.ExerciseID, &m.Sets, &m.Reps, &m.Weight, &unixInt)
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
	//add a new workout
	stmt, err := ws.env.Db.Prepare("INSERT INTO workout(date) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

    now := time.Now()
	_, err = stmt.Exec(now.Unix())
	if err != nil {
		return err
	}

	return nil
}

func (ws *WorkoutStore) UpdateWorkout(e model.Workout) ([]model.Movement, error) {
	var movements []model.Movement
	return movements, nil
}

// func (ws *WorkoutStore) DeleteWorkout(id int) ([]model.Workout, error) {
func (ws *WorkoutStore) DeleteWorkout(id int) error {
	//var emptyWorkouts []model.Workout

	//todo: prevent deleting a workout if it's the only one
	stmt, err := ws.env.Db.Prepare("DELETE FROM workout WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	return err

	/*
		_, err = stmt.Exec(id)
		if err != nil {
			return emptyWorkouts, err
		}

		workouts := ws.GetWorkouts()

		return workouts, nil
	*/
}

func (ws *WorkoutStore) GetPlaceHolderExercise() (int, error) {
	var id int

	stmt, err := ws.env.Db.Prepare("SELECT id FROM exercise LIMIT 1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ws *WorkoutStore) AddMovement(workoutId int) ([]model.Movement, error) {
	var movements []model.Movement
	var id int

	placeholder, err := ws.GetPlaceHolderExercise()
	if err != nil {
		return movements, err
	}

	m := model.Movement{
		WorkoutID:  workoutId,
		ExerciseID: placeholder,
		Sets:       3,
		Reps:       5,
		Weight:     10,
		Date:       time.Now(),
	}

	q := `
        INSERT
        INTO movement(workout_id, exercise_id, sets, reps, weight, date)
        VALUES(?,?,?,?,?,?)
        RETURNING id
    `
	stmt, err := ws.env.Db.Prepare(q)
	if err != nil {
		return movements, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(m.WorkoutID, m.ExerciseID, m.Sets, m.Reps, m.Weight, m.Date).Scan(&id)
	if err != nil {
		return movements, err
	}

	//get updated workout to return
	workout, _ := ws.GetWorkout(workoutId)
	movements = workout

	return movements, nil
}

func (ws *WorkoutStore) UpdateMovement(m model.Movement) ([]model.Movement, error) {
	var emptyMovements []model.Movement

	// update a movement
	updateQ := `
        UPDATE movement
        SET exercise_id=?, sets=?, reps=?, weight=?
        WHERE id=?
    `

	updateStmt, err := ws.env.Db.Prepare(updateQ)
	if err != nil {
		return emptyMovements, err
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(m.ExerciseID, m.Sets, m.Reps, m.Weight, m.ID)
	if err != nil {
		return emptyMovements, err
	}

	// return updated movements for a workout
	moves, err := ws.GetWorkout(m.ID)

	return moves, err
}

func (ws *WorkoutStore) DeleteMovement(id int, workoutId int) ([]model.Movement, error) {
	var emptyMovements []model.Movement

	stmt, err := ws.env.Db.Prepare("DELETE FROM movement WHERE id=?")
	if err != nil {
		return emptyMovements, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return emptyMovements, err
	}

	movements, err := ws.GetWorkout(workoutId)

	return movements, err
}
