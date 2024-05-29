package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	cache    int
	cacheMTX sync.Mutex
	expTime  = time.NewTicker(5 * time.Second)
)

func server() int {
	select {
	case <-expTime.C:
		cacheMTX.Lock()
		cache = DownStream()
		cacheMTX.Unlock()
		expTime.Reset(5 * time.Second)
		return cache
	default:
		return cache
	}
}

func DownStream() int {
	fmt.Println("down stream started.")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Println("down stream called.")
	return rand.Intn(100)
}

func main() {
	count := 200
	for i := 0; i < count; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
			fmt.Printf("%3d %3d\n", i, server())
		}(i)
	}
	time.Sleep(5 * time.Minute)
}
