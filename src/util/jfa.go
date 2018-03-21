package util

import (
	"fmt"
	"image"
	"math"
	"sync"
	"time"
)

func length2(x, y int) int {
	return x*x + y*y
}

func searchNearestPixel(x, y int, swap *SwapBuffer, level int) {
	step := int(math.Exp2(float64(level)))
	minDistance := 999999999
	currentPixel := swap.Get(x, y)
	for ny := -1; ny <= 1; ny++ {
		for nx := -1; nx <= 1; nx++ {
			sampleX, sampleY := swap.ClampCoord(x+nx*step, y+ny*step)
			pixel := swap.Get(sampleX, sampleY)
			var nearestX, nearestY int
			if currentPixel.Inside != pixel.Inside {
				// ピクセルの内部外部が異なる場合は、サンプリング点で距離を計算する
				nearestX = sampleX
				nearestY = sampleY
			} else if pixel.Nearest != nil {
				// 同じですでに近傍ピクセルを見つけている場合は近傍ピクセルとの距離を計算する
				nearestX = pixel.Nearest.X
				nearestY = pixel.Nearest.Y
			} else {
				// それ以外は無視
				continue
			}
			sampleDist := length2(nearestX-x, nearestY-y)
			if sampleDist < minDistance {
				minDistance = sampleDist
				currentPixel.Nearest = &image.Point{
					X: nearestX,
					Y: nearestY,
				}
			}
		}
	}
	swap.Set(x, y, currentPixel)
}

func JumpFlooding(img image.Image) *SwapBuffer {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	swap := NewSwapBuffer(width, height)
	swap.InitActiveBuffer(img)
	swap.Swap()
	maxLevel := int(math.Log2(float64(bounds.Max.X))) - 1
	beforeAt := time.Now()

	for level := maxLevel; level >= 0; level-- {
		wg := &sync.WaitGroup{}
		for y := 0; y < height; y++ {
			wg.Add(1)
			go func(yy int) {
				for x := 0; x < width; x++ {
					searchNearestPixel(x, yy, swap, level)
				}
				wg.Done()
			}(y)
		}
		wg.Wait()
		swap.Swap()
	}
	afterAt := time.Now()
	fmt.Println(afterAt.Sub(beforeAt).Nanoseconds())
	return swap
}
