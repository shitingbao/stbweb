package rediser

import "github.com/go-redis/redis"

//AddGeoMesber 增加要给地理位置对象
func AddGeoMesber(rd *redis.Client, key, memberName string, longitude, latitude, dist float64, geohash int64) {
	gm := redis.GeoLocation{
		Name:      memberName,
		Longitude: longitude,
		Latitude:  latitude,
		Dist:      dist,
		GeoHash:   geohash,
	}
	rd.GeoAdd(key, &gm)
}

//GetGeoMember 获取一个地理位置成员
func GetGeoMember(rd *redis.Client, key, memberName string) ([]*redis.GeoPos, error) {
	return rd.GeoPos(key, memberName).Result()
}
