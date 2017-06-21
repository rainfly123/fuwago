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
		Data []VideoResp `json:"data"`
	}
	temp := strings.Split(geohash, "-")
	longitude, _ := strconv.ParseFloat(temp[0], 32)
	latitude, _ := strconv.ParseFloat(temp[1], 32)
	jsonres := JsonResponseData{JsonResponse{0, "OK"}, QueryVideo(longitude, latitude, class)}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}

func queryStrVideo(w http.ResponseWriter, req *http.Request) {
	geohash := req.FormValue("geohash")
	if len(geohash) < 5 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	type JsonResponseData struct {
		JsonResponse
		Data []VideoResp `json:"data"`
	}
	temp := strings.Split(geohash, "-")
	longitude, _ := strconv.ParseFloat(temp[0], 32)
	latitude, _ := strconv.ParseFloat(temp[1], 32)
	jsonres := JsonResponseData{JsonResponse{0, "OK"}, QueryStrVideo(longitude, latitude)}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}

func queryV2Handler(w http.ResponseWriter, req *http.Request) {
	geohash := req.FormValue("geohash")
	radius := req.FormValue("radius")
	biggest := req.FormValue("biggest")
	if len(geohash) < 5 || len(biggest) < 1 || len(radius) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	type JsonResponseData struct {
		JsonResponse
		Data map[string]interface{} `json:"data"`
	}
	temp := strings.Split(geohash, "-")
	longitude, _ := strconv.ParseFloat(temp[0], 32)
	latitude, _ := strconv.ParseFloat(temp[1], 32)
	radiusint, _ := strconv.Atoi(radius)
	big, _ := strconv.Atoi(biggest)
	if big == 0 {
		big = 999999999
	}
	jsonres := JsonResponseData{JsonResponse{0, "OK"}, QueryV2(longitude, latitude, uint32(radiusint), big)}
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
	http.HandleFunc("/querystrvideo", queryStrVideo)
	http.HandleFunc("/queryv2", queryV2Handler)
	log.Fatal(http.ListenAndServe(":9999", nil))

}
