package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func rgba(r uint32, g uint32, b uint32, a uint32) (uint8, uint8, uint8, uint8) {
	return uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)
}

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

	img2 := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, a := rgba(col.RGBA())
			if r == 255 && g == 255 && b == 255 {
				a = 0
			}
			img2.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	out, _ := os.Create("out.png")
	defer out.Close()

	png.Encode(out, img2)
	//util.JumpFlooding(img)
}

// TODO:
// JFAで最近のピクセルを計算する
// ピクセルまでの距離を計算する
