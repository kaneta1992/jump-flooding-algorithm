package util

import (
	"image"
	"math"
)

func searchNearestPixel(x, y int, swap *SwapBuffer, level int) {
	// TODO:
	// exp2(level)をオフセットとする
	step := int(math.Exp2(float64(level)))
	minDistance := 9999.9
	var nearestPoint image.Point
	// 縦3横3の合計8ピクセルから最も近いピクセルを探す、その際ピクセル情報がないピクセルは無視する(この場合だと透過値0のピクセル)
	for ny := -1; ny <= 1; ny++ {
		for nx := -1; nx <= 1; nx++ {
			sampleX := x + nx*step
			sampleY := y + ny*step
			pixel := swap.Get(sampleX, sampleY)
			// TODO: サンプル点と現在のピクセル座標の距離を計算
			sampleDist := 1.0
			if pixel.Nearest != nil && sampleDist < minDistance {
				minDistance = sampleDist
				nearestPoint = *pixel.Nearest
			}
		}
	}
	swap.Set(x, y, swap.Get(nearestPoint.X, nearestPoint.Y))
}

func JumpFlooding(img image.Image) {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	swap := NewSwapBuffer(width, height)
	swap.InitActiveBuffer(img)
	swap.Swap()
	max_level := int(math.Log2(float64(bounds.Max.X))) - 1
	for level := max_level; level >= 0; level-- {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				searchNearestPixel(x, y, swap, level)
			}
		}
		swap.Swap()
	}
}
