package main

import (
	"fmt"
	"os"

	minCDN "github.com/ZhiYuanHuang/minCDN/cmd"
)

func main() {
	fmt.Println("hello minCDN")

	minCDN.Main(os.Args)

	fmt.Println("bye")
}
