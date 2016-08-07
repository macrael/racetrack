package main

import (
	"fmt"
	"net/http"
	// "encoding/json"
	// "time"
)

type Player struct {
	Name string		`json:"name"`
	TeamName string	`json:"team_name"`
	// queens
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLO PLayers")
}