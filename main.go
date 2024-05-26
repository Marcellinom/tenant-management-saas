package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Halo")
	})
	log.Fatal(http.ListenAndServe(os.Getenv("APP_URL"), nil))
}
