package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/kaneta1992/jump-flooding-algorithm/src/util"
)

func main() {
	file, _ := os.Open("./test.png")
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	fmt.Println(bounds.Max.X)
	fmt.Println(bounds.Max.Y)

	jfa := util.NewJFA(img)
	voronoiImage := jfa.CalcVoronol()

	out, _ := os.Create("out.png")
	defer out.Close()
	png.Encode(out, voronoiImage)
}
