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
    Key string              `json:"key" redis:"-"`
    Year int                `json:"year" redis:"year"`
    Title string            `json:"title" redis:"title"`
    
    Queens []Queen          `json:"queens"`
    Players []Player        `json:"players"`
    PlayTypes []PlayType    `json:"play_types"`
    Episodes []Episode      `json:"episodes"`
    Plays []Play            `json:"plays"` // TODO: should this be hanging off everything?

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

    season_key := mux.Vars(r)["season_id"]
    if season_key == "current" {
        season_key, err = redis.String(conn.Do("GET", "current_season_key"))
        if err != nil {
            fmt.Println("ERRORED", err)
            return
        }
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

    // Queens
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

    // PlayTypes
    season_play_types_key := fmt.Sprintf("%s:play_types", season_key)
    play_type_keys, err := redis.Strings(conn.Do("SMEMBERS", season_play_types_key))
    if err != nil {
        fmt.Printf("ERR3", err)
        return
    }

    season.PlayTypes = []PlayType{}
    for _, play_type_key := range play_type_keys {
        fmt.Println(play_type_key)
        var play_type PlayType

        play_type.Key = play_type_key
        v, err := redis.Values(conn.Do("HGETALL", play_type_key))
        if err != nil {
            fmt.Printf("ERR11", err)
            return
        }
        err = redis.ScanStruct(v, &play_type)
        if err != nil {
            fmt.Printf("ERR22", err)
            return
        }

        season.PlayTypes = append(season.PlayTypes, play_type)
    }

    // Episodes
    season_episodes_key := fmt.Sprintf("%s:episodes", season_key)
    episode_keys, err := redis.Strings(conn.Do("SMEMBERS", season_episodes_key))
    if err != nil {
        fmt.Printf("ERR3", err)
        return
    }

    season.Episodes = []Episode{}
    for _, episode_key := range episode_keys {
        fmt.Println(episode_key)
        var episode Episode

        episode.Key = episode_key
        v, err := redis.Values(conn.Do("HGETALL", episode_key))
        if err != nil {
            fmt.Printf("ERR11", err)
            return
        }
        err = redis.ScanStruct(v, &episode)
        if err != nil {
            fmt.Printf("ERR22", err)
            return
        }

        season.Episodes = append(season.Episodes, episode)
    }

    // Plays
    season_plays_key := fmt.Sprintf("%s:plays", season_key)
    play_keys, err := redis.Strings(conn.Do("SMEMBERS", season_plays_key))
    if err != nil {
        fmt.Printf("ERR3", err)
        return
    }

    season.Plays = []Play{}
    for _, play_key := range play_keys {
        fmt.Println(play_key)
        var play Play

        play.Key = play_key
        v, err := redis.Values(conn.Do("HGETALL", play_key))
        if err != nil {
            fmt.Printf("ERR1", err)
            return
        }
        err = redis.ScanStruct(v, &play)
        if err != nil {
            fmt.Printf("ERR2", err)
            return
        }

        season.Plays = append(season.Plays, play)
    }
    

    fmt.Println("No more")

    json, _ := json.MarshalIndent(season, "", "    ")
    fmt.Fprintf(w, string(json))
}

