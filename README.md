# Valve
simple key based rate limiter API using in-memory store.

**NOTE**: This package is provided "as is" with no guarantee. Use it at your own risk and always test it yourself before using it in a production environment. If you find any issues, please [create a new issue](https://github.com/twiny/valve/issues/new).

## Install
`go get github.com/twiny/valve`

## API
```go
package main

import (
	"fmt"

	"github.com/twiny/valve"
)

func main() {
    // rate:    10 requests per second
    // bursts:  5 permits bursts
    // ttl:     reset every 10 minutes - if ttl less then 1 minutes, it will be set to 1 minute.
	limiter := valve.NewLimiter(10, 5, 10*time.Minute)
	defer limiter.Close()
	key := "127.0.0.1" // some key

	for i := 0; i < 100; i++ {
		if !limiter.Allow(key) {
			fmt.Println("slow down")
		}
		fmt.Println("passed")
	}
}
// output:
/*
passed
passed
passed
passed
passed
slow down
passed
slow down
passed
slow down
passed
...
*/
``` 