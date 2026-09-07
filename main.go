package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Ascii_Art(input string) string {
	data, err := os.ReadFile("standard.txt")
	if err != nil {
		fmt.Println("ERROR", err)
		return ""
	}

	lines := strings.Split(string(data), "\n")
	font := make(map[rune][]string)

	for ch := ' '; ch <= '~'; ch++ {
		start := (int(ch) - 32) * 9
		font[ch] = lines[start+1 : start+9]
	}

	inputfile := strings.ReplaceAll(input, "\\n", "\n")
	word := strings.Split(inputfile, "\n")

	var result strings.Builder
	for i, words := range word {
		if words == "" {
			if i != len(word)-1 {
				result.WriteString("\n")
			}
			continue
		}
		for row := 0; row < 8; row++ {
			for _, ch := range words {
				result.WriteString(font[ch][row])
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `
	<html>
		<body>
			<h1>ASCII Art Generator</h1>
			<form action="/submit" method="POST">
				<input type="text" name="userText" placeholder="Enter text...">
				<label for="color">Color:</label>
				<input type="color" id="color" name="color" value="#000000">
				<button type="submit">Generate!</button>
			</form>
		</body>
	</html>`
	fmt.Fprint(w, html)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	userText := r.FormValue("userText")
	color := r.FormValue("color")
	result := Ascii_Art(userText)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	<html>
		<body>
			<pre style="color:%s">%s</pre>
			<a href="/">Go back</a>
		</body>
	</html>`, color, result)
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)

	log.Println("Server running on http://localhost:4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}
