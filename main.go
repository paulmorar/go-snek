package main

import (
	"os"
	"time"

	"gopkg.in/veandco/go-sdl2.v0/sdl"
)

const (
	windowWidth    = 800
	windowHeight   = 600
	baseRectWidth  = 20
	baseRectHeight = 20
)

type drawable interface {
	draw(*sdl.Renderer) error
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Snek game!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	events := make(chan sdl.Event)

	go func() {
		snek := newSnek(baseRectWidth, baseRectHeight)
		mouse := newMouse(windowWidth, windowHeight, baseRectWidth, baseRectHeight, snek.body)
		dir := direction(right)

		tick := time.Tick(100 * time.Millisecond)
		for {
			select {
			case <-tick:
				if err := draw(renderer, snek, mouse); err != nil {
					panic(err)
				}
				update(snek, mouse, dir)
			case event := <-events:
				if _, ok := event.(*sdl.QuitEvent); ok {
					os.Exit(0)
				}
				handleKeyEvent(event, &dir)
			}
		}
	}()

	for {
		events <- sdl.WaitEvent()
	}
}

func draw(renderer *sdl.Renderer, d ...drawable) error {
	if err := renderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE); err != nil {
		return err
	}

	if err := renderer.Clear(); err != nil {
		return err
	}

	for _, v := range d {
		if err := v.draw(renderer); err != nil {
			return err
		}
	}

	renderer.Present()
	return nil
}

func update(s *snek, m *mouse, dir direction) {
	s.update(dir)

	if s.isDead() {
		(*s) = *newSnek(baseRectWidth, baseRectHeight)
		m.eat()
	}

	if s.canEat(m.rect) {
		s.grow()
		m.eat()
	}

	m.update(s.body)
}

func handleKeyEvent(e sdl.Event, dir *direction) {
	v, ok := e.(*sdl.KeyDownEvent)
	if !ok {
		return
	}

	switch {
	case v.Keysym.Sym == sdl.K_RIGHT && *dir != left:
		*dir = right
	case v.Keysym.Sym == sdl.K_LEFT && *dir != right:
		*dir = left
	case v.Keysym.Sym == sdl.K_DOWN && *dir != up:
		*dir = down
	case v.Keysym.Sym == sdl.K_UP && *dir != down:
		*dir = up
	}
}
