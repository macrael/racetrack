package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"

	"github.com/gorilla/mux"
)

var Seasons []Season

type Season struct {
	Id uint					`json:"id"`
	Year int				`json:"year"`
	Title string			`json:"title"`
	
	Queens []Queen			`json:"queens"`
	Players []Player		`json:"players"`
	ScoreTypes []ScoreType	`json:"score_types"`
	Episodes []Episode		`json:"episodes"`

	CreatedAt time.Time		`json:"created_at"`
}

func GetSeasons(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getSeasons: %v", Seasons)
	json, _ := json.MarshalIndent(Seasons, "", "    ")
	fmt.Fprintf(w, string(json))
}

func GetSeason(w http.ResponseWriter, r *http.Request) {
	season_id := mux.Vars(r)["season_id"]
	fmt.Println("id; ", season_id)

	fmt.Fprintf(w, "hhihi")
}

// current season?
