package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(wr http.ResponseWriter, rq *http.Request) {
		fmt.Fprint(wr, "This is server Mux")
	})
	mux.HandleFunc("/post", func(wr http.ResponseWriter, rq *http.Request) {
		data, err := io.ReadAll(rq.Body)
		if err != nil {
			return
		}
		var requestData map[string]string
		err = json.Unmarshal(data, &requestData)
		if err != nil {
			return
		}
		fmt.Println("Data from client: ", requestData)
		fmt.Fprint(wr, "Data successfully received by server")
	})
	server := http.Server{
		Addr:    "localhost:9876",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	err := server.ListenAndServeTLS("../cert.pem", "../key.pem")
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		fmt.Println(err)
		return
	}
}
