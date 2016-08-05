package main

import (
	"fmt"
	"time"
	"net/http"

	"github.com/gorilla/mux"
)

func dummySeason() Season {
	s := Season{}
	s.Id = 1
	s.Year = 2016
	s.Title = "The one Bob wins"

	bob := Queen{"Bob the drag queen"}
	chichi := Queen{"Chi chi"}

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


func main() {
	fmt.Println("in main")

	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/seasons", GetSeasons)


	fmt.Println("testing")

	s := dummySeason()
	Seasons = []Season{s}
	fmt.Println(s)

	http.ListenAndServe(":8080", r)

}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is good")
	fmt.Fprintf(w, "Good Morning! %v", Seasons)
}

