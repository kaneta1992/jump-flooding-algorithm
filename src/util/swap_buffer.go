package util

import (
	"image"
	"image/color"
	"math"
)

type PixelInfo struct {
	Inside  bool
	Color   color.RGBA
	Nearest *image.Point
}

type SwapBuffer struct {
	buffers      [2][]PixelInfo
	activeBuffer int
	width        int
	height       int
}

func NewSwapBuffer(w, h int) *SwapBuffer {
	s := &SwapBuffer{activeBuffer: 0, width: w, height: h}
	for index := range s.buffers {
		s.buffers[index] = make([]PixelInfo, w*h)
	}
	return s
}

func (s *SwapBuffer) Swap() {
	s.activeBuffer = (s.activeBuffer + 1) % 2
}

func (s *SwapBuffer) getSize() (int, int) {
	return s.width, s.height
}

func (s *SwapBuffer) getIndex(x, y int) int {
	clampedX, clampedY := s.ClampCoord(x, y)
	return clampedY*s.width + clampedX
}

func (s *SwapBuffer) ClampCoord(x, y int) (int, int) {
	clampedX := int(math.Min(math.Max(0.0, float64(x)), float64(s.width-1)))
	clampedY := int(math.Min(math.Max(0.0, float64(y)), float64(s.height-1)))
	return clampedX, clampedY
}

func (s *SwapBuffer) Get(x, y int) PixelInfo {
	index := s.getIndex(x, y)
	prev := (s.activeBuffer + 1) % 2
	return s.buffers[prev][index]
}

func (s *SwapBuffer) Set(x, y int, info PixelInfo) {
	index := s.getIndex(x, y)
	s.buffers[s.activeBuffer][index] = info
}

func (s *SwapBuffer) SetNearest(x, y int, nearest *image.Point) {
	index := s.getIndex(x, y)
	s.buffers[s.activeBuffer][index].Nearest = nearest
}

func color2RGBA(col color.Color) color.RGBA {
	r, g, b, a := col.RGBA()
	return color.RGBA{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}

func (s *SwapBuffer) InitActiveBuffer(img image.Image) {
	bounds := img.Bounds()

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			col := color2RGBA(img.At(x, y))
			pixel := PixelInfo{
				Inside:  col.A >= 254,
				Color:   col,
				Nearest: nil,
			}
			s.Set(x, y, pixel)
		}
	}
}
