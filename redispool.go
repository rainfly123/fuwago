package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
)

var Clients *pool.Pool

func InitRedis() {
	var err error
	Clients, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		// handle error
	}

}

func QueryVideo(longitude, latitude float32, classid string) {
	conn, err := Clients.Get()
	if err != nil {
		// handle error
	}
	r := conn.Cmd("AUTH", "aaa11bbb22")
	r = conn.Cmd("GEORADIUS", "fuwa_c", longitude, latitude, 10, "km")
	r.List()
	l, _ := r.List()
	for _, elemStr := range l {
		fmt.Println(elemStr)
	}
	Clients.Put(conn)
}

func main() {
	InitRedis()
	QueryVideo(113.29, 23.08, "1")
}
