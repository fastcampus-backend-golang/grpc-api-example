package data

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var mtx sync.Mutex

// stockConfigs adalah map yang berisi kode saham sebagai kunci dan nilai boolean yang menunjukkan apakah saham tersebut diaktifkan atau tidak
var stockConfigs = map[string]bool{
	"AAPL": true,
	"AMZN": true,
	"GOOG": true,
	"META": true,
	"MSFT": true,
	"NFLX": true,
}

// stockPrices adalah map yang berisi riwayat harga saham
var stockPrices = map[string][]StockPrice{}

// StockPrice adalah struct yang berisi kode saham, harga, dan timestamp
type StockPrice struct {
	Code      string
	Price     int64
	Timestamp time.Time
}

// init adalah fungsi yang menginisialisasi pembaruan saham untuk pertama kalinya
func init() {
	for code, isEnabled := range stockConfigs {
		if !isEnabled {
			continue
		}

		log.Printf("Saham %s diaktifkan", code)
		go updateStock(code)
	}
}

// ToggleStock adalah fungsi yang mengaktifkan atau menonaktifkan saham
func ToggleStock(code string, isEnabled bool) {
	// kunci mutex
	mtx.Lock()

	// periksa apakah saham sudah diaktifkan atau dinonaktifkan
	if isEnabled == stockConfigs[code] {
		mtx.Unlock()
		return
	}

	log.Printf("Mengubah status %s menjadi %t", code, isEnabled)

	// jika saham belum diaktifkan dan diubah menjadi diaktifkan, trigger updateStock
	if !stockConfigs[code] && isEnabled {
		stockConfigs[code] = true

		// buka kunci mutex
		mtx.Unlock()

		// trigger updateStock
		go updateStock(code)

		return
	}

	// jika saham sudah diaktifkan dan diubah menjadi dinonaktifkan, set saham menjadi dinonaktifkan
	stockConfigs[code] = false

	// buka kunci mutex
	mtx.Unlock()
}

// updateStock adalah fungsi yang memperbarui harga saham setiap detik
func updateStock(code string) {
	for {
		// kunci mutex
		mtx.Lock()

		// jika saham dinonaktifkan, hentikan perulangan
		if !stockConfigs[code] {
			// buka kunci mutex
			mtx.Unlock()

			break
		}

		// tidur selama 1 detik untuk mensimulasikan pembaruan
		time.Sleep(1 * time.Second)

		// periksa apakah saham ada
		current, exists := stockPrices[code]
		if !exists {
			// jika saham tidak ada, tambahkan saham dengan harga awal
			hargaAwal := 10000
			waktuAwal := time.Now()

			stockPrices[code] = []StockPrice{{Code: code, Price: int64(hargaAwal), Timestamp: waktuAwal}}

			log.Printf("Saham %s ditambahkan untuk pertama kali dengan harga %d", code, hargaAwal)

			// buka kunci mutex
			mtx.Unlock()

			continue
		}

		// jika saham ada, dapatkan harga terakhir dan acak harga berikutnya
		itemTerakhir := current[len(stockPrices[code])-1]
		harga := randomizePrice(itemTerakhir.Price)

		// tambahkan harga baru ke saham
		stockPrices[code] = append(stockPrices[code], StockPrice{
			Code:      code,
			Price:     harga,
			Timestamp: time.Now(),
		})

		log.Printf("Saham %s ditambahkan dengan harga %d", code, harga)

		// buka kunci mutex
		mtx.Unlock()
	}
}

// randomizePrice adalah fungsi yang mengacak harga saham
func randomizePrice(harga int64) int64 {
	// tentukan apakah akan ditambahkan atau dikurangkan
	operasi := rand.Intn(2)

	// tentukan jumlah yang akan ditambahkan atau dikurangkan
	jumlah := rand.Int63n(100)

	// jika operasi adalah 0, tambahkan jumlah ke harga
	if operasi == 0 {
		return harga + jumlah
	}

	// jika operasi adalah 1, kurangkan jumlah dari harga
	return harga - jumlah
}

// GetStockConfig adalah fungsi untuk mengembalikan konfigurasi saham saat ini
func GetStockConfig() map[string]bool {
	return stockConfigs
}

func GetStockPrice(code string) []StockPrice {
	return stockPrices[code]
}
