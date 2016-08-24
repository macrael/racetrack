package main

import (
    "fmt"
    "time"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/gorilla/mux"
)

type Play struct {
    Key      string             `json:"key" redis:"-"`
    SeasonKey string            `json:"season_key" redis:"season_key"`
    QueenKey string             `json:"queen_key" redis:"queen_key"`
    PlayTypeKey string          `json:"play_type_key" redis:"play_type_key"`
    EpisodeKey string           `json:"episode_key" redis:"episode_key"`
    Timestamp int64             `json:"timestamp" redis:"timestamp"`
}

func PostPlay(w http.ResponseWriter, r *http.Request) {
    fmt.Println("PostPlay")

    var new_play Play
    b, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(b))
    json.Unmarshal(b, &new_play)


    season_key := mux.Vars(r)["season_id"]
    new_play.SeasonKey = season_key
    new_play.Timestamp = time.Now().Unix()

    success := AddObject("play", new_play)
    if success {
        fmt.Fprintf(w, "201 Created")
    } else {
        // TODO set the response code correct
        fmt.Fprintf(w, "ERR")
    }
}

func DeletePlay(w http.ResponseWriter, r *http.Request) {
    season_key := mux.Vars(r)["season_id"]
    play_key := mux.Vars(r)["play_key"]

    fmt.Println("DELTE play: ", play_key)

    var del_play Play
    del_play.Key = play_key
    success := GetObject(play_key, &del_play)
    fmt.Println("GOT QUEEN: ", del_play)
    if (season_key != del_play.SeasonKey) {
        fmt.Println("Hey, you are trying todelerte a play that doesn't belong here.")
        success = false
    }

    if (success) {
        success = DeleteObject("play", play_key, del_play)
    }

    if (success) {
        fmt.Fprintf(w, "200 OK (deleted)")
    } else {
        fmt.Fprintf(w, "KEY DOES NOT EXIST")
    }

}
