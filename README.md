# go-challenge-2

In this Golang-based challenge, we will demonstrate the correct way of implementing an efficient caching mechanism to handle simple HTTP requests. The main() function will act as the client, sending requests to the Server struct. The Server will respond to each request by retrieving the data from a cache.

For this scenario, the response data is a simple integer value. We will introduce a simple rule to invalidate the cache, which is to clear the cache every 5 seconds. Whenever the cache is invalidated, the Server will call another service named DownStream to fetch the updated data.

To simulate a task that takes some time to complete, we will introduce a random time.Sleep() call within the DownStream function. Additionally, we will include some logging statements throughout the code to showcase the correctness and flow of our solution.

The key aspects of this implementation will be:

Designing the Server struct with a cache to store the response data.
Implementing a mechanism to invalidate the cache at the appropriate intervals.
Handling the communication between the Server and the DownStream service to fetch the updated data.
Ensuring the correct handling of concurrent requests and thread-safe access to the cache.
Incorporating the logging statements to provide visibility into the system's behavior.
By following this approach, we will create a Golang-based solution that efficiently caches the responses and gracefully handles the dynamic nature of the data source.

# Solution 1

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

The issue with the provided code is related to the use of time.Ticker and the way the cache invalidation logic is implemented.

1. Incorrect usage of time.Ticker: The expTime variable is declared as a time.Ticker, but it is not being properly used to manage the cache expiration. The time.Ticker should be used to trigger the cache invalidation at the correct intervals, but in the current implementation, it is not being reset or stopped when the program exits.

2. Potential race conditions: The code uses a sync.Mutex to protect the cache variable, but the locking and unlocking are not done correctly. The cacheMTX.Lock() and cacheMTX.Unlock() calls should be placed around the entire block of code that accesses the cache variable, not just the DownStream() call.

3. Inefficient cache invalidation: The cache invalidation logic is implemented within the server() function, but it is not being triggered correctly. The select statement should be checking the expTime.C channel for new ticks, but instead, it is only checking the default case, which means the cache will never be invalidated.

4. Incorrect goroutine creation: The main() function creates 200 goroutines, each of which calls the server() function. However, the server() function does not have any mechanism to handle concurrent requests. This could lead to race conditions and incorrect cache updates.
   Here's a corrected version of the code that addresses these issues:

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
		// Use the cached value
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

# Solution 2

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

The main issue with the provided code is the way the cache invalidation logic is implemented. The code has a race condition that can lead to incorrect cache updates.

Here's the breakdown of the problem:

1. Race condition in cache update: The server() function checks the expTime and updates the cache only if the current time is after the expTime. However, this check is not atomic, and there is a race condition between the time when the check is performed and the time when the cache is updated. Consider the following scenario:

- Multiple goroutines call the server() function concurrently.
- One of the goroutines checks the expTime and finds that it has expired.
- Before that goroutine can acquire the cacheMTX lock and update the cache, another goroutine may have already updated the cache.
- The first goroutine then proceeds to update the cache, overwriting the value set by the other goroutine.

2. Inefficient cache invalidation: The cache invalidation logic is implemented within the server() function, but it is not being triggered correctly. The code checks the expTime on every call to the server() function, which can be inefficient if the cache is accessed frequently.

3. Incorrect goroutine creation: The main() function creates 200 goroutines, each of which calls the server() function. However, the server() function does not have any mechanism to handle concurrent requests. This could lead to race conditions and incorrect cache updates.
   Here's a corrected version of the code that addresses these issues:

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
