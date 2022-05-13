package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"
)

var couter int

var etag_rec string
var lastModify_rec time.Time

func (api objectAPIHandlers) GetImageHandle(w http.ResponseWriter, r *http.Request) {

	ifModifiedSinceStr := r.Header.Get(xhttp.IfModifiedSince)
	var objInfo ResObjectInfo
	if ifModifiedSinceStr != "" {
		ifModifiedSinceTime, timeErr := time.Parse(time.RFC1123, ifModifiedSinceStr)
		if timeErr != nil {
			log.Println(timeErr)
			writeResponse(w, http.StatusBadRequest, nil, GlobalDefaultImgeContentType)
			return
		}
		objInfo.LastModify = ifModifiedSinceTime
	}

	ifNoneMatchStr := r.Header.Get(xhttp.IfNoneMatch)
	if ifNoneMatchStr != "" {
		objInfo.Etag = ifNoneMatchStr
	}

	objectAPI := api.ObjectAPI()
	if objectAPI == nil {
		writeErrorResponse(w)
		return
	}

	imgUrl := r.URL.Path
	imgUrl = strings.TrimPrefix(imgUrl, "/")
	splitedUrl := strings.Split(imgUrl, "/")

	if len(splitedUrl) != 5 {
		log.Printf("url path error,url:%s", imgUrl)
		writeResponse(w, http.StatusNotFound, nil, GlobalDefaultImgeContentType)
		return
	}

	recModule := strings.ToLower(splitedUrl[1])

	if moduleEnableState, ok := globalStaticRecsourceMap[recModule]; !ok {

		log.Printf("recModule not supported,recModule:%s", recModule)
		writeResponse(w, http.StatusNotFound, nil, GlobalDefaultImgeContentType)
		return
	} else {
		if !moduleEnableState {

			log.Printf("recModule not enable,recModule:%s", recModule)
			writeResponse(w, http.StatusNotFound, nil, GlobalDefaultImgeContentType)
			return
		}
	}

	splitedRecFile := strings.Split(splitedUrl[4], ".")
	if len(splitedRecFile) != 2 {

		log.Printf("rec file error,rec file:%s", splitedUrl[4])
		writeResponse(w, http.StatusNotFound, nil, GlobalDefaultImgeContentType)
		return
	}
	recId := 0
	if maybeRecId, err := strconv.Atoi(splitedRecFile[0]); err != nil {

		log.Printf("rec file name error,name:%s", splitedRecFile[0])
		writeResponse(w, http.StatusNotFound, nil, GlobalDefaultImgeContentType)
		return
	} else {
		recId = maybeRecId
	}

	bucketName := fmt.Sprintf("%s-%d", recModule, recId%10)
	recFileBucketPath := fmt.Sprintf("%s/%s/%s", splitedUrl[2], splitedUrl[3], splitedUrl[4])
	fmt.Printf("rec belong bucket %s\n", bucketName)

	ctx := context.Background()

	resBuffer, statusCode, err := objectAPI.GetResInfo(ctx, bucketName, recFileBucketPath, imgUrl, w.Header(), readLock, objInfo)

	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, nil, GlobalDefaultImgeContentType)
		return
	}

	var chunk []byte
	if resBuffer != nil && resBuffer.Size() > 0 {
		buf := make([]byte, 1024)
		for {
			n, err := resBuffer.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("read buf fail", err)
				return
			}
			//说明读取结束
			if n == 0 {
				break
			}
			//读取到最终的缓冲区中
			chunk = append(chunk, buf[:n]...)
		}
	}

	writeResponse(w, statusCode, chunk, mimeNone)

	// return

	// endpoint := "172.30.237.255:9000"
	// //accessKeyID := "zion"
	// //secretAccessKey := "hzy123456"

	// endpoint = "172.24.205.20:9000"
	// accessKeyID := "minioadmin"
	// secretAccessKey := "minioadmin"
	// useSSL := false

	// minioClient, err := minio.New(endpoint, &minio.Options{
	// 	Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	// 	Secure: useSSL,
	// })

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// //object, err := minioClient.GetObject(ctx, "test", "vanGogh/paints/606px-Van_Gogh_-_Starry_Night.jpg", minio.GetObjectOptions{})
	// object, err := minioClient.GetObject(ctx, "test", "门配置/门报警配置2.jpg", minio.GetObjectOptions{})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// tmpBuffer, err2 := ioutil.ReadAll(object)
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }

	// objInfo, _ := object.Stat()

	// etag_rec = objInfo.ETag
	// lastModify_rec = objInfo.LastModified

	// w.Header().Set(xhttp.ETag, objInfo.ETag)
	// w.Header().Set(xhttp.LastModified, objInfo.LastModified.Format(time.RFC1123))
	// w.Header().Set(xhttp.CacheControl, fmt.Sprintf("private, max-age=%d", math.MaxUint32))
	// //w.Header().Set(xhttp,)
	// //bw := bufio.NewWriter(w)
	// //io.Copy(bw, object)

	// // inlineBytes := make([]byte, 10)
	// // buffer := bytes.NewBuffer(inlineBytes)
	// // io.Copy(buffer, object)

	// // io.Copy(w, buffer)

	// writeResponse(w, http.StatusOK, tmpBuffer, GlobalDefaultImgeContentType)

	// // for object := range objectCh {
	// // 	if object.Err != nil {
	// // 		fmt.Println(object.Err)
	// // 		return
	// // 	}
	// // 	fmt.Println(object)
	// // }

	// // w.Header().Set("Content-Type", "application/json")
	// // w.WriteHeader(http.StatusOK)
	// // var obj struct {
	// // 	Ss string
	// // 	Id int
	// // }
	// // obj.Ss = "dsfdsf"
	// // couter++
	// // obj.Id = recId
	// // json.NewEncoder(w).Encode(obj)
}
