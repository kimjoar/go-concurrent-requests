package main

import (
    "net/http";
    "io/ioutil";
    "time";
    "fmt";
)

type Result struct {
    Url           string
    StatusCode    int
    ContentLength int64
    ResponseTime  time.Duration
}

func fetch(url string) Result {
    t0 := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("resp failed")
    }
    defer resp.Body.Close()

    ioutil.ReadAll(resp.Body)

    return Result{url, resp.StatusCode, resp.ContentLength, time.Since(t0)}
}

func concurrentRequests(url string, concurrency int) {
    c := make(chan bool, concurrency)
    res := make(chan Result, concurrency)

    for {
        select {
        case c <- true:
            go func() {
                res <- fetch(url);
                <- c
            }()
        case result := <- res:
            fmt.Println(result)
        case <- time.After(time.Second):
            fmt.Println("timeout")
        }
    }
}

func main() {
    go concurrentRequests("http://localhost:3000", 10)

    var input string
    fmt.Scanln(&input)
}
