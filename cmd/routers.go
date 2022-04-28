package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func configureServerHandler() (http.Handler, error) {
	router := mux.NewRouter().SkipClean(true).UseEncodedPath()

	registerAPIRouter(router)

	return router, nil
}
