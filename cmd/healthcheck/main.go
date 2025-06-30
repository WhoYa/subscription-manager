package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	url := fmt.Sprintf("http://127.0.0.1:%s/healthz", port)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stdout, "GET error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stdout, "Status %d: %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	os.Exit(0)
}
