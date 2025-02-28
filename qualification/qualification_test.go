// qualification/qualification_test.go
package qualification

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"calculator/calculator"
)

func TestCalculatorQualification(t *testing.T) {
	// This test simulates a full user workflow: entering numbers, selecting an operation,
	// and verifying that the resulting page contains the correct output.

	// Set up the handler and server.
	calcHandler := calculator.NewCalculatorHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", calcHandler.HandleCalculator)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Simulate a POST request as if a user submitted the form.
	form := url.Values{}
	form.Set("num1", "20")
	form.Set("num2", "4")
	form.Set("operation", "divide")

	res, err := http.Post(ts.URL+"/", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Qualification test POST failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Qualification test expected status 200, got %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	// Check that the page output includes the expected result.
	if !strings.Contains(string(body), "20.00 ÷ 4.00 = 5.00") {
		t.Error("Qualification test failed: expected result not found in the response")
	}
}

func TestCalculatorMultipleOperations(t *testing.T) {
	// Simulate a full workflow with multiple operations.
	calcHandler := calculator.NewCalculatorHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", calcHandler.HandleCalculator)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	operations := []struct {
		num1      string
		num2      string
		operation string
		expected  string
	}{
		{"3", "2", "add", "3.00 + 2.00 = 5.00"},
		{"3", "2", "subtract", "3.00 - 2.00 = 1.00"},
		{"3", "2", "multiply", "3.00 × 2.00 = 6.00"},
		{"6", "3", "divide", "6.00 ÷ 3.00 = 2.00"},
	}

	for _, op := range operations {
		form := url.Values{}
		form.Set("num1", op.num1)
		form.Set("num2", op.num2)
		form.Set("operation", op.operation)

		res, err := http.Post(ts.URL+"/", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("POST request failed for operation %s: %v", op.operation, err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for operation %s, got %d", op.operation, res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read response body for operation %s: %v", op.operation, err)
		}
		if !strings.Contains(string(body), op.expected) {
			t.Errorf("For operation %s, expected result %q not found in response", op.operation, op.expected)
		}
	}
}
