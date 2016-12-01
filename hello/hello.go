//Output command line args
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	//helloWorld()
	//echo1()
	//echo2()
	//echo3()
	//echo4()
	echoExercise1()
}

func helloWorld() {
	fmt.Printf("Hello, World\n")
}

func echo1() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

var s string

func echo2() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " 🐱 "))
}

func echo4() {
	fmt.Println(os.Args[1:])
}

func echoExercise1() {
	for index, arg := range os.Args[1:] {
		fmt.Println(index, " "+arg)
	}
}
