package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	// "bufio"
	// "bytes"
	// "context"
	// "encoding/json"
	// "fmt"
	// "io"
	// "net/http"
	// "os"
	// "time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1. Receive data (GET method)")
		fmt.Println("2. Send data (POST method)")
		fmt.Printf(">> ")
		scanner.Scan()
		opt := scanner.Text()
		client := &http.Client{
			Timeout:   3 * time.Second,
			Transport: &http.Transport{DisableKeepAlives: true},
		}
		http.DefaultClient = client
		switch opt {
		case "1":
			receiveData()
		case "2":
			sendData()
		}
	}
}

func loadTLSConfig() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	caCert, err := os.ReadFile("../cert.pem")
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs: certPool,
	}, nil
}

func receiveData() {
	tlsConfig, err := loadTLSConfig()
	if err != nil {
		return
	}
	cx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequestWithContext(cx, "GET", "https://localhost:9876", nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("Data from server: ", string(data))
}

func sendData() {
	data := map[string]string{
		"Name": "Test",
		"Age":  "10",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	tlsConfig, err := loadTLSConfig()
	if err != nil {
		return
	}
	cx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequestWithContext(cx, "POST", "https://localhost:9876/post", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("Server response: ", string(body))
}
