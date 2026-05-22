package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Car struct {
	Make  string
	Model string
	Price string
}

var (
	appVersion = flag.Int("version", 2, "App version")
)

func main() {
	flag.Parse()

	cars := []Car{
		{Make: "Toyota", Model: "Camry", Price: "$25,000"},
		{Make: "Honda", Model: "Civic", Price: "$22,000"},
		{Make: "Tesla", Model: "Model 3", Price: "$40,000"},
		{Make: "Ford", Model: "Mustang", Price: "$35,000"},
	}

	tmpl := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Car Catalog</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        .car { border: 1px solid #ccc; padding: 15px; margin-bottom: 10px; border-radius: 5px; max-width: 400px; }
        .buy-btn { background-color: #28a745; color: white; padding: 8px 12px; text-decoration: none; border-radius: 4px; display: inline-block; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Car Catalog</h1>
    {{range .Cars}}
    <div class="car">
        <h3>{{.Make}} {{.Model}}</h3>
        <p>Price: {{.Price}}</p>
        {{if eq $.Version 2}}
        <a href="buy?model={{.Model}}" class="buy-btn">Buy Now</a>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`))

	// Handle the root. Traefik sends /user01 here as "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Safety check: ensure we only handle the root path
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data := struct {
			Version int
			Cars    []Car
		}{
			Version: *appVersion,
			Cars:    cars,
		}
		tmpl.Execute(w, data)
	})

	// Handle the buy route. Traefik sends /user01/buy here as "/buy"
	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		model := r.URL.Query().Get("model")
		// Use relative path for the "Back" link too
		fmt.Fprintf(w, "<h1>Successful bought!</h1><p>You bought a %s.</p><a href='./'>Back</a>", model)
	})

	fmt.Printf("Starting server on :8080 with version %d\n", *appVersion)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
