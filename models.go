package main

import "time"


type Scored struct {
    Queen *Queen                `json:"queen"`
    PlayType *PlayType         `json:"play_type"`
    Timestamp time.Time     `json:"timestamp"`
}
