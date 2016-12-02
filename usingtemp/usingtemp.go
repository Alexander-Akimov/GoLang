package main

import (
	"fmt"

	"github.com/akimov/tempconv"
)

func main() {
	// var c tempconv.Celsius
	// celGrad := tempconv.Celsius(c)

	// fmt.Print(celGrad)

	fmt.Printf("%s\n", tempconv.BoilingC)
	boilinfF := tempconv.CToF(tempconv.BoilingC)
	fmt.Println(boilinfF)
	fmt.Printf("%g\n", boilinfF)
}
