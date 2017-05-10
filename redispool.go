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
func main() {
	InitRedis()
	conn, err := Clients.Get()
	if err != nil {
		// handle error
	}

	r := conn.Cmd("GET", "globalid")
	c, _ := r.Str()
	fmt.Println(c)

	r = conn.Cmd("GEORADIUS", "fuwa_c", 30000)
	fmt.Println(r)

	Clients.Put(conn)
}
