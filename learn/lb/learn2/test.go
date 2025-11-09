package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type LoadBalancer struct {
	// Add fields if needed
}

// TODO: Implement ServeHTTP method here
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Log remote address
	fmt.Printf("Received request from %s\n", r.RemoteAddr)

	// 2. Log method and path
	fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)

	// 3. Loop through and log ALL headers
	for key, values := range r.Header {
		// How do you print each header?
		// Format: "Header-Name: value1, value2"
		fmt.Printf("%s: %s\n", key, strings.Join(values, ", "))
	}

	// 4. Send response back
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received by load balancer\n"))
}

func main() {
	lb := &LoadBalancer{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: lb,
	}

	log.Println("Load balancer starting on :8080")

	// What method do you call on 'server'?
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err) // or handle appropriately
	}
}
