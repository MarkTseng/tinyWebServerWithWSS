package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	serverURL  = "wss://192.168.1.20:7777"
	backendURL = "wss://127.0.0.1:8888"
)

func main() {
	// CA
	rootCert, err := ioutil.ReadFile("../ssldata/ca_bundle.crt")
	if err != nil {
		fmt.Println(err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCert) {
		fmt.Println("no root certificate detected")
	}

	// websocket proxy
	supportedSubProtocols := []string{"test-protocol"}
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: supportedSubProtocols,
	}

	u, _ := url.Parse(backendURL)
	proxy := websocketproxy.NewProxy(u)
	proxy.Upgrader = upgrader
	// setup tls
	proxy.Dialer = websocket.DefaultDialer
	proxy.Dialer.TLSClientConfig = &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxy)
	if err := http.ListenAndServeTLS(":7777", "../ssldata/certificate.crt", "../ssldata/private.key", mux); err != nil {
		fmt.Printf("ListenAndServe: %s", err)
	}

}
