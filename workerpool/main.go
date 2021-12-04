package workerpool

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

var (
	maxMessageSize = 8 * 1024
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

// func addMessageHandler(ctx context.Context, connection Conn, userID string) Conn {
// 	connection.OnMessage(func(message []byte) error {
// 		var baseMsg MessageBase
// 		if err := json.Unmarshal(message, &baseMsg); err != nil {
// 			fmt.Println("invalid message")
// 			return err
// 		}

// 		switch baseMsg.Type {

// 		}

// 		return nil
// 	})

// 	return connection
// }

func main() {
	flag.Parse()

	http.HandleFunc("/ping", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "pong")
	})

	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("user-id")
		if userID != "" {
			fmt.Fprintf(rw, "invalid user id")
			return
		}

		log.Printf("user %s registered", userID)
	})

	log.Println("server started at", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Serve", err)
	}
}
