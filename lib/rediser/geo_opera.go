package rediser

import (
	"github.com/go-redis/redis"
)

var (
	disM  = "m"  // 表示单位为米。
	disKM = "km" //表示单位为千米。
	disMI = "mi" //表示单位为英里。
	disFT = "ft" //表示单位为英尺。
)

// GeoMember geo成员
// Name成员名称，longitude, latitude经度，维度，dist距离，geohash
type GeoMember struct {
	Name      string
	Longitude float64
	Latitude  float64
	// Dist:      dist,
	// GeoHash:   geohash,
}

//GeoRadiusQuerys geo参数
//query中参数介绍
//key关键码，member 成员名称，unit单位，sort（asc/desc），count（相当于sql的limit）反馈数量,radius与指定位置的距离，longitude, latitude经度纬度
//withdist withcoord 返回位置名称的中心距离 和经纬度，WithGeoHash以52位有符号整数的形式,返回位置元素经过原始geohash编码的有序集合分值
//store将返回结果的地理位置信息保存到指定键
//StoreDist将返回结果距离中心节点的距离保存到指定键
type GeoRadiusQuerys struct {
	Key       string
	Member    string
	Longitude float64
	Latitude  float64
	Radius    float64
	// Can be m, km, ft, or mi. Default is km.
	Unit        string
	WithCoord   bool
	WithDist    bool
	WithGeoHash bool
	Count       int
	// Can be ASC or DESC. Default is no sort order.
	Sort      string
	Store     string
	StoreDist string
}

//AddGeoMesber 增加要给地理位置对象,可输入多个自定义的GeoMember对象
//key关键，memberName成员名称，longitude, latitude经度，维度，dist距离，geohash
func AddGeoMesber(rd *redis.Client, key string, gm ...*GeoMember) error {
	geoLoc := []*redis.GeoLocation{}
	for _, v := range gm {
		geoLocation := redis.GeoLocation{
			Name:      v.Name,
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
			// Dist:      dist,
			// GeoHash:   geohash,
		}
		geoLoc = append(geoLoc, &geoLocation)
	}
	return rd.GeoAdd(key, geoLoc...).Err()
}

//GetGeoMember 获取一个地理位置成员，返回一个经纬度对象数组
func GetGeoMember(rd *redis.Client, key, memberName string) ([]*redis.GeoPos, error) {
	return rd.GeoPos(key, memberName).Result()
}

//GetMemberDistance 获取两个成员之间的距离，输入关键key，对应下的两个成员名称，以及距离的单位，返回距离和error
func GetMemberDistance(rd *redis.Client, key, member1, member2, unit string) (float64, error) {
	return rd.GeoDist(key, member1, member2, unit).Result()
}

//GetRadius 使用经纬度获取成员附近成员,
//返回包含内容的对象数组和error
func GetRadius(rd *redis.Client, geoQuery GeoRadiusQuerys) ([]redis.GeoLocation, error) {
	rg := &redis.GeoRadiusQuery{
		Radius: geoQuery.Radius,
		// Can be m, km, ft, or mi. Default is km.
		Unit:        geoQuery.Unit,
		WithCoord:   geoQuery.WithCoord,
		WithDist:    geoQuery.WithDist,
		WithGeoHash: geoQuery.WithGeoHash,
		Count:       geoQuery.Count,
		// Can be ASC or DESC. Default is no sort order.
		Sort:      geoQuery.Sort,
		Store:     geoQuery.Store,
		StoreDist: geoQuery.StoreDist,
	}
	return rd.GeoRadius(geoQuery.Key, geoQuery.Longitude, geoQuery.Latitude, rg).Result()
}

//GetRadiusMember 使用名称获取成员附近成员，参数内容与GetRadius相同，使用成员名称获取附近的其他成员
func GetRadiusMember(rd *redis.Client, geoQuery GeoRadiusQuerys) ([]redis.GeoLocation, error) {
	gq := redis.GeoRadiusQuery{
		Radius: geoQuery.Radius,
		// Can be m, km, ft, or mi. Default is km.
		Unit:        geoQuery.Unit,
		WithCoord:   geoQuery.WithCoord,
		WithDist:    geoQuery.WithDist,
		WithGeoHash: geoQuery.WithGeoHash,
		Count:       geoQuery.Count,
		// Can be ASC or DESC. Default is no sort order.
		Sort:      geoQuery.Sort,
		Store:     geoQuery.Store,
		StoreDist: geoQuery.StoreDist,
	}
	return rd.GeoRadiusByMember(geoQuery.Key, geoQuery.Member, &gq).Result()
}

//GetGeoHash 将二维经纬度转为一维字符串，字符串越长表示位置更精确,两个字符串越相似表示距离越近
//输入多个成员，反馈对应顺序的结果集和error
func GetGeoHash(rd *redis.Client, key string, member ...string) ([]string, error) {
	return rd.GeoHash(key, member...).Result()
}

//DeleteGeoMember 删除一个geo成员
//这里注意：GEO没有提供删除成员的命令，但是因为GEO的底层实现是zset，所以可以借用zrem命令实现对地理位置信息的删除
func DeleteGeoMember(rd *redis.Client, key string, member ...string) (int64, error) {
	zm := []interface{}{}
	for _, v := range member {
		zm = append(zm, v)
	}
	return rd.ZRem(key, zm...).Result()
}
