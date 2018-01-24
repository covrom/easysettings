package main

import (
	"net/http"
)

func main() {
	App.Init()

	http.ListenAndServe(App.Host+":"+App.Port, nil)
}
