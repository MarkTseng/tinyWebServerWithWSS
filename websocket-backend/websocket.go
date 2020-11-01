package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"time"
)

func main() {
	supportedSubProtocols := []string{"test-protocol"}
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: supportedSubProtocols,
	}

	// backend echo server
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Don't upgrade if original host header isn't preserved
		/*
			if r.Host != "127.0.0.1:7777" {
				fmt.Printf("Host header set incorrectly.  Expecting 127.0.0.1:7777 got %s", r.Host)
				return
			}
		*/
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			messageType, p, err := conn.ReadMessage()
			fmt.Println(messageType)
			fmt.Printf("%s\n", p)
			if err != nil {
				fmt.Printf("read err\n")
				return
			}

			if err = conn.WriteMessage(messageType, p); err != nil {
				fmt.Printf("write err\n")
				return
			}
		}
	})

	// wss:
	//err := http.ListenAndServeTLS("192.168.1.20:7777", "../ssldata/certificate.crt", "../ssldata/private.key", mux2)
	err := http.ListenAndServeTLS("127.0.0.1:8888", "../ssldata/certificate.crt", "../ssldata/private.key", mux2)
	// ws
	//err := http.ListenAndServe(":8888", mux2)
	if err != nil {
		fmt.Println("ListenAndServe or ListenAndServeTLS ", err)
	}

}
