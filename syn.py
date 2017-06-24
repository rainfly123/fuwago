#!/usr/bin/env python
#-*- coding: utf-8 -*- 
import redis  

pool = redis.ConnectionPool(host='127.0.0.1', port=6379, db=0, password="aaa11bbb22")  
r = redis.Redis(connection_pool=pool)  

def QueryFuwaNew():
    for fuwa in near:
        fuwaid, distance = fuwa[0], fuwa[1]
        detail, pos, pic, idd = r.hmget(fuwaid, "detail", "pos", "pic", "id")
        geohash = r.geopos("fuwa_c", fuwaid)
        geo = "%f-%f"%(geohash[0][0], geohash[0][1])

fuwas = r.zrange("fuwa_c", 0, -1)
for fuwa in fuwas:
    creator = r.hget(fuwa, "creator")
    pos = r.geopos("fuwa_c", fuwa)[0]
    r.geoadd("fuwa_c_" + creator, pos[0], pos[1], fuwa)
    
fuwas = r.zrange("fuwa_i", 0, -1)
for fuwa in fuwas:
    creator = r.hget(fuwa, "creator")
    pos = r.geopos("fuwa_i", fuwa)[0]
    r.geoadd("fuwa_i_" + creator, pos[0], pos[1], fuwa)
 
