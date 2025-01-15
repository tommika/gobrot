// Copyright (c) 2019, 2024 Thomas Mikalsen. Subject to the MIT License
package mandelbrot

import (
	"log"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type ImageGenFunc func(view *View, rc sdl.Rect)

type View struct {
	imageGenFunc ImageGenFunc
	window       *sdl.Window
	surf         *sdl.Surface
	surf2        *sdl.Surface // secondary image used for fast scrolling
	palette      Palette
	xcenter      float64
	ycenter      float64
	scale        float64
	bearing      int
	maxorbit     int
	fastscroll   bool
	maxThreads   int
	valid        bool
}

type Waypoint struct {
	X, Y     float64
	Scale    float64
	Bearing  int
	MaxOrbit int
}

func InitView(window *sdl.Window, palette Palette, igf ImageGenFunc) View {
	if igf == nil {
		panic("Image generation function cannot be nil")
	}
	numCPU := runtime.NumCPU()
	log.Printf("NumCPU=%d", numCPU)
	view := View{
		imageGenFunc: igf,
		window:       window,
		palette:      palette,
		valid:        false,
		fastscroll:   true,
		maxThreads:   numCPU}
	return view
}

func (view *View) Invalidate() {
	view.valid = false
}

func (view *View) Paint() {
	if !view.valid {
		view.generateImage()
		view.valid = true
	}
	surf, err := view.window.GetSurface()
	if err != nil {
		log.Printf("Failed to get window surface: %v\n", err)
	} else {
		view.surf.Blit(nil, surf, nil)
	}
	view.window.UpdateSurface()
}

func (view *View) generateImage() bool {
	w, h := view.window.GetSize()
	if view.surf == nil || view.surf.W != w || view.surf.H != h {
		// need to reallocate image
		if view.surf != nil {
			view.surf.Free()
		}
		view.surf = CreateSurface(w, h)
		if view.surf == nil {
			return false
		}
	}
	var wg sync.WaitGroup
	t0 := time.Now()
	dy := int(h) / view.maxThreads
	for y := 0; y < int(h); y += dy {
		wg.Add(1)
		rc := sdl.Rect{X: 0, Y: int32(y), W: w, H: int32(min(dy, int(h)))}
		go func() {
			view.imageGenFunc(view, rc)
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed := time.Since(t0)
	log.Printf("generateImage: elapsed=%v", elapsed)
	return true
}

func (view *View) GotoWaypoint(wp Waypoint) {
	view.xcenter = wp.X
	view.ycenter = wp.Y
	view.scale = wp.Scale
	view.bearing = wp.Bearing
	view.maxorbit = wp.MaxOrbit
	view.Invalidate()
}

func (view *View) GetCurrentWaypoint() Waypoint {
	return Waypoint{
		X:        view.xcenter,
		Y:        view.ycenter,
		Scale:    view.scale,
		Bearing:  view.bearing,
		MaxOrbit: view.maxorbit}
}
func (view *View) Zoom(dxy int) {
	w, _ := view.window.GetSize()
	dx := float64(w)
	view.scale = (dx * view.scale) / (dx + float64(dxy))
	view.Invalidate()
}

func (view *View) Rotate(deg int) {
	view.bearing = (view.bearing + 360 + deg) % 360
	view.Invalidate()
}

func (view *View) Scroll(dxscroll, dyscroll int32) {
	if dxscroll == 0 && dyscroll == 0 {
		return
	}
	// update center of view
	dx := float64(dxscroll) / view.scale
	dy := float64(dyscroll) / view.scale
	if view.bearing != 0 {
		rad := toRadians(view.bearing)
		if rad != 0 {
			dx, dy = rotate(dx, dy, rad)
		}
	}
	view.xcenter += dx
	view.ycenter += dy

	if !view.fastscroll || !view.valid {
		view.Invalidate()
		return
	}

	t0 := time.Now()
	w, h := view.window.GetSize()
	if view.surf2 == nil || view.surf2.W != w || view.surf2.H != h {
		// need to reallocate image
		if view.surf2 != nil {
			view.surf2.Free()
		}
		view.surf2 = CreateSurface(w, h)
		if view.surf2 == nil {
			view.fastscroll = false
			view.Invalidate()
			return
		}
	}

	// "scroll" the image by copying current image second image
	rcsrc := sdl.Rect{X: dxscroll, Y: -dyscroll, W: w, H: h}
	rcdst := sdl.Rect{X: 0, Y: 0, W: w, H: h}
	view.surf.Blit(&rcsrc, view.surf2, &rcdst)

	// swap the primary and secondary surfaces
	view.surf, view.surf2 = view.surf2, view.surf

	// render the areas that have been exposed by the scroll
	var wg sync.WaitGroup
	if dxscroll != 0 {
		rc := sdl.Rect{Y: 0, H: h}
		if dxscroll < 0 {
			rc.X = 0
			rc.W = -dxscroll
		} else {
			rc.X = w - dxscroll
			rc.W = dxscroll
		}
		if view.maxThreads <= 1 {
			view.imageGenFunc(view, rc)
		} else {
			wg.Add(1)
			go func() {
				view.imageGenFunc(view, rc)
				wg.Done()
			}()
		}
	}

	if dyscroll != 0 {
		rc := sdl.Rect{X: 0, W: w}
		if dyscroll < 0 {
			rc.Y = h + dyscroll
			rc.H = -dyscroll
		} else {
			rc.Y = 0
			rc.H = dyscroll
		}
		// TODO  subtract the area that was rendered in step above (if any)
		if view.maxThreads <= 1 {
			view.imageGenFunc(view, rc)
		} else {
			wg.Add(1)
			go func() {
				view.imageGenFunc(view, rc)
				wg.Done()
			}()
		}
	}

	wg.Wait()
	elapsed := time.Since(t0)
	log.Printf("Scroll: elapsed=%v", elapsed)
}

func (view *View) FromScreen(xS, yS int32) (float64, float64) {
	W, H := view.window.GetSize()
	x := (float64(xS) - (float64(W) / 2)) / view.scale
	y := ((float64(H) / 2) - float64(yS)) / view.scale
	if view.bearing != 0 {
		rad := toRadians(view.bearing)
		if rad != 0 {
			x, y = rotate(x, y, rad)
		}
	}
	return x + view.xcenter, y + view.ycenter
}

func rotate(x, y float64, angle float64) (float64, float64) {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	xn := x*cos - y*sin
	yn := x*sin + y*cos
	return xn, yn
}

func toRadians(deg int) float64 {
	rad := math.Pi * float64(deg) / 180.0
	return rad
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}
