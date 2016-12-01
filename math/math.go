package main

import (
	"fmt"
)

func main() {
	fmt.Println(gcd(28, 64))
	fmt.Println(28 % 29)
	// var r = nil
	// fmt.Println(r)

}
func gcd(x, y int) int { //Алгоритм Евклида http://younglinux.info/algorithm/euclidean
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
