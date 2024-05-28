package main

import (
	"sync"
	"time"

	"github.com/madeindra/stock-grpc/data"
)

func main() {
	var mtx sync.Mutex

	data.InitUpdateStocks(&mtx)
	time.Sleep(3 * time.Second)
	data.ToggleStock("AAPL", false, &mtx)
	time.Sleep(3 * time.Second)
	data.ToggleStock("AAPL", true, &mtx)
	time.Sleep(6 * time.Second)
}
