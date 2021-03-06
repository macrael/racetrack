package main

import (
    "os"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

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

    port := os.Getenv("PORT")

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

    api.HandleFunc("/seasons/{season_id}/plays", PostPlay).Methods("POST")
    api.HandleFunc("/seasons/{season_id}/plays/{play_key}", DeletePlay).Methods("DELETE")

    api.HandleFunc("/seasons/{season_id}/players", PostPlayer).Methods("POST")
    api.HandleFunc("/seasons/{season_id}/players/{player_key}", DeletePlayer).Methods("DELETE")

    r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./static/"))))
    
    // Serve our single page site to any valid url // and really, maybe any non-api url, let 404 be handled by the app. 
    r.HandleFunc("/", ServeIndex)
    r.HandleFunc("/queens", ServeIndex)
    r.HandleFunc("/edit", ServeIndex)
    r.HandleFunc("/edit/queens", ServeIndex)
    r.HandleFunc("/edit/play_types", ServeIndex)
    r.HandleFunc("/edit/episodes", ServeIndex)
    r.HandleFunc("/edit/episodes/{episode_id}", ServeIndex)
    r.HandleFunc("/edit/players", ServeIndex)


    fmt.Println("testing")

    http.ListenAndServe(":" + port, r)
    fmt.Println("we returned?")

}



func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("this is good")
    fmt.Fprintf(w, "Good Morning! %v", Seasons)
}

