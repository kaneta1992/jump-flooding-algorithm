package util

import (
	"image"
	"math"
)

func searchNearestPixel(swap *SwapBuffer, level int) {
	// TODO:
	// exp2(level)をオフセットとする
	// 縦3横3の合計8ピクセルから最も近いピクセルを探す、その際ピクセル情報がないピクセルは無視する(この場合だと透過値0のピクセル)
}

func JumpFlooding(img image.Image) {
	bounds := img.Bounds()
	swap := NewSwapBuffer(bounds.Max.X, bounds.Max.Y)
	swap.InitActiveBuffer(img)
	swap.Swap()
	max_level := int(math.Log2(float64(bounds.Max.X))) - 1
	for level := max_level; level >= 0; level-- {
		searchNearestPixel(swap, level)
		swap.Swap()
	}
}
