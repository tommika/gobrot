package mandelbrot

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultPaletteSize = 32
)

// Controller is the main controller
func Controller() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Printf("Failed to initialize SDL: err=%v\n", err)
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Go Mandelbrot Set Explorer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Printf("Failed to create SDL window: err=%v\n", err)
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		log.Printf("Failed to retreive window surface: err=%v\n", err)
		panic(err)
	}
	log.Printf("Window surface format is %s\n", DescribeSurfaceFormat(surface))

	wpHome := Waypoint{
		X:        -0.75,
		Y:        0.0,
		Scale:    float64(surface.W) / 3.5,
		MaxOrbit: 100}

	var palR, palG, palB = 0.0, 1.0, 1.0

	view := InitView(window, CreatePalette(defaultPaletteSize, palR, palG, palB), func(view *View, rc sdl.Rect) {
		GenerateMandelbrot(view, rc)
	})
	view.GotoWaypoint(wpHome)

	quit := false
	update := true
	dragging := false
	var xdrag, ydrag int32

	for !quit {
		// For now, we'll wait for events (WaitEvent vs PollEvent, since
		// everything will happen in response to some input from the user
		event := sdl.WaitEvent()
		switch evt := event.(type) {
		case *sdl.QuitEvent:
			quit = true
		case *sdl.MouseButtonEvent:
			if evt.Type == sdl.MOUSEBUTTONUP &&
				evt.Button == sdl.BUTTON_LEFT {
				dragging = false
			} else if evt.Type == sdl.MOUSEBUTTONDOWN &&
				evt.Button == sdl.BUTTON_LEFT {
				if evt.Clicks == 1 {
					dragging = true
					xdrag = evt.X
					ydrag = evt.Y
				} else if evt.Clicks >= 2 {
					kms := sdl.GetModState()
					wp := view.GetCurrentWaypoint()
					wp.X, wp.Y = view.FromScreen(evt.X, evt.Y)
					if (kms & sdl.KMOD_SHIFT) != 0 {
						wp.Scale /= 2
					} else {
						wp.Scale *= 2
					}
					view.GotoWaypoint(wp)
					update = true
				}
			}
		case *sdl.MouseMotionEvent:
			if dragging {
				view.Scroll(xdrag-evt.X, evt.Y-ydrag)
				xdrag = evt.X
				ydrag = evt.Y
				update = true
			}
		case *sdl.MouseWheelEvent:
			kms := sdl.GetModState()
			if (kms & sdl.KMOD_SHIFT) != 0 {
				view.Rotate(int(evt.Y) * 4)
			} else if (kms & sdl.KMOD_CTRL) != 0 {
				view.maxorbit += int(evt.Y)
				view.Invalidate()
			} else if (kms & sdl.KMOD_ALT) != 0 {
				n := max(2, len(view.palette)+int(evt.Y))
				view.palette = CreatePalette(n, palR, palG, palB)
				view.Invalidate()
			} else {
				view.Zoom(int(-20 * evt.Y))
			}
			update = true
		case *sdl.KeyboardEvent:
			if evt.Type == sdl.KEYDOWN {
				log.Printf("evt.Keysym.Sym=%v", evt.Keysym.Sym)
				switch evt.Keysym.Sym {
				case sdl.K_ESCAPE:
					quit = true
					break
				case sdl.K_HOME:
					view.palette = CreatePalette(defaultPaletteSize, palR, palG, palB)
					view.GotoWaypoint(wpHome)
					update = true
					break
				case sdl.K_UP:
					view.Scroll(0, 10)
					update = true
					break
				case sdl.K_DOWN:
					view.Scroll(0, -10)
					update = true
					break
				case sdl.K_LEFT:
					view.Scroll(-10, 0)
					update = true
					break
				case sdl.K_RIGHT:
					view.Scroll(10, 0)
					update = true
					break
				case sdl.K_PAGEUP:
					view.Zoom(-50)
					update = true
					break
				case sdl.K_PAGEDOWN:
					view.Zoom(+50)
					update = true
					break
				// FIXME - not getting sdl.K_PLUS on windows or Mac. Seems to not
				// translate a Shift+Equal(=) to a Plus(+)
				case sdl.K_PLUS, sdl.K_KP_PLUS, sdl.K_EQUALS:
					view.maxorbit += 10
					view.Invalidate()
					update = true
					break
				case sdl.K_MINUS, sdl.K_KP_MINUS:
					view.maxorbit -= 10
					view.Invalidate()
					update = true
					break
				}
			}
		}
		if update {
			view.Paint()
			update = false
		}
	}
}
