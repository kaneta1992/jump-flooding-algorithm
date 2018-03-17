package util

import (
	"image"
	"image/draw"
)

// TODO:
// buffersをJFA用の構造体に変更する

type SwapBuffer struct {
	buffers      [2]*image.RGBA
	activeBuffer int
}

func NewSwapBuffer(width, height int) *SwapBuffer {
	s := &SwapBuffer{activeBuffer: 0}
	for index := range s.buffers {
		s.buffers[index] = image.NewRGBA(image.Rect(0, 0, width, height))
	}
	return s
}

func (s *SwapBuffer) Swap() {
	s.activeBuffer = (s.activeBuffer + 1) % 2
}

func (s *SwapBuffer) GetActiveBuffer() *image.RGBA {
	return s.buffers[s.activeBuffer]
}

func (s *SwapBuffer) GetPrevBuffer() *image.RGBA {
	index := (s.activeBuffer + 1) % 2
	return s.buffers[index]
}

func (s *SwapBuffer) SetActiveBuffer(img image.Image) {
	activeBuffer := s.GetActiveBuffer()
	draw.Draw(activeBuffer, activeBuffer.Bounds(), img, image.ZP, draw.Src)
}
