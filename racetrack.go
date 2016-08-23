package main

import (
    "fmt"
    "time"
    "net/http"

    "github.com/gorilla/mux"
)

func dummySeason() Season {
    s := Season{}
    s.Key = "season:2016"
    s.Year = 2016
    s.Title = "The one Bob wins"

    bob := Queen{"12", "season:2016", "Bob the drag queen"}
    chichi := Queen{"13", "season:2016", "Chi chi"}

    queens := []Queen{bob, chichi}
    s.Queens = queens

    players := []Player{Player{"MacRae", "Come a little closer"}, Player{"Oscar", "The House Always Wins"}}
    s.Players = players

    scoring := []PlayType{}
    cries := PlayType{"play_type:1", "season:2016", "cries", -2}
    shade := PlayType{"play_type:2", "season:2016", "throws shade", 4}
    scoring = append(scoring, cries)
    scoring = append(scoring, shade)
    s.PlayTypes = scoring


    one := Episode{}
    one.Number = 1
    scores := []Scored{}
    scores = append(scores, Scored{&bob, &cries, time.Now()})
    //one.Scores = scores

    two := Episode{}
    two.Number = 2
    scores = []Scored{}
    scores = append(scores, Scored{&bob, &shade, time.Now()})
    scores = append(scores, Scored{&chichi, &cries, time.Now()})
    //two.Scores = scores

    s.Episodes = []Episode{one, two}

    return s
}



func ServeIndex(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Have the site. ")
    http.ServeFile(w, r, "./static/index.html")
}

func DebugCalled(w http.ResponseWriter, r *http.Request) {
    fmt.Println("A debug call!")
    fmt.Fprintf(w, "debug this")
}

func main() {
    fmt.Println("in main")

    r := mux.NewRouter().StrictSlash(true)

    r.HandleFunc("/debug", DebugCalled)

    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/seasons", GetSeasons)
    api.HandleFunc("/seasons/{season_id}", GetSeason)
    api.HandleFunc("/seasons/{season_id}/queens", PostQueen).Methods("POST")
    api.HandleFunc("/seasons/{season_id}/queens/{queen_key}", DeleteQueen).Methods("DELETE")
    api.HandleFunc("/seasons/{season_id}/play_types", PostPlayType).Methods("POST")
    api.HandleFunc("/seasons/{season_id}/play_types/{play_type_key}", DeletePlayType).Methods("DELETE")

    api.HandleFunc("/seasons/{season_id}/episodes", PostEpisode).Methods("POST")
    api.HandleFunc("/seasons/{season_id}/episodes/{episode_key}", DeleteEpisode).Methods("DELETE")

    r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./static/"))))
    
    // Serve our single page site to any valid url
    r.HandleFunc("/", ServeIndex)
    r.HandleFunc("/queen", ServeIndex)
    r.HandleFunc("/edit", ServeIndex)
    r.HandleFunc("/edit/queens", ServeIndex)
    r.HandleFunc("/edit/play_types", ServeIndex)
    r.HandleFunc("/edit/episodes", ServeIndex)
    r.HandleFunc("/edit/episodes/{episode_id}", ServeIndex)


    fmt.Println("testing")

    s := dummySeason()
    Seasons = []Season{s}
    fmt.Println(s)

    // http.Handle("/", r)
    http.ListenAndServe(":8080", r)

}



func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("this is good")
    fmt.Fprintf(w, "Good Morning! %v", Seasons)
}

