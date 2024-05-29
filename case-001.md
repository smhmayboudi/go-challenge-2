# Problem 001

```GO
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
```

Sample RUN:

```SHELL
down stream started.
104  77
 69  77
 62  77
 44  77
176  77
  5  77
 35  77
157  77
down stream called.
```

# Solution 001

```GO
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
	expTimer *time.Ticker
)

func server() int {
	cacheMTX.Lock()
	defer cacheMTX.Unlock()

	select {
	case <-expTimer.C:
		cache = DownStream()
	default:
		Use the cached value
	}

	return cache
}

func DownStream() int {
	fmt.Println("down stream started.")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Println("down stream called.")
	return rand.Intn(100)
}

func main() {
	expTimer = time.NewTicker(5 * time.Second)
	defer expTimer.Stop()

	var wg sync.WaitGroup
	count := 200
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
			fmt.Printf("%3d %3d\n", i, server())
		}(i)
	}
	wg.Wait()
}
```

Sample RUN:

```SHELL
down stream started.
down stream called.
```