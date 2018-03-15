package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

func main() {
	file, _ := os.Open("./test.png")
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Width)
	fmt.Println(config.Height)
}
