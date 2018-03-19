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
	var nearestPoint image.Point
	for ny := -1; ny <= 1; ny++ {
		for nx := -1; nx <= 1; nx++ {
			sampleX := x + nx*step
			sampleY := y + ny*step
			pixel := swap.Get(sampleX, sampleY)
			if pixel.Nearest != nil {
				sampleDist := length2(pixel.Nearest.X-x, pixel.Nearest.Y-y)
				if sampleDist < minDistance {
					minDistance = sampleDist
					nearestPoint = *pixel.Nearest
				}
			}
		}
	}
	swap.Set(x, y, swap.Get(nearestPoint.X, nearestPoint.Y))
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
