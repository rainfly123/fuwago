package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"menteslibres.net/gosexy/redis"
	"net/http"
	"runtime"
	//	"strings"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Fuwa struct {
	Distance  float32 `json:"distance"`
	Pic       string  `json:"pic"`
	Gid       string  `json:"gid"`
	Geo       string  `json:"geo"`
	Pos       string  `json:"pos"`
	Id        string  `json:"id"`
	Detail    string  `json:"detail"`
	Avatar    string  `json:"avatar"`
	Name      string  `json:"name"`
	Gender    string  `json:"gender"`
	Signature string  `json:"signature"`
	Location  string  `json:"location"`
	Video     string  `json:"video"`
	Hider     string  `json:"hider"`
}

func QueryHandle(w http.ResponseWriter, req *http.Request) {
	geohash := req.FormValue("geohash")
	radius := req.FormValue("redius")
	if len(geohash) < 5 || len(radius) < 1 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool
	var fuwas []Fuwa

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	kkey := "group_" + groupid
	ls, _ := client.HMGet(kkey, "creator", "name", "notice", "snap")
	fmt.Println(ls)
	gkey := "groupmembers_" + groupid
	members, _ := client.SMembers(gkey)
	for _, member := range members {
		lss, _ := client.HMGet("user_"+member, "nick", "snap")
		fmt.Println(lss)
	}
	client.Close()
	type JsonResponseData struct {
		JsonResponse
		Data []Fuwa `json:data`
	}
	jsonres := JsonResponseData{0, "OK", fuwas}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}

func main() {

	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile)
	InitRedis()

	http.HandleFunc("/grpinfo", grpinfoHandle)
	log.Fatal(http.ListenAndServe(":1688", nil))

}
