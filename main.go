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
	sdfImage := jfa.CalcSDF(8.0)

	voronoiFile, _ := os.Create("voronoi.png")
	defer voronoiFile.Close()
	png.Encode(voronoiFile, voronoiImage)

	sdfFile, _ := os.Create("sdf.png")
	defer sdfFile.Close()
	png.Encode(sdfFile, sdfImage)
}
