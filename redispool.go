package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	//"github.com/mediocregopher/radix.v2/redis"
	"math"
	"strconv"
)

var Clients *pool.Pool

func InitRedis() {
	var err error
	Clients, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		// handle error
	}

}

func EarthDistance(lat1, lng1, lat2, lng2 float64) int {
	radius := 6371000.0 // 6378137
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return int(dist * radius)
}

func QueryVideo(longitude, latitude float64, classid string) {
	conn, err := Clients.Get()
	if err != nil {
		// handle error
	}
	r := conn.Cmd("AUTH", "aaa11bbb22")
	//r = conn.Cmd("GEORADIUS", "fuwa_c", longitude, latitude, 10, "km")
	//r = conn.Cmd("HMGET", "fuwa_c_2294", "name", "pos")
	r = conn.Cmd("ZREVRANGE", "video_"+classid, 0, 4)
	filemd5s, _ := r.List()
	//for _, elemStr := range filemd5s {
	//	fmt.Println(elemStr)
	//}
	var distances = map[int]int{0: 0, 1: 0, 2: 0, 3: 0, 4: 0}
	r = conn.Cmd("GEOPOS", "video_g_"+classid, filemd5s)
	posa, _ := r.Array()
	for i, elem := range posa {
		pos, _ := elem.List()
		lonti, _ := strconv.ParseFloat(pos[0], 32)
		lati, _ := strconv.ParseFloat(pos[1], 32)
		//fmt.Println(lonti, lati)
		dis := EarthDistance(lati, lonti, latitude, longitude)
		distances[i] = dis
	}
	fmt.Println(distances)

	defer Clients.Put(conn)
}

func main() {
	InitRedis()
	QueryVideo(113.29, 23.08, "1")
}
