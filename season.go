package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/garyburd/redigo/redis"
)

var Seasons []Season

type Season struct {
	Key string				`json:"key"`
	Year int				`json:"year" redis:"year"`
	Title string			`json:"title" redis:"title"`
	
	Queens []Queen			`json:"queens"`
	Players []Player		`json:"players"`
	ScoreTypes []ScoreType	`json:"score_types"`
	Episodes []Episode		`json:"episodes"`

}

func GetSeasons(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getSeasons: %v", Seasons)
	json, _ := json.MarshalIndent(Seasons, "", "    ")
	fmt.Fprintf(w, string(json))
}

func GetSeason(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
	    fmt.Println("ERRORED CONNECTING")
	    return
	}
	defer conn.Close()

	season_id := mux.Vars(r)["season_id"]
	season_key := fmt.Sprintf("season:%s", season_id)

	var season Season
	season.Key = season_key
	v, err := redis.Values(conn.Do("HGETALL", season_key))
	if err != nil {
		fmt.Printf("ERR1", err)
		return
	}
	err = redis.ScanStruct(v, &season)
	if err != nil {
		fmt.Printf("ERR2", err)
		return
	}
	json, _ := json.MarshalIndent(season, "", "    ")
	fmt.Fprintf(w, string(json))
}

func GetCurrentSeason(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start Current Season")

	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
	    fmt.Println("ERRORED CONNECTING")
	    return
	}
	defer conn.Close()

	season_key, err := redis.String(conn.Do("GET", "current_season_key"))
	if err != nil {
		fmt.Println("ERRORED", err)
	    return
	}

	var season Season
	season.Key = season_key
	v, err := redis.Values(conn.Do("HGETALL", season_key))
	if err != nil {
		fmt.Printf("ERR1", err)
		return
	}
	err = redis.ScanStruct(v, &season)
	if err != nil {
		fmt.Printf("ERR2", err)
		return
	}
	json, _ := json.MarshalIndent(season, "", "    ")
	fmt.Fprintf(w, string(json))	

}

// current season?
