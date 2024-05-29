package main

import (
	"time"

	"github.com/madeindra/stock-grpc/data"
)

func main() {
	time.Sleep(3 * time.Second)
	data.ToggleStock("AAPL", false)
	time.Sleep(3 * time.Second)
	data.ToggleStock("AAPL", true)
	time.Sleep(6 * time.Second)
}
