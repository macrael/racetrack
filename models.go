package main

import "time"

type Episode struct {
    Number uint     `json:"number"`
    Scores []Scored `json:"scores"`
    // summary?
}

type Scored struct {
    Queen *Queen                `json:"queen"`
    PlayType *PlayType         `json:"play_type"`
    Timestamp time.Time     `json:"timestamp"`
}