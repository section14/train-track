package model

import "time"

type Exercise struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Movement struct {
	ID         int       `json:"id"`
	ExerciseID int       `json:"exerciseId"`
	WorkoutID  int       `json:"workoutId"`
	Sets       int       `json:"sets"`
	Reps       int       `json:"reps"`
	Date       time.Time `json:"date"`
}

type Workout struct {
	ID int `json:"id"`
	//Movements []Movement `json:"movements"`
	Date time.Time `json:"date"`
}
