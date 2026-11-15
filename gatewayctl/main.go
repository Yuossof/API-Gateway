package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "activate" {
		fmt.Println("usage: gatewayctl activate <version>")
		os.Exit(1)
	}

	requestBody := map[string]string{
		"version": os.Args[2],
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		DialContext: func(
			ctx context.Context,
			_, _ string,
		) (net.Conn, error) {
			return net.Dial("unix", "/tmp/gateway.sock")
		},
	}

	client := &http.Client{Transport: transport}

	request, err := http.NewRequest(
		http.MethodPost,
		"http://unix/activate",
		bytes.NewReader(body),
	)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	io.Copy(os.Stdout, response.Body)
}
