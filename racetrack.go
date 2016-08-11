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

	scoring := []ScoreType{}
	cries := ScoreType{"cries", -2}
	shade := ScoreType{"throws shade", 4}
	scoring = append(scoring, cries)
	scoring = append(scoring, shade)
	s.ScoreTypes = scoring


	one := Episode{}
	one.Number = 1
	scores := []Scored{}
	scores = append(scores, Scored{&bob, &cries, time.Now()})
	one.Scores = scores

	two := Episode{}
	two.Number = 2
	scores = []Scored{}
	scores = append(scores, Scored{&bob, &shade, time.Now()})
	scores = append(scores, Scored{&chichi, &cries, time.Now()})
	two.Scores = scores

	s.Episodes = []Episode{one, two}

	return s
}



func ServeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Have the site. ")
    http.ServeFile(w, r, "./static/index.html")
}

func main() {
	fmt.Println("in main")

	r := mux.NewRouter().StrictSlash(true)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/seasons", GetSeasons)
	// api.HandleFunc("/seasons/current", GetCurrentSeason)
	api.HandleFunc("/seasons/{season_id}", GetSeason)
	api.HandleFunc("/seasons/{season_id}/queens", PostQueen).Methods("POST")
	api.HandleFunc("/seasons/{season_id}/queens/{queen_key}", DeleteQueen).Methods("DELETE")

	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./static/"))))
	
	// Serve our single page site to any valid url
	r.HandleFunc("/", ServeIndex)
	r.HandleFunc("/queen", ServeIndex)
	r.HandleFunc("/edit", ServeIndex)
	r.HandleFunc("/edit/queens", ServeIndex)


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

