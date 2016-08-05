package main

import "time"

type Queen struct{
	Name string		`json:"name"`
}

type Player struct {
	Name string		`json:"name"`
	TeamName string	`json:"team_name"`
}

type Episode struct {
	Number uint		`json:"number"`
	Scores []Scored	`json:"scores"`
}

type ScoreType struct {
	Name string		`json:"name"` // e.g. Cries
	Effect int  	`json:"effect"` // +/- 5 pts. 
}

type Scored struct {
	Queen *Queen			`json:"queen"`
	ScoreType *ScoreType	`json:"sore_type"`
	Timestamp time.Time		`json:"timestamp"`
}