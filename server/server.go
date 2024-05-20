package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
