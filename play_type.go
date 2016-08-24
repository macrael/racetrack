package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/gorilla/mux"
)

type PlayType struct {
    Key string                  `json:"key" redis:"-"`
    SeasonKey string            `json:"season_key" redis:"season_key"`
    Action string               `json:"action" redis:"action"` // e.g. Cries
    Effect int                  `json:"effect" redis:"effect"` // +/- 5 pts. 
}

func PostPlayType(w http.ResponseWriter, r *http.Request) {
    fmt.Println("post play type")

    var new_play_type PlayType
    b, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(b))
    json.Unmarshal(b, &new_play_type)

    season_key := mux.Vars(r)["season_id"]
    new_play_type.SeasonKey = season_key

    success := AddObject("play_type", new_play_type)
    if success {
        fmt.Fprintf(w, "201 Created")
    } else {
        // TODO set the response code correct
        fmt.Fprintf(w, "ERR")
    }

}

func DeletePlayType(w http.ResponseWriter, r *http.Request){

    season_key := mux.Vars(r)["season_id"] //TODO: Verify that the play to be deleted belongs to this season
    play_type_key := mux.Vars(r)["play_type_key"]

    fmt.Println("DELTE Play: ", play_type_key)

    var del_play_type PlayType
    del_play_type.Key = play_type_key
    success := GetObject(play_type_key, &del_play_type)
    fmt.Println("GOT QUEEN: ", del_play_type)
    if (season_key != del_play_type.SeasonKey) {
        fmt.Println("Hey, you are trying todelerte a play_type that doesn't belong here.")
        success = false
    }

    if (success) {
        success = DeleteObject("play_type", play_type_key, del_play_type)
    }
    
    if (success) {
        fmt.Fprintf(w, "200 OK (deleted)")
    } else {
        fmt.Fprintf(w, "KEY DOES NOT EXIST")
    }
}
