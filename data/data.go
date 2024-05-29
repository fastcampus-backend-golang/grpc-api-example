package data

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var mtx sync.Mutex

// stockConfigs is a map that contains the stock code as the value whether the stock is enabled or not
var stockConfigs = map[string]bool{
	"AAPL": true,
	"AMZN": true,
	"GOOG": true,
	"META": true,
	"MSFT": true,
	"NFLX": true,
}

// stockPrices is a map that contains the price history of the stock
var stockPrices = map[string][]StockPrice{}

// StockPrice is a struct that contains the stock code, price, and timestamp
type StockPrice struct {
	Code      string
	Price     int64
	Timestamp time.Time
}

// init is a function that initializes the stock update for the first time
func init() {
	for code, isEnabled := range stockConfigs {
		if !isEnabled {
			continue
		}

		log.Printf("Stock %s is enabled", code)
		go updateStock(code)
	}
}

// ToggleStock is a function that toggles the stock to be enabled or disabled
func ToggleStock(code string, isEnabled bool) {
	// lock the mutex
	mtx.Lock()

	// check if the stock is already enabled or disabled
	if isEnabled == stockConfigs[code] {
		mtx.Unlock()
		return
	}

	log.Printf("Toggling %s to %t", code, isEnabled)

	// if the stock is not enabled and the stock is toggled to enable, trigger the updateStock
	if !stockConfigs[code] && isEnabled {
		stockConfigs[code] = true

		// unlock the mutex
		mtx.Unlock()

		// trigger the updateStock
		go updateStock(code)

		return
	}

	// if the stock is enabled and the stock is toggled to disable, set the stock to be disabled
	stockConfigs[code] = false

	// unlock the mutex
	mtx.Unlock()
}

// updateStock is a function that updates the stock price every second
func updateStock(code string) {
	for {
		// lock the mutex
		mtx.Lock()

		// if the stock is disabled, break the loop
		if !stockConfigs[code] {
			// unlock the mutex
			mtx.Unlock()

			break
		}

		// sleep for 1 second to simulate the update
		time.Sleep(1 * time.Second)

		// check if the stock exists
		current, exists := stockPrices[code]
		if !exists {
			// if the stock does not exist, add the stock with the initial price
			initialPrice := 10000
			initialTimestamp := time.Now()

			stockPrices[code] = []StockPrice{{Code: code, Price: int64(initialPrice), Timestamp: initialTimestamp}}

			log.Printf("Stock %s first time added with price %d", code, initialPrice)

			// unlock the mutex
			mtx.Unlock()

			continue
		}

		// if the stock exists, get the last price and randomize the next price
		lastItem := current[len(stockPrices[code])-1]
		price := randomizePrice(lastItem.Price)

		// append the new price to the stock
		stockPrices[code] = append(stockPrices[code], StockPrice{
			Code:      code,
			Price:     price,
			Timestamp: time.Now(),
		})

		log.Printf("Stock %s added with price %d", code, price)

		// unlock the mutex
		mtx.Unlock()
	}
}

// randomizePrice is a function that randomizes the price of the stock
func randomizePrice(price int64) int64 {
	// determine to add or subtract
	operation := rand.Intn(2)

	// determine the amount to add or subtract
	amount := rand.Int63n(100)

	// if operation is 0, add the amount to the price
	if operation == 0 {
		return price + amount
	}

	// if operation is 1, subtract the amount from the price
	return price - amount
}
