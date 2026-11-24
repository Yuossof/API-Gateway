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
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "activate" {
		fmt.Println("usage: gatewayctl activate <path-to-config.yaml>")
		os.Exit(1)
	}

	// تحويل المسار لمسار مطلق عشان يشتغل صح بغض النظر عن working directory بتاع الـ daemon
	absPath, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	requestBody := map[string]string{
		"path": absPath,
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
