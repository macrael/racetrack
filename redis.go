package main

import (
    "os"
    "fmt"
    "regexp"
    "net/url"
	"strings"
    "reflect"

    "github.com/garyburd/redigo/redis"
)

func GetObject(object_key string, object interface{}) bool {
    // Runtime confirm that we were passed a pointer, not sure how to do this in the type system.
    d := reflect.ValueOf(object)
	if d.Kind() != reflect.Ptr {
        fmt.Println("HEY, you gotta pass a pointer to get anything out of GET")
        return false
    }

    conn, err := Connect()
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
func AddObject(object_name string, object interface{}) (string, bool) {
    conn, err := Connect()
    if err != nil {
        fmt.Println("ERRORED CONNECTING")
        return "", false
    }
    defer conn.Close()

    new_object_id, err := conn.Do("INCR", fmt.Sprintf("new_%v_id", object_name))
    if err != nil {
        fmt.Println("ERR getting INCR key", err)
        return "", false
    }
    new_object_key := fmt.Sprintf("%v:%v", object_name, new_object_id)
    fmt.Println("Got a new key:", new_object_key)

    args := redis.Args{new_object_key}.AddFlat(object)
    fmt.Println("CALL", args)
    _, err = conn.Do("HMSET", args...)
    if err != nil {
        fmt.Println("no set new obj", err)
        return "", false
    }

    // Reflectively add all things that end with Key to their matching babies.
    t := reflect.TypeOf(object)
    v := reflect.ValueOf(object)
    key_regex, _ := regexp.Compile(".+Key$")

    fmt.Println("HEYOOO Lets REFLECT on our life's choices")
    for i := 0; i < v.NumField(); i++ {
        ft := t.Field(i)
        fv := v.Field(i)
        fmt.Println("HIHIHI: ", string(ft.Tag.Get("belongs_to")))

        belongs_to_tag := ft.Tag.Get("belongs_to")
        // not sure what to do with this tag besides using it to squash.
        // probably, replace object_name with something else. 
        squash := string(belongs_to_tag) == "-"

        if (key_regex.MatchString(ft.Name) && !squash) {
            belongs_to_key := fmt.Sprintf("%s:%ss", fv.Interface(), object_name)
            fmt.Println("ActualBelongtoKey: ", belongs_to_key)

            _, err = conn.Do("SADD", belongs_to_key, new_object_key)
            if err != nil {
                fmt.Println("didn't add to set", err)
                // this is where we should roll back, roll it ALL back
                return "", false
            }
        }
    }
    fmt.Println("well that was reflective")

    return new_object_key, true
}

func DeleteObject(object_name string, object_key string, object interface{}) bool {
    conn, err := Connect()
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
        
        belongs_to_tag := ft.Tag.Get("belongs_to")
        // not sure what to do with this tag besides using it to squash.
        // probably, replace object_name with something else. 
        squash := string(belongs_to_tag) == "-"

        if (key_regex.MatchString(ft.Name) && !squash) {
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


// these two functions lifted from https://github.com/soveran/redisurl/blob/master/redisurl.go
func Connect() (redis.Conn, error) {
	return ConnectToURL(os.Getenv("REDIS_URL"))
}

func ConnectToURL(s string) (c redis.Conn, err error) {
	redisURL, err := url.Parse(s)

	if err != nil {
		return
	}

	auth := ""

	if redisURL.User != nil {
		if password, ok := redisURL.User.Password(); ok {
			auth = password
		}
	}

	c, err = redis.Dial("tcp", redisURL.Host)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(auth) > 0 {
		_, err = c.Do("AUTH", auth)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if len(redisURL.Path) > 1 {
		db := strings.TrimPrefix(redisURL.Path, "/")
		c.Do("SELECT", db)
	}
    fmt.Println("Connected?")
	return
}
