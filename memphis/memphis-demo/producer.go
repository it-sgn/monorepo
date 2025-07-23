package main

import (
	"fmt"
	"os"
	"os/signal" // Import untuk menangani sinyal sistem
	"syscall"   // Import untuk sinyal sistem
	"time"      // Import untuk time.Sleep

	"github.com/memphisdev/memphis.go"
)

func main() {
	// Konfigurasi koneksi Memphis
	memphisHost := "localhost"
	memphisUser := "root"
	memphisPassword := "memphis"
	stationName := "Attendance-System" // Nama station
	producerName := "att-log"          // Nama producer

	// Menghubungkan ke Memphis
	conn, err := memphis.Connect(memphisHost, memphisUser, memphis.Password(memphisPassword))
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close() // Pastikan koneksi ditutup saat program berakhir

	// Membuat Producer
	p, err := conn.CreateProducer(stationName, producerName)
	if err != nil {
		fmt.Printf("Create Producer failed: %v\n", err)
		os.Exit(1)
	}

	// Menyiapkan Headers (opsional, bisa dihapus jika tidak diperlukan)
	hdrs := memphis.Headers{}
	hdrs.New()
	err = hdrs.Add("Coba", "123")
	if err != nil {
		fmt.Printf("Header failed: %v\n", err)
		os.Exit(1)
	}

	// Membuat channel untuk menangani sinyal penghentian (Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Mendaftarkan sinyal SIGINT dan SIGTERM

	// Loop tak terbatas untuk mengirim pesan
	messageCounter := 0
	fmt.Println("Producer running... Press Ctrl+C to stop.")

	for {
		select {
		case <-sigs:
			// Menerima sinyal penghentian, keluar dari loop
			fmt.Println("\nReceived stop signal. Shutting down producer...")
			return // Keluar dari fungsi main
		default:
			// Mengirim pesan secara berkala
			messageCounter++
			message := fmt.Sprintf("You have a message! (Message #%d at %s)", messageCounter, time.Now().Format(time.RFC3339))
			err = p.Produce([]byte(message), memphis.MsgHeaders(hdrs))
			if err != nil {
				fmt.Printf("Produce failed: %v\n", err)
				// Lanjutkan meskipun ada error, atau tambahkan logika retry/exit
			} else {
				fmt.Printf("Produced message: \"%s\"\n", message)
			}

			// Jeda sebentar sebelum mengirim pesan berikutnya
			time.Sleep(2 * time.Second) // Mengirim pesan setiap 2 detik
		}
	}
}

// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/memphisdev/memphis.go"
// )

// func main() {
// 	conn, err := memphis.Connect("localhost", "root", memphis.Password("memphis"))
// 	if err != nil {
// 		fmt.Printf("Connection failed: %v", err)
// 		os.Exit(1)
// 	}
// 	defer conn.Close()
// 	p, err := conn.CreateProducer("Attendance-System", "att-log")

// 	if err != nil {
// 		fmt.Printf("Create Producer failed: %v", err)
// 		os.Exit(1)
// 	}

// 	hdrs := memphis.Headers{}
// 	hdrs.New()
// 	err = hdrs.Add("Coba", "123")

// 	if err != nil {
// 		fmt.Printf("Header failed: %v", err)
// 		os.Exit(1)
// 	}

// 	err = p.Produce([]byte("You have a message!"), memphis.MsgHeaders(hdrs))

// 	if err != nil {
// 		fmt.Printf("Produce failed: %v", err)
// 		os.Exit(1)
// 	}
// }
