package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/coocood/freecache"

	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOObjects struct {
	nsMutex     *nsLockMap
	minioClient *minio.Client
	memoryCache *freecache.Cache
}

func NewMinIOObjectLayer(endPoint, accessKeyID, secretAccessKey string, useSSL bool) (ObjectLayer, error) {

	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	fs := &MinIOObjects{
		nsMutex:     newNSLock(),
		minioClient: minioClient,
		memoryCache: newMemoryCacheFn(),
	}

	return fs, nil
}

func (fs *MinIOObjects) NewNSLock(recUri string) RWLocker {
	return fs.nsMutex.NewNSLock(recUri)
}

func (fs *MinIOObjects) GetResInfo(ctx context.Context, bucket, object, resUri string, h http.Header, lockType LockType, clientObjInfo ResObjectInfo) (resBuffer *bufio.Reader, statusCode int, err error) {

	if ok, cacheResObj := fs.getCache(resUri); ok {
		if !clientObjInfo.LastModify.IsZero() {
			if cacheResObj.ObjInfo.LastModify.Equal(clientObjInfo.LastModify) {
				if clientObjInfo.Etag == "" {
					statusCode = http.StatusNotModified
					return
				} else {
					if cacheResObj.ObjInfo.Etag == clientObjInfo.Etag {
						statusCode = http.StatusNotModified
						return
					}
				}
			}
		} else {
			if clientObjInfo.Etag != "" && cacheResObj.ObjInfo.Etag == clientObjInfo.Etag {
				statusCode = http.StatusNotModified
				return
			}
		}
	}

	nsUnLocker := func() {}
	if lockType != noLock {
		lock := fs.NewNSLock(resUri)
		switch lockType {
		case writeLock:
			lkctx, err := lock.GetLock(ctx, globalOperationTimeout)
			if err != nil {
				return nil, http.StatusBadGateway, err
			}
			ctx = lkctx.Context()
			nsUnLocker = func() { lock.Unlock(lkctx.cancel) }
		case readLock:
			lkctx, err := lock.GetRLock(ctx, globalOperationTimeout)
			if err != nil {
				return nil, http.StatusBadGateway, err
			}
			ctx = lkctx.Context()
			nsUnLocker = func() { lock.RUnlock(lkctx.Cancel) }
		}
	}

	if ok, cacheResObj := fs.getCache(resUri); ok {
		if !clientObjInfo.LastModify.IsZero() {
			if cacheResObj.ObjInfo.LastModify.Equal(clientObjInfo.LastModify) {
				if clientObjInfo.Etag == "" {
					statusCode = http.StatusNotModified
					nsUnLocker()
					return
				} else {
					if cacheResObj.ObjInfo.Etag == clientObjInfo.Etag {
						statusCode = http.StatusNotModified
						nsUnLocker()
						return
					}
				}
			}
		} else {
			if clientObjInfo.Etag != "" && cacheResObj.ObjInfo.Etag == clientObjInfo.Etag {
				statusCode = http.StatusNotModified
				nsUnLocker()
				return
			}
		}
	}

	var minioObject *minio.Object
	minioObject, err = fs.minioClient.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		nsUnLocker()
		log.Println(err)
		return nil, http.StatusNotFound, err
	}

	var minioObjInfo minio.ObjectInfo
	minioObjInfo, err = minioObject.Stat()
	if err != nil {
		nsUnLocker()
		log.Println(err)
		return nil, http.StatusBadGateway, err
	}

	h.Set(xhttp.ETag, minioObjInfo.ETag)
	h.Set(xhttp.LastModified, minioObjInfo.LastModified.Format(time.RFC1123))
	h.Set(xhttp.CacheControl, fmt.Sprintf("private, max-age=%d", math.MaxUint32))
	h.Set(xhttp.ContentType, minioObjInfo.ContentType)

	buffer := bytes.NewBuffer([]byte{})

	_, err = io.Copy(buffer, minioObject)
	if err != nil {
		nsUnLocker()
		log.Println(err)
		return nil, http.StatusBadGateway, err
	}

	resObj := &ResObject{
		ObjInfo: ResObjectInfo{
			Etag:        minioObjInfo.ETag,
			LastModify:  minioObjInfo.LastModified,
			ContentType: minioObjInfo.ContentType,
		},
	}

	inlineBytes := make([]byte, buffer.Len())
	buffer.Read(inlineBytes)

	go fs.setCache(resUri, resObj, inlineBytes)

	respBuff := bytes.NewBuffer(inlineBytes)
	resBuffer = bufio.NewReader(respBuff)
	nsUnLocker()
	return resBuffer, http.StatusOK, nil
}

func (fs *MinIOObjects) getCache(resUri string) (ok bool, resObj *ResObject) {
	cacheKeyBytes := []byte(resUri)
	cacheValueBytes, err := fs.memoryCache.Get(cacheKeyBytes)
	if err != nil {
		ok = false
		_ = cacheValueBytes
		//log.Println(err)
		return
	}

	if len(cacheValueBytes) <= 0 {
		ok = false
		log.Printf("uri:%s,cache hit but len is %d", resUri, len(cacheValueBytes))
		return
	}

	cacheValue := bytes.NewBuffer(cacheValueBytes)
	dec := gob.NewDecoder(cacheValue)
	err = dec.Decode(&resObj)
	if err != nil {
		ok = false
		log.Println(err)
		return
	}

	ok = true

	return
}

func (fs *MinIOObjects) setCache(resUri string, resObj *ResObject, originByteData []byte) {
	tmpChan := make(chan []byte)

	go func() {
		//fmt.Println("set cache running")
		cacheBytes := make([]byte, len(originByteData))
		copy(cacheBytes, originByteData)
		tmpChan <- cacheBytes
	}()

	cacheKeyBytes := []byte(resUri)

	//tmpByteData := <-tmpChan
	resObj.ByteData = <-tmpChan
	close(tmpChan)

	cacheValue := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(cacheValue)
	err := enc.Encode(resObj)

	if err != nil {
		log.Println(err)
		return
	}

	cacheValueBytes := make([]byte, cacheValue.Len())
	cacheValue.Read(cacheValueBytes)

	expire := 60 * 3
	err = fs.memoryCache.Set(cacheKeyBytes, cacheValueBytes, expire)
	if err != nil {
		log.Println(err)
	}
}
