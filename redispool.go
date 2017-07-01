package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	//"encoding/json"
	"math"
	"sort"
	"strconv"
)

var Clients *pool.Pool

func InitRedis() {
	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		if err = client.Cmd("AUTH", "aaa11bbb22").Err; err != nil {
		}
		return client, nil
	}
	Clients, _ = pool.NewCustom("tcp", "127.0.0.1:6379", 10, df)
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
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Avatar   string `json:"avatar"`
	Userid   string `json:"userid"`
	Video    string `json:"video"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Distance string `json:"distance"`
	Filemd5  string `json:"filemd5"`
}
type Fuwa struct {
	Detail    string  `json:"detail"`
	Pos       string  `json:"pos"`
	Pic       string  `json:"pic"`
	Name      string  `json:"name"`
	Avatar    string  `json:"avatar"`
	Gender    string  `json:"gender"`
	Signature string  `json:"signature"`
	Location  string  `json:"location"`
	Video     string  `json:"video"`
	Hider     string  `json:"hider"`
	Geo       string  `json:"geo"`
	Distance  float32 `json:"distance"`
}

type nearFuwa struct {
	Fuwa
	Id  string `json:"id"`
	Gid string `json:"gid"`
}

type farFuwa struct {
	Fuwa
	Number uint32 `json:"number"`
}

func QueryVideo(longitude, latitude float64, classid string) []VideoResp {
	var results []VideoResp
	var nearvideo []string
	var distances []string

	conn, err := Clients.Get()
	if err != nil {
		return results
	}
	r := conn.Cmd("ZREVRANGE", "video_"+classid, 0, 9)
	filemd5s, _ := r.List()
	r = conn.Cmd("GEOPOS", "video_g_"+classid, filemd5s)
	posa, _ := r.Array()
	for i, elem := range posa {
		pos, _ := elem.List()
		lonti, _ := strconv.ParseFloat(pos[0], 32)
		lati, _ := strconv.ParseFloat(pos[1], 32)
		dis := EarthDistance(lati, lonti, latitude, longitude)
		if dis <= 20000 {
			distances = append(distances, strconv.Itoa(dis))
			nearvideo = append(nearvideo, filemd5s[i])
		}
	}
	total := len(nearvideo)
	if total == 0 {
		for i, elem := range posa {
			pos, _ := elem.List()
			lonti, _ := strconv.ParseFloat(pos[0], 32)
			lati, _ := strconv.ParseFloat(pos[1], 32)
			dis := EarthDistance(lati, lonti, latitude, longitude)
			distances = append(distances, strconv.Itoa(dis))
			nearvideo = append(nearvideo, filemd5s[i])
		}
	}
	total = len(nearvideo)
	for i, filemd5 := range nearvideo {
		r = conn.Cmd("HMGET", filemd5, "name", "gender", "avatar", "userid", "video", "width", "height")
		resp, _ := r.List()
		temp := VideoResp{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], distances[i], filemd5}
		results = append(results, temp)
	}

	r = conn.Cmd("GEORADIUS", "video_g_"+classid, longitude, latitude, 20000, "m", "withdist", "count", "100", "ASC")
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
	var nearvideo []string
	var distances []string

	conn, err := Clients.Get()
	if err != nil {
		return results
	}
	r := conn.Cmd("ZREVRANGE", "video_i", 0, 9)
	filemd5s, _ := r.List()
	r = conn.Cmd("GEOPOS", "video_g_i", filemd5s)
	posa, _ := r.Array()
	for i, elem := range posa {
		pos, _ := elem.List()
		lonti, _ := strconv.ParseFloat(pos[0], 32)
		lati, _ := strconv.ParseFloat(pos[1], 32)
		dis := EarthDistance(lati, lonti, latitude, longitude)
		if dis <= 20000 {
			distances = append(distances, strconv.Itoa(dis))
			nearvideo = append(nearvideo, filemd5s[i])
		}
	}
	total := len(nearvideo)
	if total == 0 {
		for i, elem := range posa {
			pos, _ := elem.List()
			lonti, _ := strconv.ParseFloat(pos[0], 32)
			lati, _ := strconv.ParseFloat(pos[1], 32)
			dis := EarthDistance(lati, lonti, latitude, longitude)
			distances = append(distances, strconv.Itoa(dis))
			nearvideo = append(nearvideo, filemd5s[i])
		}
	}
	total = len(nearvideo)
	for i, filemd5 := range nearvideo {
		r = conn.Cmd("HMGET", filemd5, "name", "gender", "avatar", "userid", "video", "width", "height")
		resp, _ := r.List()
		temp := VideoResp{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], distances[i], filemd5}
		results = append(results, temp)
	}

	r = conn.Cmd("GEORADIUS", "video_g_i", longitude, latitude, 20000, "m", "withdist", "count", "100", "ASC")
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

func (a ByFuwagid) Len() int      { return len(a) }
func (a ByFuwagid) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByFuwagid) Less(i, j int) bool {
	m, _ := strconv.Atoi(a[i].Fuwagid[7:])
	n, _ := strconv.Atoi(a[j].Fuwagid[7:])
	return m < n
}

func QueryV2(longitude, latitude float64, radius uint32, biggest int) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	var nresponse ByFuwagid
	var fresponse ByFuwagid

	result := make(map[string]interface{}, 2)
	conn, err := Clients.Get()
	if err != nil {
		return result
	}

	r := conn.Cmd("GEORADIUS", "fuwa_c", longitude, latitude, radius, "m", "withdist", "count", "500")
	nelem, _ := r.Array()
	for _, elem := range nelem {
		temp, _ := elem.List()
		howfar, _ := strconv.ParseFloat(temp[1], 32)
		if howfar < HOWFAR {
			fuwagidn, _ := strconv.Atoi(temp[0][7:])
			if fuwagidn < biggest && len(nresponse) <= 100 {
				fuwa := GEORADIUSRESP{temp[0], temp[1]}
				nresponse = append(nresponse, fuwa)
			}
		} else if len(fresponse) <= 300 {
			fuwa := GEORADIUSRESP{temp[0], temp[1]}
			fresponse = append(fresponse, fuwa)
		}

	}
	sort.Sort(sort.Reverse(nresponse))
	for _, v := range nresponse {
		var geo string
		r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
			"gender", "signature", "location", "video", "owner", "id")
		resp, _ := r.List()

		r = conn.Cmd("GEOPOS", "fuwa_c", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		dis, _ := strconv.ParseFloat(v.Distance, 32)

		temp := nearFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
			resp[8], resp[9], geo, float32(dis)}, resp[10], v.Fuwagid}
		nearfuwa = append(nearfuwa, temp)
	}
	for _, v := range fresponse {
		var geo string
		var has bool
		has = false
		r = conn.Cmd("GEOPOS", "fuwa_c", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		for i, shit := range farfuwa {
			if shit.Geo == geo {
				farfuwa[i].Number += 1
				has = true
			}
		}
		if has == false {

			r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
				"gender", "signature", "location", "video", "owner")
			resp, _ := r.List()

			dis, _ := strconv.ParseFloat(v.Distance, 32)
			temp := farFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
				resp[8], resp[9], geo, float32(dis)}, 1}
			farfuwa = append(farfuwa, temp)
		}
	}
	result["near"] = nearfuwa
	result["far"] = farfuwa
	defer Clients.Put(conn)
	return result
}

func QueryStrV2(longitude, latitude float64, radius uint32, biggest int) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	var nresponse ByFuwagid
	var fresponse ByFuwagid

	result := make(map[string]interface{}, 2)
	conn, err := Clients.Get()
	if err != nil {
		return result
	}

	r := conn.Cmd("GEORADIUS", "fuwa_i", longitude, latitude, radius, "m", "withdist", "count", "500")
	nelem, _ := r.Array()
	for _, elem := range nelem {
		temp, _ := elem.List()
		howfar, _ := strconv.ParseFloat(temp[1], 32)
		if howfar < HOWFAR {
			fuwagidn, _ := strconv.Atoi(temp[0][7:])
			if fuwagidn < biggest && len(nresponse) <= 100 {
				fuwa := GEORADIUSRESP{temp[0], temp[1]}
				nresponse = append(nresponse, fuwa)
			}
		} else if len(fresponse) <= 300 {
			fuwa := GEORADIUSRESP{temp[0], temp[1]}
			fresponse = append(fresponse, fuwa)
		}

	}
	sort.Sort(sort.Reverse(nresponse))
	for _, v := range nresponse {
		var geo string
		r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
			"gender", "signature", "location", "video", "owner", "id")
		resp, _ := r.List()

		r = conn.Cmd("GEOPOS", "fuwa_i", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		dis, _ := strconv.ParseFloat(v.Distance, 32)

		temp := nearFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
			resp[8], resp[9], geo, float32(dis)}, resp[10], v.Fuwagid}
		nearfuwa = append(nearfuwa, temp)
	}
	for _, v := range fresponse {
		var geo string
		var has bool
		has = false
		r = conn.Cmd("GEOPOS", "fuwa_i", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		for i, shit := range farfuwa {
			if shit.Geo == geo {
				farfuwa[i].Number += 1
				has = true
			}
		}
		if has == false {

			r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
				"gender", "signature", "location", "video", "owner")
			resp, _ := r.List()

			dis, _ := strconv.ParseFloat(v.Distance, 32)
			temp := farFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
				resp[8], resp[9], geo, float32(dis)}, 1}
			farfuwa = append(farfuwa, temp)
		}
	}
	result["near"] = nearfuwa
	result["far"] = farfuwa
	defer Clients.Put(conn)
	return result
}

func QueryV3(longitude, latitude float64, radius uint32, biggest int, creator string) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	var nresponse ByFuwagid
	var fresponse ByFuwagid

	if radius > 20000 {
		radius = 20000
	}
	result := make(map[string]interface{}, 2)
	conn, err := Clients.Get()
	if err != nil {
		return result
	}

	r := conn.Cmd("GEORADIUS", "fuwa_c_"+creator, longitude, latitude, radius, "m", "withdist", "count", "500")
	nelem, _ := r.Array()
	for _, elem := range nelem {
		temp, _ := elem.List()
		howfar, _ := strconv.ParseFloat(temp[1], 32)
		if howfar < HOWFAR {
			fuwagidn, _ := strconv.Atoi(temp[0][7:])
			if fuwagidn < biggest && len(nresponse) <= 100 {
				fuwa := GEORADIUSRESP{temp[0], temp[1]}
				nresponse = append(nresponse, fuwa)
			}
		} else if len(fresponse) <= 300 {
			fuwa := GEORADIUSRESP{temp[0], temp[1]}
			fresponse = append(fresponse, fuwa)
		}

	}
	sort.Sort(sort.Reverse(nresponse))
	for _, v := range nresponse {
		var geo string
		r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
			"gender", "signature", "location", "video", "owner", "id")
		resp, _ := r.List()

		r = conn.Cmd("GEOPOS", "fuwa_c", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		dis, _ := strconv.ParseFloat(v.Distance, 32)

		temp := nearFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
			resp[8], resp[9], geo, float32(dis)}, resp[10], v.Fuwagid}
		nearfuwa = append(nearfuwa, temp)
	}
	for _, v := range fresponse {
		var geo string
		var has bool
		has = false
		r = conn.Cmd("GEOPOS", "fuwa_c", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		for i, shit := range farfuwa {
			if shit.Geo == geo {
				farfuwa[i].Number += 1
				has = true
			}
		}
		if has == false {

			r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
				"gender", "signature", "location", "video", "owner")
			resp, _ := r.List()

			dis, _ := strconv.ParseFloat(v.Distance, 32)
			temp := farFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
				resp[8], resp[9], geo, float32(dis)}, 1}
			farfuwa = append(farfuwa, temp)
		}
	}
	result["near"] = nearfuwa
	result["far"] = farfuwa
	defer Clients.Put(conn)
	return result
}

func QueryStrV3(longitude, latitude float64, radius uint32, biggest int, creator string) map[string]interface{} {
	var farfuwa []farFuwa
	var nearfuwa []nearFuwa
	var nresponse ByFuwagid
	var fresponse ByFuwagid

	if radius > 20000 {
		radius = 20000
	}
	result := make(map[string]interface{}, 2)
	conn, err := Clients.Get()
	if err != nil {
		return result
	}

	r := conn.Cmd("GEORADIUS", "fuwa_i_"+creator, longitude, latitude, radius, "m", "withdist", "count", "500")
	nelem, _ := r.Array()
	for _, elem := range nelem {
		temp, _ := elem.List()
		howfar, _ := strconv.ParseFloat(temp[1], 32)
		if howfar < HOWFAR {
			fuwagidn, _ := strconv.Atoi(temp[0][7:])
			if fuwagidn < biggest && len(nresponse) <= 100 {
				fuwa := GEORADIUSRESP{temp[0], temp[1]}
				nresponse = append(nresponse, fuwa)
			}
		} else if len(fresponse) <= 300 {
			fuwa := GEORADIUSRESP{temp[0], temp[1]}
			fresponse = append(fresponse, fuwa)
		}

	}
	sort.Sort(sort.Reverse(nresponse))
	for _, v := range nresponse {
		var geo string
		r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
			"gender", "signature", "location", "video", "owner", "id")
		resp, _ := r.List()

		r = conn.Cmd("GEOPOS", "fuwa_i", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		dis, _ := strconv.ParseFloat(v.Distance, 32)

		temp := nearFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
			resp[8], resp[9], geo, float32(dis)}, resp[10], v.Fuwagid}
		nearfuwa = append(nearfuwa, temp)
	}
	for _, v := range fresponse {
		var geo string
		var has bool
		has = false
		r = conn.Cmd("GEOPOS", "fuwa_i", v.Fuwagid)
		posa, _ := r.Array()
		for _, elem := range posa {
			pos, _ := elem.List()
			long, _ := strconv.ParseFloat(pos[0], 32)
			lat, _ := strconv.ParseFloat(pos[1], 32)
			geo = fmt.Sprintf("%f-%f", long, lat)
		}
		for i, shit := range farfuwa {
			if shit.Geo == geo {
				farfuwa[i].Number += 1
				has = true
			}
		}
		if has == false {

			r = conn.Cmd("HMGET", v.Fuwagid, "detail", "pos", "pic", "name", "avatar",
				"gender", "signature", "location", "video", "owner")
			resp, _ := r.List()

			dis, _ := strconv.ParseFloat(v.Distance, 32)
			temp := farFuwa{Fuwa{resp[0], resp[1], resp[2], resp[3], resp[4], resp[5], resp[6], resp[7],
				resp[8], resp[9], geo, float32(dis)}, 1}
			farfuwa = append(farfuwa, temp)
		}
	}
	result["near"] = nearfuwa
	result["far"] = farfuwa
	defer Clients.Put(conn)
	return result
}

/*
func main() {
	InitRedis()
	//fmt.Println(QueryVideo(113.301, 23.0827, "1"))
	//b, _ := json.Marshal(QueryV2(0, 0, 0, "0"))
	//fmt.Println(string(b))
	b, _ := json.Marshal(QueryV2(113.301, 23.0827, 10000, "300000"))
	fmt.Println(string(b))
}
*/
