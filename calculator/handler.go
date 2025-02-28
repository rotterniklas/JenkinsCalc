// handler.go
package calculator

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// PageData represents the data passed to the HTML template
type PageData struct {
	Title     string
	Result    string
	FirstNum  string
	SecondNum string
	Operation string
}

// CalculatorHandler handles calculator functionality
type CalculatorHandler struct {
	tmpl *template.Template
}

// NewCalculatorHandler creates a new calculator handler
func NewCalculatorHandler() *CalculatorHandler {
	tmpl := template.Must(
		template.New("calculator").Funcs(template.FuncMap{
			"safe": func(s string) template.HTML {
				return template.HTML(s) // Marks string as safe HTML
			},
		}).Parse(htmlTemplate),
	)
	return &CalculatorHandler{tmpl: tmpl}
}

// HandleCalculator processes calculator requests
func (h *CalculatorHandler) HandleCalculator(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Go Calculator",
	}

	if r.Method == http.MethodPost {
		// Parse form values
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Get form values
		num1Str := r.FormValue("num1")
		num2Str := r.FormValue("num2")
		operation := r.FormValue("operation")

		// Save the input values for display
		data.FirstNum = num1Str
		data.SecondNum = num2Str
		data.Operation = operation

		// Perform calculation
		result, err := Calculate(num1Str, num2Str, operation)
		if err != nil {
			data.Result = fmt.Sprintf("Error: %s", err.Error())
		} else {
			data.Result = result
		}
	}

	// Render the template
	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Calculate performs the calculation based on the provided values and operation
// Exported for testing
func Calculate(num1Str, num2Str, operation string) (string, error) {
	// Convert string values to float64
	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)

	// Check for conversion errors
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("please enter valid numbers")
	}

	// Perform the calculation based on the selected operation
	switch operation {
	case "add":
		return fmt.Sprintf("%.2f + %.2f = %.2f", num1, num2, num1+num2), nil
	case "subtract":
		return fmt.Sprintf("%.2f - %.2f = %.2f", num1, num2, num1-num2), nil
	case "multiply":
		return fmt.Sprintf("%.2f × %.2f = %.2f", num1, num2, num1*num2), nil
	case "divide":
		if num2 == 0 {
			return "", fmt.Errorf("cannot divide by zero")
		}
		return fmt.Sprintf("%.2f ÷ %.2f = %.2f", num1, num2, num1/num2), nil
	default:
		return "", fmt.Errorf("invalid operation")
	}
}

// ServeCSS serves the CSS file
func ServeCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(styleCSS))
}
