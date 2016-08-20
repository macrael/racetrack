package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/gorilla/mux"
    "github.com/garyburd/redigo/redis"
)

type PlayType struct {
    Key string                  `json:"key" redis:"-"`
    SeasonKey string     `json:"season_key" redis:"season_key"`
    Action string             `json:"action" redis:"action"` // e.g. Cries
    Effect int                    `json:"effect" redis:"effect"` // +/- 5 pts. 
}

func PostPlayType(w http.ResponseWriter, r *http.Request) {
    fmt.Println("post play type")
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("ERRORED CONNECTING")
        return
    }
    defer conn.Close()

    var new_play_type PlayType
    b, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(b))
    json.Unmarshal(b, &new_play_type)

    season_key := mux.Vars(r)["season_id"]
    new_play_type.SeasonKey = season_key

    new_play_type_id, err := conn.Do("INCR", "new_play_type_id")
    // if err
    new_play_type_key := fmt.Sprintf("play_type:%v", new_play_type_id) // little sketch, shoiuld maybe validate the input here.

    args := redis.Args{new_play_type_key}.AddFlat(new_play_type)
    fmt.Println("CALL", args)
    _, err = conn.Do("HMSET", args...)
    if err != nil {
        fmt.Println("no set new platyep", err)
        return
    }

    season_play_types_key := fmt.Sprintf("%s:play_types", season_key)
    fmt.Println(season_play_types_key, new_play_type_key)
    _, err = conn.Do("SADD", season_play_types_key, new_play_type_key)
    if err != nil {
        fmt.Println("didn't add to play_types set", err)
        return
    }

    fmt.Fprintf(w, "201 Created")

}

func DeletePlayType(w http.ResponseWriter, r *http.Request){
    fmt.Println("del play type")
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("ERRORED CONNECTING")
        return
    }
    defer conn.Close()

    season_key := mux.Vars(r)["season_id"] //TODO: Verify that the play to be deleted belongs to this season
    play_type_key := mux.Vars(r)["play_type_key"]

    fmt.Println("DELTE Play: ", play_type_key)

    success, err := redis.Bool(conn.Do("DEL", play_type_key))
    if err != nil {
        fmt.Println("donezo")
        return
    }

    season_play_types_key := fmt.Sprintf("%s:play_types", season_key)
    success, err = redis.Bool(conn.Do("SREM", season_play_types_key, play_type_key))
    if err != nil {
        fmt.Println("Lonzores")
        return
    }

    if (success) {
        fmt.Fprintf(w, "200 OK (deleted)")
    } else {
        fmt.Fprintf(w, "KEY DOES NOT EXIST")
    }
}
