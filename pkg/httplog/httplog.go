package httplog

import (
	"log"
	"net/http"
)

func Log(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h(writer, request)
		log.Default().Printf("%s %s", request.Method, request.URL.Path)
	}
}
