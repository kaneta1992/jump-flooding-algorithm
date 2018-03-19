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

	swap := util.JumpFlooding(img)
	img2 := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			point := swap.Get(x, y).Nearest
			col := swap.Get(point.X, point.Y).Color
			img2.Set(x, y, col)
		}
	}
	out, _ := os.Create("out.png")
	defer out.Close()
	png.Encode(out, img2)
}

// TODO:
// ピクセルまでの距離を計算する
