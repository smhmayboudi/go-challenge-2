# Problem 002

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
	expTime  time.Time
)

func server() int {
	if time.Now().After(expTime) {
		cacheMTX.Lock()
		if time.Now().After(expTime) {
			cache = DownStream()
			expTime = time.Now().Add(5 * time.Second)
		}
		cacheMTX.Unlock()
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
down stream called.
```

# Solution 002

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
	expTime  time.Time
)

func server() int {
	cacheMTX.Lock()
	defer cacheMTX.Unlock()

	if time.Now().After(expTime) {
		cache = DownStream()
		expTime = time.Now().Add(5 * time.Second)
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
