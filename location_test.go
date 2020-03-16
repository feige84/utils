package utils

import (
	"fmt"
	"testing"
)

func TestLoc(t *testing.T) {
	currentLat := 24.335931
	currentLng := 109.448426
	//两点距离差不多在600米左右
	destLat := 24.330819
	destLng := 109.451471
	fmt.Println(EarthDistance(currentLat, currentLng, destLat, destLng))

	leftTop, rightTop, leftBottom, rightBottom := SquarePoint(currentLat, currentLng, 5000)
	Pr(leftTop)
	Pr(rightTop)
	Pr(leftBottom)
	Pr(rightBottom)
}
