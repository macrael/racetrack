package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/gorilla/mux"
)

type Queen struct{
    Key string          `json:"key" redis:"-"`
    SeasonKey string    `json:"season_key" redis:"season_key"`
    Name string         `json:"name" redis:"name"`
}

func PostQueen(w http.ResponseWriter, r *http.Request) {
    fmt.Println("PostQueen")

    var new_queen Queen
    b, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(b))
    json.Unmarshal(b, &new_queen)

    season_key := mux.Vars(r)["season_id"]
    new_queen.SeasonKey = season_key

    _, success := AddObject("queen", new_queen)
    if success {
        fmt.Fprintf(w, "201 Created")
    } else {
        // TODO set the response code correct
        fmt.Fprintf(w, "ERR")
    }
}

func DeleteQueen(w http.ResponseWriter, r *http.Request) {
    season_key := mux.Vars(r)["season_id"]
    queen_key := mux.Vars(r)["queen_key"]

    fmt.Println("DELTE QUEN: ", queen_key)

    var del_queen Queen
    del_queen.Key = queen_key
    success := GetObject(queen_key, &del_queen)
    fmt.Println("GOT QUEEN: ", del_queen)
    if (season_key != del_queen.SeasonKey) {
        fmt.Println("Hey, you are trying todelerte a queen that doesn't belong here.")
        success = false
    }

    if (success) {
        success = DeleteObject("queen", queen_key, del_queen)
    }

    //TODO: If you delete a queen, we don't clean up references to it in players (at least)...

    if (success) {
        fmt.Fprintf(w, "200 OK (deleted)")
    } else {
        fmt.Fprintf(w, "KEY DOES NOT EXIST")
    }

}
