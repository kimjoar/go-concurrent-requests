package concurrentRequests

import (
    "net/http";
    "io/ioutil";
    "time";
    "fmt";
)

type requests struct {
    uri string
    concurrency int
    results []Result
    closing chan chan error
}

type Requester interface {
    Close() error
    Results() []Result
}

func ConcurrentRequests(uri string, concurrency int) Requester {
    r := &requests{
        uri: uri,
        concurrency: concurrency,
        closing: make(chan chan error),
    }

    go r.loop()

    return r
}

func (r *requests) Close() error {
    errc := make(chan error)
    r.closing <- errc
    return <-errc
}

func (r *requests) Results() []Result {
    return r.results
}

type Result struct {
    Url           string
    StatusCode    int
    ContentLength int64
    ResponseTime  time.Duration
}

func (r *requests) fetch() Result {
    t0 := time.Now()
    resp, err := http.Get(r.uri)
    if err != nil {
        fmt.Println("resp failed")
    }
    defer resp.Body.Close()

    ioutil.ReadAll(resp.Body)

    return Result{r.uri, resp.StatusCode, resp.ContentLength, time.Since(t0)}
}

func (r *requests) loop() {
    c := make(chan bool, r.concurrency)
    res := make(chan Result, r.concurrency)
    var err error
    open := true

    for {
        select {
        case c <- true:
            go func() {
                result := r.fetch()
                if open {
                    res <- result
                    <- c
                }
            }()
        case result := <- res:
            r.results = append(r.results, result)
            fmt.Println(result)
        case errc := <-r.closing:
            open = false
            errc <- err
            close(res)
            close(c)
            return
        }
    }
}

func nsToMs(ns int64) int64 {
    return ns / 1000000
}

func totalResponseTime(results []Result) (time int64) {
    for _, val := range results {
        time += nsToMs(val.ResponseTime.Nanoseconds())
    }
    return
}

func ConcurrentRequestsWithStats(uri string, concurrency int, duration time.Duration) {
    fetching := ConcurrentRequests(uri, concurrency)

    // Stop fetching after some time
    time.AfterFunc(duration, func() {
        fetching.Close()
        ms := nsToMs(duration.Nanoseconds())
        reqs := len(fetching.Results())
        totalResponseTimeInMs := totalResponseTime(fetching.Results())

        fmt.Println("-----------")
        fmt.Println("Done!")
        fmt.Println("- Total Requests: ", reqs)
        fmt.Println("- Elapsed (ms): ", ms)
        fmt.Println("- reqs/s: ", float64(reqs) / duration.Seconds())
        fmt.Println("- Avg response time (ms): ", totalResponseTimeInMs / int64(reqs))
    })
}

