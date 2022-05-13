package cmd

import (
	"os"
	"time"

	"github.com/dustin/go-humanize"

	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"
	"github.com/coocood/freecache"
)

const (
	GlobalMinCDNDefaultPort = "9006"

	GlobalImgeContentTypePrefix  = "image/"
	GlobalDefaultImgeContentType = "image/jpeg"

	GlobalMemoryCacheSize = 1 * humanize.GiByte
)

var (
	globalOSSignalCh = make(chan os.Signal, 1)

	globalMinCDNAddr = ""

	globalMinCDNPort = GlobalMinCDNDefaultPort

	globalHTTPServer        *xhttp.Server
	globalHTTPServerErrorCh = make(chan error)
	globalMemoryCache       *freecache.Cache

	globalUseETCD bool

	globalETCDAddress = ""

	globalStaticRecsourceMap = map[string]bool{
		"designer": true,
		"stranger": true,
	}

	globalOperationTimeout = newDynamicTimeout(10*time.Minute, 5*time.Minute)

	globalMinioEndPoint = ""

	globalMinioAccessKeyID = ""

	globalMinioSecretAccessID = ""

	globalMemoryCacheSize = GlobalMemoryCacheSize
)
