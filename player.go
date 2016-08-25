package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/gorilla/mux"
    "github.com/garyburd/redigo/redis"
)

type Player struct {
    Key      string       `json:"key" redis:"-"`
    SeasonKey string      `json:"season_key" redis:"season_key"`
    Name string           `json:"name" redis:"name"`
    TeamName string       `json:"team_name" redis:"team_name"`
    QueenKeys []string    `json:"queen_keys" redis:"-"`
    WinnerKey string      `json:"winner_key" redis:"winner_key" belongs_to:"-"`
}

func PostPlayer(w http.ResponseWriter, r *http.Request) {
    fmt.Println("PostPlayer")

    var new_player Player
    b, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(b))
    json.Unmarshal(b, &new_player)

    fmt.Println("NewPlayer:", new_player)

    season_key := mux.Vars(r)["season_id"]
    new_player.SeasonKey = season_key

    new_player_key, success := AddObject("player", new_player)
    new_player.Key = new_player_key

    if (success) {
        conn, err := redis.Dial("tcp", ":6379")
        if err != nil {
            fmt.Println("ERRORED CONNECTING")
            success = false
        }
        defer conn.Close()

        // add the queen keys to redis
        for _, queen_key := range new_player.QueenKeys {
            player_queens_key := fmt.Sprintf("%s:queens", new_player.Key)
            fmt.Println("adding qs to ", player_queens_key)

            _, err = conn.Do("SADD", player_queens_key, queen_key)
            if err != nil {
                fmt.Println("didn't add player's queen to set!!", err)
                // this is where we should roll back, roll it ALL back
                success = false
            }
        }
    }

    if (success) {
        fmt.Fprintf(w, "201 Created")
    } else {
        fmt.Fprintf(w, "Error, not created")
    }
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {

    season_key := mux.Vars(r)["season_id"]
    player_key := mux.Vars(r)["player_key"]

    fmt.Println("DELTE PLAR: ", player_key)

    var del_player Player
    del_player.Key = player_key
    success := GetObject(player_key, &del_player)
    fmt.Println("GOT Player: ", del_player)
    if (success && season_key != del_player.SeasonKey) {
        fmt.Println("Hey, you are trying todelerte a player that doesn't belong here.")
        success = false
    }

    if (success) {
        success = DeleteObject("player", player_key, del_player)
        if (success) {
            conn, err := redis.Dial("tcp", ":6379")
            if err != nil {
                fmt.Println("ERRORED CONNECTING")
                success = false
                // this is dumb, and probably dumb above too. 
            }
            defer conn.Close()
            //delete the queenkeys
            player_queens_key := fmt.Sprintf("%s:queens", del_player.Key)
            _, err = redis.Bool(conn.Do("DEL", player_queens_key))
            if err != nil {
                fmt.Println("donrrrezo")
                success = false
            }
        }
    }

    if (success) {
        fmt.Fprintf(w, "200 OK (deleted)")
    } else {
        fmt.Fprintf(w, "KEY DOES NOT EXIST")
    }

}
