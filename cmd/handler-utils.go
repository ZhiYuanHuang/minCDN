package cmd

import (
	"fmt"
	"net/http"
)

func errorResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Print(w)
}

func httpTraceHdrs(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.ServeHTTP(w, r)
		return
	}
}
