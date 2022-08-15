package main

import (
	"fmt"
	"net/http"

	"github.com/mplewis/gemocities"
)

func main() {
	srv := gemocities.BuildServer()
	fmt.Println("Server listening on :8888")
	http.ListenAndServe(":8888", srv)
}
