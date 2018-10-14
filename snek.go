package main

import "gopkg.in/veandco/go-sdl2.v0/sdl"

type direction int

const (
	left direction = iota + 1
	right
	up
	down
)

type snek struct {
	body              []sdl.Rect
	dead              bool
	previousDirection direction
}

func newSnek(w, h int32) *snek {
	return &snek{
		body: []sdl.Rect{
			{X: 200, Y: 200, W: w, H: h},
			{X: 180, Y: 200, W: w, H: h},
			{X: 160, Y: 200, W: w, H: h},
			{X: 140, Y: 200, W: w, H: h},
		},
	}
}

func (s snek) isDead() bool {
	return s.dead
}

func (s *snek) update(dir direction) {
	if s.isDead() {
		return
	}

	for i := len(s.body) - 1; i >= 1; i-- {
		s.body[i].X, s.body[i].Y = s.body[i-1].X, s.body[i-1].Y
	}

	head := s.head()

	switch {
	case dir == right && s.previousDirection != left:
		head.X += baseRectWidth
		if head.X > windowWidth-baseRectWidth {
			head.X = 0
		}
	case dir == left && s.previousDirection != right:
		head.X -= baseRectWidth
		if head.X < 0 {
			head.X = windowWidth - baseRectWidth
		}
	case dir == down && s.previousDirection != up:
		head.Y += baseRectHeight
		if head.Y > windowHeight-baseRectHeight {
			head.Y = 0
		}
	case dir == up && s.previousDirection != down:
		head.Y -= baseRectHeight
		if head.Y < 0 {
			head.Y = windowHeight - baseRectHeight
		}
	}

	for i := 1; i < len(s.body); i++ {
		if head.HasIntersection(&s.body[i]) {
			s.dead = true
			return
		}
	}

	s.previousDirection = dir
}

func (s snek) draw(renderer *sdl.Renderer) error {
	if s.isDead() {
		return nil
	}

	if err := renderer.SetDrawColor(255, 255, 255, sdl.ALPHA_OPAQUE); err != nil {
		return err
	}

	if err := renderer.FillRects(s.body); err != nil {
		return err
	}

	return renderer.DrawRects(s.body)
}

func (s snek) canEat(mouse *sdl.Rect) bool {
	head := s.head()
	return head.HasIntersection(mouse)
}

func (s *snek) grow() {
	prevX, prevY := s.body[len(s.body)-1].X, s.body[len(s.body)-1].Y
	s.body = append(s.body, sdl.Rect{W: baseRectWidth, H: baseRectHeight, X: prevX, Y: prevY})
}

func (s *snek) head() *sdl.Rect {
	return &s.body[0]
}
