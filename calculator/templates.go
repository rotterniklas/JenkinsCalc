// templates.go
package calculator

// HTML template for the calculator page
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <h1>Go Calculator Test1</h1>
    <div class="calculator">
        <form method="POST" action="/">
            <input type="number" name="num1" placeholder="Enter first number" value="{{.FirstNum}}" required>
            <select name="operation">
                <option value="add" {{if eq .Operation "add"}}selected{{end}}>Addition (+)</option>
                <option value="subtract" {{if eq .Operation "subtract"}}selected{{end}}>Subtraction (-)</option>
                <option value="multiply" {{if eq .Operation "multiply"}}selected{{end}}>Multiplication (×)</option>
                <option value="divide" {{if eq .Operation "divide"}}selected{{end}}>Division (÷)</option>
            </select>
            <input type="number" name="num2" placeholder="Enter second number" value="{{.SecondNum}}" required>
            <button type="submit">Calculate</button>
        </form>
		{{if .Result}}
		<div class="result">
			Result: {{.Result | safe}} <!-- Use | safe here -->
		</div>
		{{end}}
    </div>
</body>
</html>
`

// CSS styles for the calculator page
const styleCSS = `
body {
    font-family: Arial, sans-serif;
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
}

h1 {
    color: #2c3e50;
    text-align: center;
}

.calculator {
    background-color: #f7f9fa;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

form {
    display: grid;
    gap: 15px;
}

input, select, button {
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
}

button {
    background-color: #3498db;
    color: white;
    border: none;
    cursor: pointer;
    transition: background-color 0.3s;
}

button:hover {
    background-color: #2980b9;
}

.result {
    margin-top: 20px;
    padding: 15px;
    background-color: #e8f4fc;
    border-radius: 4px;
    text-align: center;
    font-size: 18px;
}
`
