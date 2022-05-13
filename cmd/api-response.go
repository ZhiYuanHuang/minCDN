package cmd

import (
	"net/http"
	"strconv"

	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"
)

func setCommonHeaders(w http.ResponseWriter) {
	// Set the "Server" http header.
	w.Header().Set(xhttp.ServerInfo, "MinCDN")

	w.Header().Set(xhttp.AcceptRanges, "bytes")

}

func writeErrorResponse(w http.ResponseWriter) {

	writeResponse(w, http.StatusBadRequest, nil, GlobalDefaultImgeContentType)
}

func writeResponse(w http.ResponseWriter, statusCode int, response []byte, mType mimeType) {
	setCommonHeaders(w)
	if mType != mimeNone {
		w.Header().Set(xhttp.ContentType, string(mType))
	}
	w.Header().Set(xhttp.ContentLength, strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	if response != nil {
		w.Write(response)
	}
}

type mimeType string

const (
	mimeNone mimeType = ""
	mimeJSON mimeType = "application/json"
	mimeXML  mimeType = "application/xml"
)
