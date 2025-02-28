// calculator/handler_test.go
package calculator

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		num1Str   string
		num2Str   string
		operation string
		want      string
		wantErr   bool
	}{
		{"2", "3", "add", "2.00 + 3.00 = 5.00", false},
		{"2", "3", "subtract", "2.00 - 3.00 = -1.00", false},
		{"2", "3", "multiply", "2.00 × 3.00 = 6.00", false},
		{"6", "3", "divide", "6.00 ÷ 3.00 = 2.00", false},
		{"6", "0", "divide", "", true},
		{"a", "3", "add", "", true},
		{"2", "b", "add", "", true},
		{"2", "3", "modulus", "", true},
	}

	for _, tt := range tests {
		got, err := Calculate(tt.num1Str, tt.num2Str, tt.operation)
		if (err != nil) != tt.wantErr {
			t.Errorf("Calculate(%q, %q, %q) error = %v, wantErr %v", tt.num1Str, tt.num2Str, tt.operation, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("Calculate(%q, %q, %q) = %q, want %q", tt.num1Str, tt.num2Str, tt.operation, got, tt.want)
		}
	}
}

func TestServeCSS(t *testing.T) {
	req := httptest.NewRequest("GET", "/static/style.css", nil)
	w := httptest.NewRecorder()
	ServeCSS(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("ServeCSS returned status %d, expected 200", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "text/css" {
		t.Errorf("ServeCSS Content-Type = %q, expected \"text/css\"", ct)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != styleCSS {
		t.Errorf("ServeCSS body does not match expected CSS")
	}
}

func TestHandleCalculatorGet(t *testing.T) {
	handler := NewCalculatorHandler()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.HandleCalculator(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("HandleCalculator GET returned status %d, expected 200", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Go Calculator") {
		t.Errorf("HandleCalculator GET response does not contain 'Go Calculator'")
	}
}

func TestHandleCalculatorPost(t *testing.T) {
	handler := NewCalculatorHandler()
	form := url.Values{}
	form.Add("num1", "5")
	form.Add("num2", "10")
	form.Add("operation", "add")

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.HandleCalculator(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("HandleCalculator POST returned status %d, expected 200", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "5.00 + 10.00 = 15.00") {
		t.Errorf("HandleCalculator POST response does not contain expected result")
	}
}

func TestCalculateEmptyFields(t *testing.T) {
	// When either input is empty, Calculate should return an error.
	_, err := Calculate("", "5", "add")
	if err == nil {
		t.Error("Expected error when first number is empty")
	}

	_, err = Calculate("5", "", "add")
	if err == nil {
		t.Error("Expected error when second number is empty")
	}
}

func TestHandleCalculatorPostMissingOperation(t *testing.T) {
	// When the operation is missing, Calculate should return an error.
	handler := NewCalculatorHandler()
	form := url.Values{}
	form.Add("num1", "10")
	form.Add("num2", "5")
	// Note: "operation" is not set

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.HandleCalculator(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("HandleCalculator POST with missing operation returned status %d, expected 200", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	// Expect error message for invalid operation
	if !strings.Contains(string(body), "invalid operation") {
		t.Error("Expected error message for missing/invalid operation not found in response")
	}
}

func TestTemplateOptionSelection(t *testing.T) {
	// Verify that the template correctly marks the selected option.
	handler := NewCalculatorHandler()
	data := PageData{
		Title:     "Test Calculator",
		Operation: "multiply",
	}
	var buf bytes.Buffer
	err := handler.tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	html := buf.String()
	// The template should mark the "multiply" option as selected.
	if !strings.Contains(html, `<option value="multiply" selected>`) {
		t.Error("Expected multiply option to be marked as selected in the rendered template")
	}
}
