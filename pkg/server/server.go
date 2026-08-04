package server

import (
	"net/http"

	"github.com/Gorthas/TODO-project/pkg/api"
)

func Run(addr string) error {
	api.Init()

	webDir := http.Dir("./web")
	handler := http.FileServer(webDir)
	http.Handle("/", handler)

	return http.ListenAndServe(addr, nil)
}
