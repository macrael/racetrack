package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/gorilla/mux"
)

type Episode struct {
    Key string          `json:"key" redis:"-"`
    SeasonKey string    `json:"season_key" redis:"season_key"`
    Number uint         `json:"number" redis:"number"`
}

func PostEpisode(w http.ResponseWriter, r *http.Request) {
    fmt.Println("PostEpisode")

    var new_episode Episode
    b, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(b))
    json.Unmarshal(b, &new_episode)

    season_key := mux.Vars(r)["season_id"]
    new_episode.SeasonKey = season_key

    _, success := AddObject("episode", new_episode)
    if success {
        fmt.Fprintf(w, "201 Created")
    } else {
        // TODO set the response code correct
        fmt.Fprintf(w, "ERR")
    }
}

func DeleteEpisode(w http.ResponseWriter, r *http.Request) {
    season_key := mux.Vars(r)["season_id"] //TODO: Verify that the episode to be deleted belongs to this season
    episode_key := mux.Vars(r)["episode_key"]

    fmt.Println("DELTE EPI: ", episode_key)

    var del_episode Episode
    del_episode.Key = episode_key
    success := GetObject(episode_key, &del_episode)
    fmt.Println("GOT QUEEN: ", del_episode)
    if (season_key != del_episode.SeasonKey) {
        fmt.Println("Hey, you are trying todelerte a episode that doesn't belong here.")
        success = false
    }

    if (success) {
        success = DeleteObject("episode", episode_key, del_episode)
    }

    if (success) {
        fmt.Fprintf(w, "200 OK (deleted)")
    } else {
        fmt.Fprintf(w, "KEY DOES NOT EXIST")
    }

}
