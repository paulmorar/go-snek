package main

import (
	"math/rand"

	"gopkg.in/veandco/go-sdl2.v0/sdl"
)

type mouse struct {
	rect                      *sdl.Rect
	eaten                     bool
	windowWidth, windowHeight int32
	w, h                      int32
}

func newMouse(windowWidth, windowHeight, w, h int32, blackList []sdl.Rect) *mouse {
	length := (windowWidth / baseRectWidth) * (windowHeight / baseRectHeight)

	tmp := make(map[sdl.Rect]struct{}, len(blackList))

	for _, r := range blackList {
		tmp[r] = struct{}{}
	}

	freeRects := make([]sdl.Rect, 0, length)

	x := int32(0)
	y := int32(0)
	for i := int32(0); i < length; i++ {
		rect := sdl.Rect{X: x, Y: y, W: baseRectWidth, H: baseRectHeight}
		if _, ok := tmp[rect]; !ok {
			freeRects = append(freeRects, rect)
		}

		x += baseRectWidth
		if x > windowWidth-baseRectWidth {
			y += baseRectHeight
			x = 0
		}
	}

	index := rand.Intn(len(freeRects) - 1)

	return &mouse{
		rect:         &freeRects[index],
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		w:            w,
		h:            h,
	}
}

func (m *mouse) update(blackList []sdl.Rect) {
	if m.isEaten() {
		(*m) = *newMouse(m.windowWidth, m.windowHeight, m.w, m.h, blackList)
	}
}

func (m mouse) draw(renderer *sdl.Renderer) error {
	if m.isEaten() {
		return nil
	}

	if err := renderer.SetDrawColor(255, 255, 255, sdl.ALPHA_OPAQUE); err != nil {
		return err
	}

	if err := renderer.FillRect(m.rect); err != nil {
		return err
	}

	return renderer.DrawRect(m.rect)
}

func (m mouse) isEaten() bool {
	return m.eaten
}

func (m *mouse) eat() {
	m.eaten = true
}
