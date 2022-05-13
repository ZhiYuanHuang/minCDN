package cmd

import (
	"bufio"
	"context"
	"net/http"
)

type ObjectLayer interface {
	NewNSLock(recUri string) RWLocker
	GetResInfo(ctx context.Context, bucket, object, resUri string, h http.Header, lockType LockType, clientObjInfo ResObjectInfo) (resBuffer *bufio.Reader, statusCode int, err error)
}

type LockType int

const (
	noLock LockType = iota
	readLock
	writeLock
)
