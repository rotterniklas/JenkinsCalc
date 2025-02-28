// main.go
package main

import (
	"calculator/calculator"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize the calculator handler
	calcHandler := calculator.NewCalculatorHandler()

	// Register routes
	http.HandleFunc("/", calcHandler.HandleCalculator)
	http.HandleFunc("/static/style.css", calculator.ServeCSS)

	// Start the server
	port := ":8090"
	fmt.Printf("Server starting at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
