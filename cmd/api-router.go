package cmd

import (
	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"
)

func newHTTPServerFn() *xhttp.Server {
	globalObjLayerMutex.RLock()
	defer globalObjLayerMutex.RUnlock()
	return globalHTTPServer
}

func setHTTPServer(h *xhttp.Server) {
	globalObjLayerMutex.Lock()
	defer globalObjLayerMutex.Unlock()
	globalHTTPServer = h
}
