package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/kaneta1992/jump-flooding-algorithm/src/util"
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

	swap := util.NewSwapBuffer(bounds.Max.X, bounds.Max.Y)
	swap_img := swap.GetActiveBuffer()
	draw.Draw(swap_img, swap_img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, a := rgba(col.RGBA())
			swap_img.Set(x, y, color.RGBA{uint8(float64(r) * 0.5), uint8(float64(g) * 0.5), uint8(float64(b) * 0.5), uint8(a)})
		}
	}
	swap_img.Set(64, 64, color.RGBA{255, 0, 255, 255})

	out, _ := os.Create("out.png")
	defer out.Close()

	png.Encode(out, swap_img)
}

// TODO:
// 配列に画像を読み込む
// 処理用のプライマリ、セカンダリバッファを用意する
// 処理ごとにバッファを切り替える仕組みを作る
// JFAで最近のピクセルを計算する
// ピクセルまでの距離を計算する
// 保存
