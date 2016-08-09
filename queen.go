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
	Key string			`json:"key" redis:"-"`
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

	season_key := mux.Vars(r)["season_id"]
	new_queen.SeasonKey = season_key

	new_queen_id, err := conn.Do("INCR", "new_queen_id")
	new_queen_key := fmt.Sprintf("queen:%v", new_queen_id) // little sketch, shoiuld maybe validate the input here.

	args := redis.Args{new_queen_key}.AddFlat(new_queen)
	fmt.Println("CALL", args)
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		fmt.Println("no set new queen", err)
		return
	}

	season_queens_key := fmt.Sprintf("%s:queens", season_key)
	fmt.Println(season_queens_key, new_queen_key)
	_, err = conn.Do("SADD", season_queens_key, new_queen_key)
	if err != nil {
		fmt.Println("didn't add to queens set", err)
		return
	}

	fmt.Fprintf(w, "201 Created")

}