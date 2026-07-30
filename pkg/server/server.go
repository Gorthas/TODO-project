package server

import "net/http"

func Run(addr string) error {
	webDir := http.Dir("./web")
	handler := http.FileServer(webDir)
	http.Handle("/", handler)

	return http.ListenAndServe(addr, nil)
}
