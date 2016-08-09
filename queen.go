package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/garyburd/redigo/redis"
)

type Queen struct{
	SeasonKey string	`json:"season_key" redis:"season_key"`
	Name string			`json:"name" redis:"name"`
}

func PostQueen(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostQueen")
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
	    fmt.Println("ERRORED CONNECTING")
	    return
	}
	defer conn.Close()

	var new_queen Queen
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
	json.Unmarshal(b, &new_queen)

	season_id := mux.Vars(r)["season_id"]
	season_key := fmt.Sprintf("season:%s", season_id)
	new_queen.SeasonKey = season_key

	new_queen_id, err := conn.Do("INCR", "new_queen_id")
	new_queen_key := fmt.Sprintf("queen:%v", new_queen_id) // little sketch, shoiuld maybe validate the input here.

	argss := redis.Args{new_queen_key}.AddFlat(new_queen)
	fmt.Println("args: ", argss)
	_, err = conn.Do("HMSET", argss...)
	if err != nil {
		fmt.Println("no set", err)
		return
	}

	// NEXT: Add new queenID to denomralized season:id:queens list(set? so we can del)

	fmt.Fprintf(w, "201 Created")

}