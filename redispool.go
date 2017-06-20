package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	//"github.com/mediocregopher/radix.v2/redis"
	"encoding/json"
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

type VideoResp struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Userid    string `json:"userid"`
	Video     string `json:"video"`
	Width     string `json:"width"`
	Height    string `json:"height"`
	Distances string `json:"distance"`
	Filemd5   string `json:"filemd5"`
}
type Fuwa struct {
	Distance  float32 `json:"distance"`
	Pic       string  `json:"pic"`
	Gid       string  `json:"gid"`
	Geo       string  `json:"geo"`
	Pos       string  `json:"pos"`
	Detail    string  `json:"detail"`
	Avatar    string  `json:"avatar"`
	Name      string  `json:"name"`
	Gender    string  `json:"gender"`
	Signature string  `json:"signature"`
	Location  string  `json:"location"`
	Video     string  `json:"video"`
	Hider     string  `json:"hider"`
}

type nearFuwa struct {
	Fuwa
	Id string `json:"id"`
}

type farFuwa struct {
	Fuwa
	Number uint32 `json:"number"`
}

func QueryVideo(longitude, latitude float64, classid string) []VideoResp {
	var results []VideoResp
	conn, err := Clients.Get()
	if err != nil {
		// handle error
	}
	r := conn.Cmd("AUTH", "aaa11bbb22")
	r = conn.Cmd("ZREVRANGE", "video_"+classid, 0, 4)
	filemd5s, _ := r.List()
	total := len(filemd5s)
	distances := make(map[int]string, 5)
	r = conn.Cmd("GEOPOS", "video_g_"+classid, filemd5s)
	posa, _ := r.Array()
	for i, elem := range posa {
		pos, _ := elem.List()
		lonti, _ := strconv.ParseFloat(pos[0], 32)
		lati, _ := strconv.ParseFloat(pos[1], 32)
		dis := EarthDistance(lati, lonti, latitude, longitude)
		distances[i] = strconv.Itoa(dis)
	}
	for i, filemd5 := range filemd5s {
		r = conn.Cmd("HMGET", filemd5, "name", "gender", "avatar", "userid", "video", "width", "height")
		resp, _ := r.List()
		temp := VideoResp{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], distances[i], filemd5}
		results = append(results, temp)
	}

	r = conn.Cmd("GEORADIUS", "video_g_"+classid, longitude, latitude, 10000, "m", "withdist", "count", "100", "ASC")
	posa, _ = r.Array()
	for _, elem := range posa {
		var had bool
		pos, _ := elem.List()
		filemd5 := pos[0]
		dis := pos[1]
		r = conn.Cmd("HMGET", filemd5, "name", "gender", "avatar", "userid", "video", "width", "height")
		resp, _ := r.List()
		temp := VideoResp{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], dis, filemd5}
		had = false
		for i := 0; i < total; i++ {
			if results[i].Filemd5 == filemd5 {
				had = true
			}
		}
		if had != true {
			results = append(results, temp)
		}
	}
	defer Clients.Put(conn)
	return results
}

func QueryStrVideo(longitude, latitude float64) []VideoResp {
	var results []VideoResp
	conn, err := Clients.Get()
	if err != nil {
		// handle error
	}
	r := conn.Cmd("AUTH", "aaa11bbb22")
	r = conn.Cmd("ZREVRANGE", "video_i", 0, 4)
	filemd5s, _ := r.List()
	total := len(filemd5s)
	distances := make(map[int]string, 5)
	r = conn.Cmd("GEOPOS", "video_g_i", filemd5s)
	posa, _ := r.Array()
	for i, elem := range posa {
		pos, _ := elem.List()
		lonti, _ := strconv.ParseFloat(pos[0], 32)
		lati, _ := strconv.ParseFloat(pos[1], 32)
		dis := EarthDistance(lati, lonti, latitude, longitude)
		distances[i] = strconv.Itoa(dis)
	}
	for i, filemd5 := range filemd5s {
		r = conn.Cmd("HMGET", filemd5, "name", "gender", "avatar", "userid", "video", "width", "height")
		resp, _ := r.List()
		temp := VideoResp{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], distances[i], filemd5}
		results = append(results, temp)
	}

	r = conn.Cmd("GEORADIUS", "video_g_i", longitude, latitude, 10000, "m", "withdist", "count", "100", "ASC")
	posa, _ = r.Array()
	for _, elem := range posa {
		var had bool
		pos, _ := elem.List()
		filemd5 := pos[0]
		dis := pos[1]
		r = conn.Cmd("HMGET", filemd5, "name", "gender", "avatar", "userid", "video", "width", "height")
		resp, _ := r.List()
		temp := VideoResp{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], dis, filemd5}
		had = false
		for i := 0; i < total; i++ {
			if results[i].Filemd5 == filemd5 {
				had = true
			}
		}
		if had != true {
			results = append(results, temp)
		}
	}
	defer Clients.Put(conn)
	return results
}

const HOWFAR = 300

type GEORADIUSRESP struct {
	Fuwagid  string
	Distance string
}
type ByFuwagid []GEORADIUSRESP

func (a ByFuwagid) Len() int           { return len(a) }
func (a ByFuwagid) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFuwagid) Less(i, j int) bool { return a[i].Fuwagid[7:] < a[j].Fuwagid[7:] }

func QueryV2(longtitude, latitude float64, radius uint32, biggest string) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	var response GEORADIUSRESP

	result := make(map[string]interface{}, 2)
	conn, err := Clients.Get()
	if err != nil {
		// handle error
	}

	r := conn.Cmd("AUTH", "aaa11bbb22")
	r = conn.Cmd("GEORADIUS", "fuwa_c", longitude, latitude, radius, "m", "withdist")
	nelem, _ = r.Array()
	for _, elem := range nelem {
		temp, _ := elem.List()
		if temp[0][7:] < biggest && temp[1] < HOWFAR {
			fuwa := GEORADIUSRESP{temp[0], temp[1]}
			response = append(response, fuwa)
		}

	}
	return result
}

func QueryStrV2(longtitude, latitude float64, radius uint32, biggest string) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	result := make(map[string]interface{}, 2)

	return result
}
func QueryV3(longtitude, latitude float64, radius uint32, biggest string, creator uint32) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	result := make(map[string]interface{}, 2)

	return result
}

func QueryStrV3(longtitude, latitude float64, radius uint32, biggest string, creator uint32) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	result := make(map[string]interface{}, 2)

	return result
}

func main() {
	//InitRedis()
	//fmt.Println(QueryVideo(113.301, 23.0827, "1"))
	b, _ := json.Marshal(QueryV2(0, 0, 0, "0"))
	fmt.Println(string(b))
}
