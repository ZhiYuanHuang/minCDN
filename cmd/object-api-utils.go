package cmd

import (
	"io"
	"log"
	"time"

	"github.com/google/uuid"
)

const SlashSeparator = "/"

func mustGetUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

type GetObjectReader struct {
	io.Reader
	ObjInfo ResObjectInfo
}

type ResObject struct {
	ByteData []byte
	ObjInfo  ResObjectInfo
}

type ResObjectInfo struct {
	Etag        string
	LastModify  time.Time
	ContentType string
}
