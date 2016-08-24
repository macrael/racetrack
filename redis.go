package main

import (
    "fmt"
    "reflect"
    "regexp"

    "github.com/garyburd/redigo/redis"
)

func GetObject(object_key string, object interface{}) bool {
    // Runtime confirm that we were passed a pointer, not sure how to do this in the type system.
    d := reflect.ValueOf(object)
	if d.Kind() != reflect.Ptr {
        fmt.Println("HEY, you gotta pass a pointer to get anything out of GET")
        return false
    }

    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("ERRORED CONNECTING")
        return false
    }
    defer conn.Close()
    
    fmt.Println("get object: ", object_key)

    v, err := redis.Values(conn.Do("HGETALL", object_key))
    if err != nil {
        fmt.Printf("ERR1", err)
        return false
    }
    err = redis.ScanStruct(v, object)
    if err != nil {
        fmt.Printf("ERR2", err)
        return false
    }

    return true
}

// TODO: actually return something different on err
func AddObject(object_name string, object interface{}) bool {
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

    // Reflectively add all things that end with Key to their matching babies.
    t := reflect.TypeOf(object)
    v := reflect.ValueOf(object)
    key_regex, _ := regexp.Compile(".+Key$")

    fmt.Println("HEYOOO Lets REFLECT on our life's choices")
    for i := 0; i < v.NumField(); i++ {
        ft := t.Field(i)
        fv := v.Field(i)

        if (key_regex.MatchString(ft.Name)) {
            belongs_to_key := fmt.Sprintf("%s:%ss", fv.Interface(), object_name)
            fmt.Println("ActualBelongtoKey: ", belongs_to_key)

            _, err = conn.Do("SADD", belongs_to_key, new_object_key)
            if err != nil {
                fmt.Println("didn't add to set", err)
                // this is where we should roll back, roll it ALL back
                return false
            }


        }

    }
    fmt.Println("well that was reflective")

    return true
}

func DeleteObject(object_name string, object_key string, object interface{}) bool {
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

    t := reflect.TypeOf(object)
    v := reflect.ValueOf(object)
    key_regex, _ := regexp.Compile(".+Key$")

    fmt.Println("HEYOOO Lets REFLECT on our life's choices")
    for i := 0; i < v.NumField(); i++ {
        ft := t.Field(i)
        fv := v.Field(i)

        if (key_regex.MatchString(ft.Name)) {
            belongs_to_key := fmt.Sprintf("%s:%ss", fv.Interface(), object_name)
            fmt.Println("ActualBelongtoKey: ", belongs_to_key)

            _, err = conn.Do("SREM", belongs_to_key, object_key)
            if err != nil {
                fmt.Println("didn't rem to set", err)
                // this is where we should roll back, roll it ALL back
                return false
            }
        }
    }
    
    return true
}
