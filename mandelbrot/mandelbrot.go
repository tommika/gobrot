// Copyright (c) 2019, 2024 Thomas Mikalsen. Subject to the MIT License
package mandelbrot

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func GenerateMandelbrot(view *View, rc sdl.Rect) {
	w, h := view.surf.W, view.surf.H
	var cosRotate, sinRotate float64
	if view.bearing != 0 {
		radRotate := toRadians(view.bearing)
		cosRotate = math.Cos(radRotate)
		sinRotate = math.Sin(radRotate)
	}
	for y := rc.Y; y < rc.Y+rc.H; y++ {
		for x := rc.X; x < rc.X+rc.W; x++ {
			if x < 0 || y < 0 || x >= view.surf.W || y >= view.surf.H {
				continue
			}
			cx := (float64(x) - (float64(w) / 2)) / view.scale
			cy := ((float64(h) / 2) - float64(y)) / view.scale
			if view.bearing != 0 {
				cx, cy = cx*cosRotate-cy*sinRotate,
					cx*sinRotate+cy*cosRotate
			}
			orbit := computeOrbit(cx+view.xcenter, cy+view.ycenter, view.maxorbit)
			var rgba ColorRGBA
			if orbit == view.maxorbit {
				rgba = BLACK
			} else {
				rgba = view.palette[orbit%len(view.palette)]
			}
			SetSurfacePixel(view.surf, int(x), int(y), rgba)
		}
	}
}

func computeOrbit(cx float64, cy float64, maxorbit int) int {
	// z = z^2 + c; c = x + y*i; z^2 = x^2 - y^2 + 2xyi
	var zx, zy float64
	var zx2, zy2 float64 // holds zx^2 and zy2
	for orbit := 0; orbit < maxorbit; orbit++ {
		zx, zy = zx2-zy2+cx, (2*zx*zy)+cy
		zx2, zy2 = zx*zx, zy*zy
		if zx2+zy2 > 4 {
			// this point has escaped - not in the set
			return orbit
		}
	}
	// reached max orbit - this point is in the set
	return maxorbit
}
