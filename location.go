package utils

import "math"

//currentLat 当前纬度
//currentLng 当前经度
//destLat 目的地纬度
//destLng 目的地经度

// 返回值的单位为米
func EarthDistance(currentLat, currentLng, destLat, destLng float64) float64 {
	radius := float64(6371000) // 6378137
	rad := math.Pi / 180.0
	currentLat = currentLat * rad
	currentLng = currentLng * rad
	destLat = destLat * rad
	destLng = destLng * rad
	theta := destLng - currentLng
	dist := math.Acos(math.Sin(currentLat)*math.Sin(destLat) + math.Cos(currentLat)*math.Cos(destLat)*math.Cos(theta))
	return dist * radius
}
