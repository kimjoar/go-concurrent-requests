Perform concurrent requests
===========================

Just playing with Go. Simple code to perform concurrent requests on a
URL.

Can be used like this, `main.go`:

```go
package main

import (
    req "github.com/kjbekkelund/go-concurrent-requests"
    "fmt";
    "time";
)

func main() {
    req.ConcurrentRequestsWithStats("http://localhost:3000", 30, 3*time.Second)

    var input string
    fmt.Scanln(&input)
}
```

Remember to call `go get .` in the same folder as the file.
