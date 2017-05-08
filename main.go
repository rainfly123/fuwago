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

func grpinfoHandle(w http.ResponseWriter, req *http.Request) {
	groupid := req.FormValue("groupid")
	if len(groupid) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool

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
	jsonres := JsonResponse{0, "OK"}
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
