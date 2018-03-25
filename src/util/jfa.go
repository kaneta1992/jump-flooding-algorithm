package util

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"
	"time"
)

type JFA struct {
	buffer *SwapBuffer
}

func clamp(min, max, value float64) float64 {
	return math.Min(math.Max(value, min), max)
}

func length(x, y int) float64 {
	return math.Sqrt(float64(x*x + y*y))
}

func length2(x, y int) int {
	return x*x + y*y
}

func distance(p1, p2 image.Point) float64 {
	x := p2.X - p1.X
	y := p2.Y - p1.Y
	return length(x, y)
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

func jumpFloodingAlgorithm(buffer *SwapBuffer) {
	width, height := buffer.getSize()
	maxLevel := int(math.Log2(float64(width))) - 1
	beforeAt := time.Now()

	for level := maxLevel; level >= 0; level-- {
		wg := &sync.WaitGroup{}
		for y := 0; y < height; y++ {
			wg.Add(1)
			go func(yy int) {
				for x := 0; x < width; x++ {
					searchNearestPixel(x, yy, buffer, level)
				}
				wg.Done()
			}(y)
		}
		wg.Wait()
		buffer.Swap()
	}
	afterAt := time.Now()
	fmt.Println(afterAt.Sub(beforeAt).Nanoseconds())
}

func NewJFA(img image.Image) *JFA {
	bounds := img.Bounds()
	jfa := &JFA{
		buffer: NewSwapBuffer(bounds.Max.X, bounds.Max.Y),
	}
	jfa.buffer.InitActiveBuffer(img)
	jfa.buffer.Swap()
	jumpFloodingAlgorithm(jfa.buffer)
	return jfa
}

func (j *JFA) createImageWithEachPixel(insideFunction, outsideFunction, hasNotNearestFunction func(pixel PixelInfo) color.RGBA) *image.RGBA {
	width, height := j.buffer.getSize()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := j.buffer.Get(x, y)
			var col color.RGBA
			if pixel.Nearest == nil {
				col = hasNotNearestFunction(pixel)
			}
			if pixel.Inside {
				col = insideFunction(pixel)
			} else {
				col = outsideFunction(pixel)
			}
			img.Set(x, y, col)
		}
	}
	return img
}

func (j *JFA) CalcVoronol() *image.RGBA {
	return j.createImageWithEachPixel(func(pixel PixelInfo) color.RGBA {
		return pixel.Color
	}, func(pixel PixelInfo) color.RGBA {
		nearest := pixel.Nearest
		return j.buffer.Get(nearest.X, nearest.Y).Color
	}, func(pixel PixelInfo) color.RGBA {
		return color.RGBA{}
	})
}

func (j *JFA) CalcSDF(maxDistance float64) *image.RGBA {
	return j.createImageWithEachPixel(func(pixel PixelInfo) color.RGBA {
		dist := distance(*pixel.Nearest, pixel.Coord)
		distColor := uint8((1.0 - clamp(0.0, 1.0, dist/maxDistance)) * 127)
		return color.RGBA{A: 255, R: distColor, G: distColor, B: distColor}
	}, func(pixel PixelInfo) color.RGBA {
		dist := distance(*pixel.Nearest, pixel.Coord)
		distColor := 128 + uint8(clamp(0.0, 1.0, dist/maxDistance)*127)
		return color.RGBA{A: 255, R: distColor, G: distColor, B: distColor}
	}, func(pixel PixelInfo) color.RGBA {
		return color.RGBA{}
	})
}
