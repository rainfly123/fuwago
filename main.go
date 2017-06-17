package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func queryVideo(w http.ResponseWriter, req *http.Request) {
	geohash := req.FormValue("geohash")
	class := req.FormValue("class")
	if len(geohash) < 5 || len(class) < 1 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	type JsonResponseData struct {
		JsonResponse
		Data []VideoResp `json:data`
	}
	temp := strings.Split(geohash, "-")
	longitude, _ := strconv.ParseFloat(temp[0], 32)
	latitude, _ := strconv.ParseFloat(temp[1], 32)
	jsonres := JsonResponseData{JsonResponse{0, "OK"}, QueryVideo(longitude, latitude, class)}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}

func main() {

	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile)
	InitRedis()
	fmt.Println("ok")
	http.HandleFunc("/queryvideo", queryVideo)
	log.Fatal(http.ListenAndServe(":1688", nil))

}
