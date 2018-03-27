package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/kaneta1992/jump-flooding-algorithm/src/jfa"
)

func main() {
	images := 4
	for i := 0; i < images; i++ {
		file, _ := os.Open(fmt.Sprintf("./testImage/test%d.png", i))
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		bounds := img.Bounds()
		fmt.Println(bounds.Max.X)
		fmt.Println(bounds.Max.Y)

		jfa := jfa.NewJFA(img)
		voronoiImage := jfa.CalcVoronol()
		sdfImage := jfa.CalcSDF(8.0)

		voronoiFile, _ := os.Create(fmt.Sprintf("./testImage/voronoi%d.png", i))
		defer voronoiFile.Close()
		png.Encode(voronoiFile, voronoiImage)

		sdfFile, _ := os.Create(fmt.Sprintf("./testImage/sdf%d.png", i))
		defer sdfFile.Close()
		png.Encode(sdfFile, sdfImage)
	}
}
