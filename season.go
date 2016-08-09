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
	Key string				`json:"key" redis:"-"`
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
	var season_key string
	if season_id == "current" {
		season_key, err = redis.String(conn.Do("GET", "current_season_key"))
		if err != nil {
			fmt.Println("ERRORED", err)
		    return
		}
	} else {
		season_key = fmt.Sprintf("season:%s", season_id)
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

	season_queens_key := fmt.Sprintf("%s:queens", season_key)
	queen_keys, err := redis.Strings(conn.Do("SMEMBERS", season_queens_key))
	if err != nil {
		fmt.Printf("ERR3", err)
		return
	}

	season.Queens = []Queen{}
	for _, queen_key := range queen_keys {
		fmt.Println(queen_key)
		var queen Queen

		queen.Key = queen_key
		v, err := redis.Values(conn.Do("HGETALL", queen_key))
		if err != nil {
			fmt.Printf("ERR1", err)
			return
		}
		err = redis.ScanStruct(v, &queen)
		if err != nil {
			fmt.Printf("ERR2", err)
			return
		}

		season.Queens = append(season.Queens, queen)

	}

	fmt.Println("No more")

	json, _ := json.MarshalIndent(season, "", "    ")
	fmt.Fprintf(w, string(json))
}

