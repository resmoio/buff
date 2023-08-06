# buff Library

The buff library provides a generic buffered data structure in Go that supports automatic flushing 
based on size or time. There are a few tests to cover multiple use cases.

This library is meant to be used for batching multiple items to a slower network sink.

## Features

* **Generic Type:** Use any type of data for buffering, not interface{}.
* **Size-Based Flushing:** Automatically flush data when buffer size is reached.
* **Time-Based Flushing:** Periodic flushing of data using a ticker.
* **Concurrency Safe:** Uses mutex locks to ensure safe concurrent access.

## Roadmap

No promises.

* **Non-uniform item sizes:** Useful for payload size limited flushes.
* **Error Handling & Retries:** Graceful error management with exponential backoff.
* **Parallel Processing:** Distribute tasks across multiple workers.
* **Metrics:** Provide performance-related statistics.
* **Performance Enhancements:** Optimize for high-throughput scenarios.

## Installation

```shell
go get github.com/resmoio/buff
```

## Example

```go
package main

import (
	"fmt"
	"time"
	"github.com/resmoio/buff"
)

func main() {
	flushFunc := func(data []int) {
		fmt.Println("Flushing data:", data)
	}

	buf := buff.NewBuffer(flushFunc, 3, time.Second * 2)
	buf.Start()

	buf.Push(1)
	buf.Push(2)
	buf.Push(3)  // This should immediately trigger the flushFunc due to size.

	time.Sleep(time.Second * 3)

	buf.Push(4)
	buf.Close()  // This will flush remaining data.
}
```