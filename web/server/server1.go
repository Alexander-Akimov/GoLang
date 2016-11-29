package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	}

	http.HandleFunc("/", handler1) // Каждый запрос вызывет обработчик

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Обработчик возврщает компонент пути из URL запроса.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(os.Stderr, "URL.Path = %q\n", r.URL.Path)
}
