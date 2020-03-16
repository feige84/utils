package utils

import "math"

const (
	EarthRadius = float64(6371000) // 6378137				//地球半径，平均半径为6371km
)

//$ad = Db::connect('databasetwo')->query("SELECT *,(2 * 6378.137 * ASIN(	SQRT(POW( SIN( PI( ) * ( " . 用户$longitude . "- 查询表.longitude ) / 360 ), 2 ) + COS( PI( ) * " . 用户$latitude . " / 180 ) * COS(  查询表.latitude * PI( ) / 180 ) * POW( SIN( PI( ) * ( " . 用户$latitude . "- 查询表.latitude ) / 360 ), 2 )))) AS distance FROM `查询表` ORDER BY	distance ASC LIMIT 1");
//使用此函数计算得到结果后，带入sql查询。

//currentLat 当前纬度
//currentLng 当前经度
//destLat 目的地纬度
//destLng 目的地经度

type Coordinate struct {
	Longitude float64 //经度
	Latitude  float64 //纬度
}

/**
*计算某个经纬度的周围某段距离的正方形的四个点
*
*@param lng float 经度
*@param lat float 纬度
*@param distance float 该点所在圆的半径，该圆与此正方形内切，单位：米
*@return array 正方形的四个点的经纬度坐标
 */
func SquarePoint(lat, lng, distance float64) (leftTop, rightTop, leftBottom, rightBottom Coordinate) {
	//rad := deg2rad(lat)
	dLng := 2 * math.Asin(math.Sin(distance/(2*EarthRadius))/math.Cos(lat*math.Pi/180))
	dLng = dLng * 180 / math.Pi //rad2deg(dLng)

	dLat := distance / EarthRadius
	dLat = dLat * 180 / math.Pi //rad2deg(dLat)

	leftTop.Latitude = lat + dLat
	leftTop.Longitude = lng - dLng

	rightTop.Latitude = lat + dLat
	rightTop.Longitude = lng + dLng

	leftBottom.Latitude = lat - dLat
	leftBottom.Longitude = lng - dLng

	rightBottom.Latitude = lat - dLat
	rightBottom.Longitude = lng + dLng
	return
}

// 返回值的单位为米
func EarthDistance(currentLat, currentLng, destLat, destLng float64) float64 {
	rad := math.Pi / 180.0
	currentLat = currentLat * rad
	currentLng = currentLng * rad
	destLat = destLat * rad
	destLng = destLng * rad
	theta := destLng - currentLng
	dist := math.Acos(math.Sin(currentLat)*math.Sin(destLat) + math.Cos(currentLat)*math.Cos(destLat)*math.Cos(theta))
	return dist * EarthRadius
}
