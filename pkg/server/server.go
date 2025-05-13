package server

import (
	"fmt"
	"net/http"
)

const (
	webDir = "web"
	port   = 7540
)

func Run() error {

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
