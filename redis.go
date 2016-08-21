package main

import (
    "fmt"

    "github.com/garyburd/redigo/redis"
)

// TODO: actually return something different on err
func AddObject(season_key string, object_name string, object interface{}) bool {
   conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("ERRORED CONNECTING")
        return false
    }
    defer conn.Close()

    new_object_id, err := conn.Do("INCR", fmt.Sprintf("new_%v_id", object_name))
    if err != nil {
        fmt.Println("ERR getting INCR key", err)
        return false
    }
    new_object_key := fmt.Sprintf("%v:%v", object_name, new_object_id)
    fmt.Println("Got a new key:", new_object_key)

    args := redis.Args{new_object_key}.AddFlat(object)
    fmt.Println("CALL", args)
    _, err = conn.Do("HMSET", args...)
    if err != nil {
        fmt.Println("no set new obj", err)
        return false
    }

    season_objects_key := fmt.Sprintf("%s:%ss", season_key, object_name)
    fmt.Println("seasonObjetKey: ", season_objects_key)
    
    _, err = conn.Do("SADD", season_objects_key, new_object_key)
    if err != nil {
        fmt.Println("didn't add to set", err)
        // this is where we should roll back
        return false
    }

    return true
}

func DeleteObject(season_key string, object_name string, object_key string) bool {
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("ERRORED CONNECTING")
        return false
    }
    defer conn.Close()

    _, err = redis.Bool(conn.Do("DEL", object_key))
    if err != nil {
        fmt.Println("donezo")
        return false
    }
    
    season_objects_key := fmt.Sprintf("%s:%ss", season_key, object_name)
    fmt.Println("trytotdeL", season_objects_key)
    _, err = redis.Bool(conn.Do("SREM", season_objects_key, object_key))
    if err != nil {
        fmt.Println("Lonzores")
        return false
    }

    return true
}
