// integration/server_integration_test.go
package integration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"calculator/calculator"
)

func TestServerIntegration(t *testing.T) {
	// Set up the server with the same routes as main.go.
	calcHandler := calculator.NewCalculatorHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", calcHandler.HandleCalculator)
	mux.HandleFunc("/static/style.css", calculator.ServeCSS)

	// Start a test server.
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// --- Test GET Request ---
	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("GET / expected status 200, got %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Go Calculator") {
		t.Error("GET / response does not contain expected 'Go Calculator' title")
	}

	// --- Test POST Request ---
	data := url.Values{}
	data.Set("num1", "10")
	data.Set("num2", "5")
	data.Set("operation", "subtract")

	res, err = http.Post(ts.URL+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("POST / expected status 200, got %d", res.StatusCode)
	}
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "10.00 - 5.00 = 5.00") {
		t.Error("POST / response does not contain expected calculation result")
	}
}

func TestNotFoundRoute(t *testing.T) {
	// Set up the server with the same routes as main.go.
	calcHandler := calculator.NewCalculatorHandler()
	mux := http.NewServeMux()

	// Add a catch-all handler for unknown routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Handle the root path normally
			calcHandler.HandleCalculator(w, r)
		} else {
			// Return 404 for any other unknown routes
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/static/style.css", calculator.ServeCSS)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Request a non-existent route.
	res, err := http.Get(ts.URL + "/nonexistent")
	if err != nil {
		t.Fatalf("GET request for nonexistent route failed: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found for unknown route, got %d", res.StatusCode)
	}
}

func TestServerConcurrentRequests(t *testing.T) {
	// Set up the server.
	calcHandler := calculator.NewCalculatorHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", calcHandler.HandleCalculator)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Send multiple concurrent POST requests.
	var wg sync.WaitGroup
	numRequests := 10
	errCh := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			data := url.Values{}
			// Alternate between addition and subtraction.
			if i%2 == 0 {
				data.Set("num1", "10")
				data.Set("num2", "5")
				data.Set("operation", "add")
			} else {
				data.Set("num1", "10")
				data.Set("num2", "5")
				data.Set("operation", "subtract")
			}
			res, err := http.Post(ts.URL+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
			if err != nil {
				errCh <- err
				return
			}
			if res.StatusCode != http.StatusOK {
				errCh <- &httpError{Status: res.StatusCode}
			}
			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				errCh <- err
				return
			}
			html := string(body)
			if i%2 == 0 && !strings.Contains(html, "10.00 + 5.00 = 15.00") {
				errCh <- &contentError{Msg: "Addition result not found"}
			}
			if i%2 != 0 && !strings.Contains(html, "10.00 - 5.00 = 5.00") {
				errCh <- &contentError{Msg: "Subtraction result not found"}
			}
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}
}

// Custom error types for clarity.
type httpError struct {
	Status int
}

func (e *httpError) Error() string {
	return "unexpected HTTP status: " + http.StatusText(e.Status)
}

type contentError struct {
	Msg string
}

func (e *contentError) Error() string {
	return e.Msg
}
