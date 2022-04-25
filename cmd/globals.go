package cmd

import (
	"os"

	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"
)

const (
	GlobalMinCDNDefaultPort = "9006"
)

var (
	globalOSSignalCh = make(chan os.Signal, 1)

	globalHTTPServer *xhttp.Server
)
